DATE = $(shell date +%FT%T%Z)
BUILD_DIR = build/bin
GIT_VER=$(shell git rev-parse HEAD)
GO = CGO_ENABLED=0 go
MODULE=$(shell head -n1 go.mod | awk -F ' ' '{print $$2}')
APPS=$(shell ls cmd)
GOLINTCI_FILE=golangci-lint
LINT_PATH=$(shell pwd)/bin
GOLINTCI=$(LINT_PATH)/$(GOLINTCI_FILE)

LDFLAGS=-ldflags "-s -w \
	 -X ${MODULE}/pkg/version.hash=${GIT_VER} \
	 -X ${MODULE}/pkg/version.buildtimestamp=${DATE}"

define go-build
GOOS=$(2) GOARCH=$(3) $(GO) build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/$(1)-$(2)-$(3) -v cmd/$(1)/*.go
endef

define splitter
$(strip $(word $(2),$(subst _, ,$(1))))
endef

.PHONY: all
all: vet lint test # test-integration
	$(foreach app,$(APPS),make build_$(app)_linux_amd64 build_$(app)_windows_amd64;)

install-tools:
	# NOTE: pin golangci-lint version to support go 1.18
	# https://github.com/golangci/golangci-lint/issues/2649#issuecomment-1092784768
	GOBIN=$(LINT_PATH) $(GO) install github.com/golangci/$(GOLINTCI_FILE)/cmd/$(GOLINTCI_FILE)@f5b92e1

.PHONY: build_%
build_%:
	$(call go-build,$(call splitter,$@,2),$(call splitter,$@,3),$(call splitter,$@,4))

.PHONY: lint
lint:
	$(GOLINTCI) run ./pkg/... ./cmd/...

.PHONY: vet
vet:
	$(GO) vet -composites=false ./pkg/... ./cmd/...

.PHONY: test
test:
	$(GO) test -v -coverprofile cover.out ./cmd/... ./pkg/...

.PHONY: test-integration
test-integration:
		$(GO) test -v ./tests/...

.PHONY: clean
clean:
	-rm -f ${BUILD_DIR}/${BINARY}-*

.PHONY: distclean
distclean:
	rm -rf ./build

.PHONY: mrproper
mrproper: distclean
	git ls-files --others | xargs rm -rf
