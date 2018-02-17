# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

PACKAGE := github.com/overwaifu/overwaifu/cmd/overwaifu

GOEXE ?= go

.PHONY: vendor build docker test clean help
.DEFAULT_GOAL := help

vendor: ## Install deps and sync vendored dependencies
	@echo "Installing vendored dependencies"
	@${GOEXE} version
	@dep ensure

build: ## Build overwaifu binary
	@echo "Building overwaifu binary"
	@${GOEXE} build ${PACKAGE}

test: ## Run tests
	@${GOEXE} test

clean: ## Clean the project
	@echo "Cleaning the project"
	@rm overwaifu.exe

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
