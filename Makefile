default: fmt lint tidy build test

# globals
GOBIN?=$(GOPATH)/bin
export GOBIN

LOCAL_BIN?=$(CURDIR)/build
NAME=app

# lint
GOLANGCI_BIN=$(GOBIN)/golangci-lint
$(GOLANGCI_BIN):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# gofumpt
GOFUMPT_BIN=$(GOBIN)/gofumpt
$(GOFUMPT_BIN):
	go install mvdan.cc/gofumpt@latest

# gci
GCI_BIN=$(GOBIN)/gci
$(GCI_BIN):
	go install github.com/daixiang0/gci@latest

# mocks
MOCKERY_BIN=$(GOBIN)/mockery
$(MOCKERY_BIN):
	go install github.com/vektra/mockery/v3@latest

.PHONY: build
build: BIN?=$(LOCAL_BIN)/$(NAME)
build:
	go build -o $(BIN) ./main.go

.PHONY: lint
lint: $(GOLANGCI_BIN)
	$(GOLANGCI_BIN) run

.PHONY: generate
generate: $(MOCKERY_BIN)
	go generate ./...

.PHONY: test
test: generate
	go test ./... -cover -race -coverprofile cover.out  && go tool cover -func cover.out

.PHONY: fmt
fmt: $(GOFUMPT_BIN) $(GCI_BIN)
	$(GOFUMPT_BIN) -l -w .
	$(GCI_BIN) write \
		-s standard \
		-s default \
		-s "prefix(github.com/DivPro/app)" \
		--skip-generated .

.PHONY: tidy
tidy:
	go mod tidy
