release:
	git tag -a $(VER) -m '$(VER)'
	git push origin --tags
	echo Tagged release with $(VER)
.PHONY: release