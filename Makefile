release:
	git tag -a v$(VER) -m 'v$(VER)'
	git push origin
	echo Tagged release with $(VER)
.PHONY: release

list-releases:
	git ls-remote --tags
.PHONY: list-release