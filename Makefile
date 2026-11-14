# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

build:
	mkdir -p bin
	go build -o ./bin/chainsmith ./cmd/chainsmith

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: ## Run tests.
	go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

.PHONY: test-e2e
test-e2e:
	go test -v ./test/e2e

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

# Same version as the go-test-coverage action in .github/workflows/ci-run-on-pr.yaml: Dependabot bumps that one,
# update this one with it. Since v2.20.0 it needs Go 1.27 to install (GOTOOLCHAIN=auto downloads it).
GO_TEST_COVERAGE_VERSION ?= v2.23.0

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@$(GO_TEST_COVERAGE_VERSION)

.PHONY: test-coverage
test-coverage: ## Run the unit tests with a coverage profile (cover.out).
	go test $$(go list ./... | grep -v /e2e) -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

.PHONY: check-coverage
check-coverage: install-go-test-coverage test-coverage
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml
