.PHONY: help
help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## builds the app
	go build  -o out/focus ./cmd/focus

.PHONY: run
run: build ## runs the app
	./out/focus

.PHONY: clean
clean: ## cleans the binary
	rm -f out/focus

.PHONY: install
install: ## installs the app
	go install github.com/nendix/focus/cmd/focus@latest

.PHONY: release
release:  ## creates a release on github (e.g. make release name=v1.3.7)
	git push
	gh release create $(name) --generate-notes
	git fetch --all

.PHONY: test
test: ## runs the tests
	go test ./...
