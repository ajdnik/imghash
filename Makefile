GOLANGCI_LINT_VERSION := v2.10.1

.PHONY: fmt vet test coverage lint lint-install vulncheck vulncheck-install all

all: fmt vet lint vulncheck test

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

lint: lint-install
	golangci-lint run ./...

lint-install:
	@which golangci-lint > /dev/null 2>&1 || { \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	}

vulncheck: vulncheck-install
	$$(go env GOPATH)/bin/govulncheck ./...

vulncheck-install:
	@test -f $$(go env GOPATH)/bin/govulncheck || { \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	}
