.DEFAULT_GOAL := help

.PHONY: test
test: ## Run tests
	@go test -v -race ./...

.PHONY: example
example: ## Run example
	@go run example/main.go $(ARGS)

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
