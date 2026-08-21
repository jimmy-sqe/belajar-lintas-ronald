.PHONY: help build run test test-unit test-integration test-cover lint fmt mocks swag-gen oapi-gen clean setup seed-postgres seed-mysql seed-mongodb buf-generate

# Tools pinned via go run @version (no global install)
# golangci-lint v2.x supports Go 1.25 source; force GOTOOLCHAIN so `go run`
# rebuilds it with the toolchain go.mod declares (the v2.5.0 module pins
# go 1.24.x in its own go.mod which would otherwise be the build toolchain).
GOLANGCI_LINT := GOTOOLCHAIN=go1.25.0 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
MOCKERY       := go run github.com/vektra/mockery/v2@v2.50.0
MIGRATE       := go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

APP_NAME ?= backend-belajar-lintas-ronald

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

setup: ## Install per-project tools and pre-commit hooks
	@command -v mise >/dev/null && mise install || echo "(mise not installed; install via https://mise.jdx.dev)"
	@command -v pre-commit >/dev/null && pre-commit install --hook-type commit-msg --hook-type pre-push --hook-type pre-commit || echo "(pre-commit not installed; install via 'mise install' or pip)"
	@git config commit.template .gitmessage
	@echo "Setup complete. Hooks installed; commit template set."

build: ## Build the binary
	go build -o bin/$(APP_NAME) .

run: ## Run the service locally
	APP_ENV=local go run . serve

test: test-unit ## Default test target (unit only)

test-unit: ## Run unit tests only
	go test -short -count=1 ./...

test-integration: ## Run integration tests (build tag 'integration')
	go test -count=1 -tags=integration ./internal/adapter/...

test-cover: ## Run tests with coverage on tested packages
	# Discover only packages that actually contain test files. On Go 1.25,
	# expanding `./...` includes packages with no tests, which triggers a
	# `go: no such tool "covdata"` failure during coverage aggregation.
	@pkgs=$$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...); \
	  go test -short -count=1 -coverprofile=cover.out $$pkgs
	@go tool cover -func=cover.out | tail -1

lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run --timeout 3m

fmt: ## Format Go sources
	gofmt -w .

mocks: ## Generate mocks (placeholder; Phase 3 adds .mockery.yaml)
	@echo "Phase 3 will configure mockery"

# boilerplate:axis=api-doc option=swag START
SWAG  := go run github.com/swaggo/swag/cmd/swag@v1.16.4

swag-gen: ## (api-doc=swag) Generate swagger docs from annotations
	$(SWAG) init -g main.go -o docs/swagger
# boilerplate:axis=api-doc option=swag END

# boilerplate:axis=api-doc option=oapi-codegen START
OAPI  := go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

oapi-gen: ## (api-doc=oapi-codegen) Generate types + server from openapi.yaml
	$(OAPI) -generate types,echo-server -package generated \
	  internal/http/apidoc/openapi/openapi.yaml \
	  > internal/http/apidoc/openapi/generated.go
# boilerplate:axis=api-doc option=oapi-codegen END

# boilerplate:axis=rpc option=grpc START
.PHONY: buf-generate
buf-generate: ## Regenerate gRPC stubs from proto/ into gen/
	buf generate
# boilerplate:axis=rpc option=grpc END

clean: ## Remove build artifacts
	rm -rf bin/ cover.out

# boilerplate:axis=persistence option=postgres START
seed-postgres: ## (persistence=postgres) Run seeds against Postgres
	APP_ENV=$${APP_ENV:-local} go run . db:seed
# boilerplate:axis=persistence option=postgres END

# boilerplate:axis=persistence option=mysql START
seed-mysql: ## (persistence=mysql) Run seeds against MySQL
	APP_ENV=$${APP_ENV:-local} go run . db:seed
# boilerplate:axis=persistence option=mysql END

# boilerplate:axis=persistence option=mongodb START
seed-mongodb: ## (persistence=mongodb) Run seeds against MongoDB
	APP_ENV=$${APP_ENV:-local} go run . db:seed
# boilerplate:axis=persistence option=mongodb END
