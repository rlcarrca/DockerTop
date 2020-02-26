BUILD_DATE=`date +%FT%T%z`

.PHONY: help
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ## - Check the version
	@echo "${BUILD_DATE}"

.PHONY: build
build: ## - Build dockertop
	go build dockertop.go dockerlist.go dockerdiff.go

.PHONY: clean
clean: ## - Remove dockertop
	@rm dockertop

