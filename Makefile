GOLANGCI_LINT_VERSION := v2.10.1

.PHONY: fmt vet test fuzz coverage lint lint-install vulncheck vulncheck-install all

FUZZ_TIME ?= 10s

all: fmt vet lint vulncheck test

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

# Retries bare "context deadline exceeded" failures: those come from a race in
# the Go fuzzing coordinator at the -fuzztime deadline, not from the target.
# A real crash always writes a failing input to testdata, so it is never retried.
fuzz:
	@grep -o 'func Fuzz[A-Za-z0-9_]*' fuzz_test.go | sed 's/^func //' | while read -r target; do \
		echo "Fuzzing $$target for $(FUZZ_TIME)..."; \
		for attempt in 1 2 3; do \
			out=$$(go test -run='^$$' -fuzz="^$${target}$$" -fuzztime=$(FUZZ_TIME) 2>&1); \
			code=$$?; \
			echo "$$out"; \
			[ $$code -eq 0 ] && break; \
			echo "$$out" | grep -q 'Failing input written to' && exit 1; \
			echo "$$out" | grep -q 'context deadline exceeded' || exit 1; \
			echo "Fuzzing flake (attempt $$attempt/3), retrying..."; \
			[ $$attempt -eq 3 ] && exit 1; \
		done; \
	done

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
