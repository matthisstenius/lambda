package lambda

import (
	"github.com/aws/aws-sdk-go/aws/session"
	sdk "github.com/aws/aws-sdk-go/service/lambda"

	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/mstenius/logger"
	"github.com/aws/aws-sdk-go/aws"
)

// Input data for current invocation
type Input struct {
	Service  string
	Resource string
	Body     interface{}
}

// Invoke Lambda function with payload based on input
func Invoke(input Input, out interface{}) error {
	sess, _ := session.NewSession()
	client := sdk.New(sess)

	body, err := json.Marshal(input.Body)
	if err != nil {
		logger.WithFields(logger.Fields{"error": err, "body": input.Body}).Error("Lambda::Invoke() error during marshal")
		return errors.New("could not marshal body")
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"resource":   input.Resource,
		"body":       string(body),
		"httpMethod": http.MethodPost,
	})

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
			"payload": output.Payload,
		}).Error("Lambda::Invoke() error during unmarshal of output payload")
		return errors.New("error during unmarshal")
	}

	if response.StatusCode != http.StatusOK {
		logger.WithFields(logger.Fields{
			"error":   err,
			"payload": output.Payload,
		}).Warning("Lambda::Invoke() none 200 status code")
		return errors.New("none 200 status code")
	}

	if err := json.Unmarshal([]byte(response.Body), out); err != nil {
		logger.WithFields(logger.Fields{
			"error":   err,
			"payload": response.Body,
		}).Error("Lambda::Invoke() error during unmarshal of response body")
		return errors.New("error during unmarshal")
	}
	return nil
}
