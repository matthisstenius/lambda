release:
	git tag -a v$(VER) -m 'v$(VER)'
	git push origin --tags
	echo Tagged release with $(VER)
.PHONY: release