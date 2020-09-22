# Make does not offer a recursive wildcard function, so here's one:
rwildcard=$(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
GO := GO111MODULE=on go
GO_NOMOD := GO111MODULE=off go
GOTEST := $(GO) test

GO_DEPENDENCIES := $(call rwildcard,./,*.go)

.PHONY: all
all: build test check

.PHONY: test
test:
	CGO_ENABLED=$(CGO_ENABLED) $(GOTEST) -short ./...

.PHONY: check
check: fmt lint sec

get-fmt-deps: ## Install test dependencies
	$(GO_NOMOD) get golang.org/x/tools/cmd/goimports

.PHONY: importfmt
importfmt: get-fmt-deps
	@echo "FORMATTING IMPORTS"
	@goimports -w $(GO_DEPENDENCIES)

.PHONY: fmt
fmt: importfmt
	@echo "FORMATTING SOURCE"
	FORMATTED=`$(GO) fmt ./...`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	$(GO_NOMOD) get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINT)
	@echo "VETTING"
	$(GO) vet ./...
	@echo "LINTING"
	$(GOLINT) -set_exit_status ./...

GOSEC := $(GOPATH)/bin/gosec
$(GOSEC):
	$(GO_NOMOD) get -u github.com/securego/gosec/cmd/gosec

.PHONY: sec
sec: $(GOSEC)
	@echo "SECURITY SCANNING"
	$(GOSEC) -quiet -fmt=csv ./...

.PHONY: clean
clean:
	rm -rf bin build release

.PHONY: build
build:
	$(GO) build -i ./...
