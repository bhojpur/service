################################################################################
# Variables                                                                    #
################################################################################

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0
BINARIES    ?= svcsvr svcutl
HA_MODE     ?= false
# Force in-memory log for placement
FORCE_INMEM ?= true

# Add latest tag if LATEST_RELEASE is true
LATEST_RELEASE ?=

PROTOC ?=protoc
# name of protoc-gen-go when protoc-gen-go --version is run.
PROTOC_GEN_GO_NAME = "protoc-gen-go"

ifdef API_VERSION
	RUNTIME_API_VERSION = $(API_VERSION)
else
	RUNTIME_API_VERSION = 1.0
endif

ifdef REL_VERSION
	APP_VERSION := $(REL_VERSION)
else
	APP_VERSION := edge
endif

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),arm64)
	TARGET_ARCH_LOCAL=arm64
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

ifeq ($(GOARCH),amd64)
	LATEST_TAG=latest
else
	LATEST_TAG=latest-$(GOARCH)
endif

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL ?= windows
   BINARY_EXT_LOCAL = .exe
   PROTOC_GEN_GO_NAME := "protoc-gen-go.exe"
endif
export GOOS ?= $(TARGET_OS_LOCAL)

PROTOC_GEN_GO_NAME+= "v1.27.1"

ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
GOLANGCI_LINT:=golangci-lint.exe
export ARCHIVE_EXT = .zip
else
BINARY_EXT_LOCAL:=
GOLANGCI_LINT:=golangci-lint
export ARCHIVE_EXT = .tar.gz
endif

export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

################################################################################
# Target: test                                                                 #
################################################################################
.PHONY: test
test:
	CGO_ENABLED=$(CGO) go test ./... $(COVERAGE_OPTS) $(BUILDMODE)

################################################################################
# Target: lint                                                                 #
################################################################################
.PHONY: lint
lint:
	# Due to https://github.com/golangci/golangci-lint/issues/580, we need to add --fix for windows
	$(GOLANGCI_LINT) run --timeout=20m

################################################################################
# Target: modtidy-all                                                          #
################################################################################
MODFILES := $(shell find . -name go.mod)

define modtidy-target
.PHONY: modtidy-$(1)
modtidy-$(1):
	cd $(shell dirname $(1)); go mod tidy -compat=1.17; cd -
endef

# Generate modtidy target action for each go.mod file
$(foreach MODFILE,$(MODFILES),$(eval $(call modtidy-target,$(MODFILE))))

# Enumerate all generated modtidy targets
# Note that the order of execution matters: root and tests/certification go.mod
# are dependencies in each certification test. This order is preserved by the
# tree walk when finding the go.mod files.
TIDY_MODFILES:=$(foreach ITEM,$(MODFILES),modtidy-$(ITEM))

# Define modtidy-all action trigger to run make on all generated modtidy targets
.PHONY: modtidy-all
modtidy-all: $(TIDY_MODFILES)

################################################################################
# Target: modtidy                                                              #
################################################################################
.PHONY: modtidy
modtidy:
	go mod tidy

################################################################################
# Target: init-proto                                                            #
################################################################################
.PHONY: init-proto
init-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

################################################################################
# Target: gen-proto                                                            #
################################################################################
GRPC_PROTOS:=version auth mvcc membership keyed network registry
PROTO_PREFIX:=github.com/bhojpur/service/pkg/core

# Generate archive files for each binary
# $(1): the binary name to be archived
define genProtoc
.PHONY: gen-proto-$(1)
gen-proto-$(1):
	$(PROTOC) --proto_path=. --proto_path=$(GOPATH)/src --proto_path=$(GOPATH)/src/googleapis --go_out=./pkg/core --go_opt=module=$(PROTO_PREFIX) --go-grpc_out=./pkg/core --go-grpc_opt=require_unimplemented_servers=false,module=$(PROTO_PREFIX) ./pkg/core/v1/$(1)/*.proto
endef

$(foreach ITEM,$(GRPC_PROTOS),$(eval $(call genProtoc,$(ITEM))))

GEN_PROTOS:=$(foreach ITEM,$(GRPC_PROTOS),gen-proto-$(ITEM))

.PHONY: gen-proto
gen-proto: check-proto-version $(GEN_PROTOS) modtidy

################################################################################
# Target: check-diff                                                           #
################################################################################
.PHONY: check-diff
check-diff:
	git diff --exit-code -- '*go.mod' # check no changes
	git diff --exit-code -- '*go.sum' # check no changes

################################################################################
# Target: check-proto-version                                                  #
################################################################################
.PHONY: check-proto-version
check-proto-version: ## Checking the version of proto related tools
	@test "$(shell protoc --version)" = "libprotoc 3.19.4" \
	|| { echo "please use protoc 3.19.4 to generate proto, see https://github.com/bhojpur/service/blob/master/pkg/core/README.md#proto-client-generation"; exit 1; }

	@test "$(shell protoc-gen-go-grpc --version)" = "protoc-gen-go-grpc 1.1.0" \
	|| { echo "please use protoc-gen-go-grpc 1.1.0 to generate proto, see https://github.com/bhojpur/service/blob/master/pkg/core/README.md#proto-client-generation"; exit 1; }

	@test "$(shell protoc-gen-go --version 2>&1)" = "$(PROTOC_GEN_GO_NAME)" \
	|| { echo "please use protoc-gen-go v1.27.0 to generate proto, see https://github.com/bhojpur/service/blob/master/pkg/core/README.md#proto-client-generation"; exit 1; }

################################################################################
# Target: check-proto-diff                                                     #
################################################################################
.PHONY: check-proto-diff
check-proto-diff:
	git diff --exit-code ./pkg/core/v1/auth/auth.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/keyed/keyedserver.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/keyed/raft_internal.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/keyed/rpc.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/keyed/rpc_grpc.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/membership/membership.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/mvcc/kv.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/network/connection.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/network/connection_grpc.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/network/connectioncontext.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/network/networkservice.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/network/networkservice_grpc.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/registry/registry.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/registry/registry_grpc.pb.go # check no changes
	git diff --exit-code ./pkg/core/v1/version/version.pb.go # check no changes

################################################################################
# Target: conf-tests                                                           #
################################################################################
.PHONY: conf-tests
conf-tests:
	CGO_ENABLED=$(CGO) go test -v -tags=conftests -count=1 ./tests/conformance

################################################################################
# Target: E2E-tests-zeebe                                                      #
################################################################################
.PHONY: e2e-tests-zeebe
e2e-tests-zeebe:
	CGO_ENABLED=$(CGO) go test -v -tags=e2etests -count=1 ./tests/e2e/bindings/zeebe/...