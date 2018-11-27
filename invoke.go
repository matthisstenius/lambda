package lambda

import (
	"github.com/aws/aws-sdk-go/aws/session"
	sdk "github.com/aws/aws-sdk-go/service/lambda"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/matthisstenius/logger"
)

// PathParameters holds path params
type PathParameters map[string]interface{}

// QueryParameters holds query string params
type QueryParameters map[string]interface{}

// AuthParameters holds auth params that
type AuthParameters map[string]interface{}

// Input data for current invocation
type Input struct {
	Service        string
	Resource       string
	Body           interface{}
	Method         string
	PathParams     PathParameters
	QueryParams    QueryParameters
	AuthParameters AuthParameters
}

// Invoke Lambda function with payload based on input
func Invoke(input Input, out interface{}) error {
	sess, _ := session.NewSession()
	client := sdk.New(sess)

	payload, err := encodePayload(input)
	if err != nil {
		logger.WithFields(logger.Fields{"error": err}).Error("Lambda::Invoke() Could not encode payload ")
		return errors.New("could not encode payload")
	}

	reqInput := sdk.InvokeInput{
		FunctionName: aws.String(input.Service),
		Payload:      payload,
	}

	req, output := client.InvokeRequest(&reqInput)
	if err := req.Send(); err != nil {
		logger.WithFields(logger.Fields{"error": err}).Error("Lambda::Invoke() Lambda client error")
		return errors.New("lambda client error")
	}

	var response struct {
		StatusCode int
		Body       string
	}
	if err := json.Unmarshal(output.Payload, &response); err != nil {
		logger.WithFields(logger.Fields{
			"error":   err,
			"payload": string(output.Payload),
		}).Error("Lambda::Invoke() error during unmarshal of output payload")
		return errors.New("error during unmarshal")
	}

	if response.StatusCode != http.StatusOK {
		logger.WithFields(logger.Fields{
			"error":   err,
			"payload": string(output.Payload),
		}).Warning("Lambda::Invoke() none 200 status code")
		return errors.New("none 200 status code")
	}

	if err := json.Unmarshal([]byte(response.Body), out); err != nil {
		logger.WithFields(logger.Fields{
			"error":   err,
			"payload": string(response.Body),
		}).Error("Lambda::Invoke() error during unmarshal of response body")
		return errors.New("error during unmarshal")
	}
	return nil
}

func encodePayload(input Input) ([]byte, error) {
	method := http.MethodGet
	if input.Method != "" {
		method = input.Method
	}

	payload := map[string]interface{}{
		"resource":   input.Resource,
		"body":       input.Body,
		"httpMethod": method,
	}

	if len(input.AuthParameters) > 0 {
		payload["requestContext"] = map[string]interface{}{
			"authorizer": input.AuthParameters,
		}
	}
	if len(input.PathParams) > 0 {
		payload["pathParameters"] = input.PathParams
	}
	if len(input.QueryParams) > 0 {
		payload["queryStringParameters"] = input.QueryParams
	}
	return json.Marshal(payload)
}
