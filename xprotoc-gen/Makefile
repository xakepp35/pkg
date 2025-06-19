PROTO_OUT?=api
INPUT_PROTO?=proto
PROTO_FILES=$(shell find $(INPUT_PROTO) -name "*.proto")

OUT=$(PWD)/bin
GO_API_OUT?=api-go
GO_OUT=$(PWD)/$(GO_API_OUT)
PATH := $(OUT):$(PATH)

PROTOC_VERSION := v30.2
PROTOC_ZIP     := protoc-30.2-linux-x86_64.zip

# FIBER OPTIONS
FIBER_OPTS_PATHS                    ?= source_relative
FIBER_OPTS_ERROR_HANDLERS_PACKAGE   ?= github.com/xakepp35/pkg/xprotoc-gen/utils
FIBER_OPTS_JSON_UNMARSHAL_PACKAGE   ?= encoding/json
FIBER_OPTS_GRPC_ERROR_FUNC          ?= HandleGRPCStatusError
FIBER_OPTS_UNMARSHAL_ERROR_FUNC     ?= HandleUnmarshalError
FIBER_OPTS_VALIDATION_ERROR_FUNC    ?= HandleValidationError

PROTOC_GEN_GO_FIBER_OPTS = \
  paths=$(FIBER_OPTS_PATHS) \
  error_handlers_package=$(FIBER_OPTS_ERROR_HANDLERS_PACKAGE) \
  json_unmarshal_package=$(FIBER_OPTS_JSON_UNMARSHAL_PACKAGE) \
  grpc_error_handle_func=$(FIBER_OPTS_GRPC_ERROR_FUNC) \
  unmarshal_error_handle_func=$(FIBER_OPTS_UNMARSHAL_ERROR_FUNC) \
  validation_error_handle_func=$(FIBER_OPTS_VALIDATION_ERROR_FUNC)

# FASTHTTP OPTIONS
FASTHTTP_OPTS_PATHS                    ?= source_relative
FASTHTTP_OPTS_ERROR_HANDLERS_PACKAGE   ?= github.com/xakepp35/pkg/xprotoc-gen/utils
FASTHTTP_OPTS_JSON_UNMARSHAL_PACKAGE   ?= encoding/json
FASTHTTP_OPTS_GRPC_ERROR_FUNC          ?= FastHTTPHandleGRPCStatusError
FASTHTTP_OPTS_UNMARSHAL_ERROR_FUNC     ?= FastHttpHandleUnmarshalError
FASTHTTP_OPTS_VALIDATION_ERROR_FUNC    ?= FastHttpHandleValidationError

PROTOC_GEN_GO_FASTHTTP_OPTS = \
  paths=$(FASTHTTP_OPTS_PATHS) \
  error_handlers_package=$(FASTHTTP_OPTS_ERROR_HANDLERS_PACKAGE) \
  json_unmarshal_package=$(FASTHTTP_OPTS_JSON_UNMARSHAL_PACKAGE) \
  grpc_error_handle_func=$(FASTHTTP_OPTS_GRPC_ERROR_FUNC) \
  unmarshal_error_handle_func=$(FASTHTTP_OPTS_UNMARSHAL_ERROR_FUNC) \
  validation_error_handle_func=$(FASTHTTP_OPTS_VALIDATION_ERROR_FUNC)


OPENAPI_V2_OPT_OUTPUT_FORMAT ?= yaml
OPENAPI_V2_OPT_DISABLE_DEFAULT_ERRORS ?= true
OPENAPI_V2_OPT_ENUMS_AS_INTS ?= false

OPENAPI_V2_OPT = \
	output_format=$(OPENAPI_V2_OPT_OUTPUT_FORMAT) \
	disable_default_errors=$(OPENAPI_V2_OPT_DISABLE_DEFAULT_ERRORS) \
	enums_as_ints=$(OPENAPI_V2_OPT_ENUMS_AS_INTS)

install_tools:
	cd tools/protoc-gen-go-fiber && go build -o $(OUT)/protoc-gen-go-fiber .
	cd tools/protoc-gen-go-fasthttp  && go build -o $(OUT)/protoc-gen-go-fasthttp .
	cd tools/protoc-gen-ts-axios && go build -o $(OUT)/protoc-gen-ts-axios .
	GOBIN=$(OUT) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1
	GOBIN=$(OUT) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(OUT) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(OUT) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.3
	GOBIN=$(OUT) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3
.PHONY: install_tools

proto_gen:
	protoc -I $(INPUT_PROTO) -I lib \
             --go_out=$(GO_OUT) --go_opt=paths=source_relative \
             --go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
             --grpc-gateway_out=$(GO_OUT) --grpc-gateway_opt=paths=source_relative \
             --openapiv2_out=$(GO_OUT) \
             $(foreach opt,$(OPENAPI_V2_OPT),--openapiv2_opt=$(opt)) \
             --validate_out=$(GO_OUT) --validate_opt=paths=source_relative,lang=go \
             --go-fiber_out=$(GO_OUT) \
             $(foreach opt,$(PROTOC_GEN_GO_FIBER_OPTS),--go-fiber_opt=$(opt)) \
			 --go-fasthttp_out=$(GO_OUT) \
			 $(foreach opt,$(PROTOC_GEN_GO_FASTHTTP_OPTS),--go-fasthttp_opt=$(opt)) \
             $(PROTO_FILES)
.PHONY: proto_gen

proto_download:
	rm -rf lib

	# Создание каталогов
	mkdir -p lib/google/api
	mkdir -p lib/google/protobuf
	mkdir -p lib/pkg/types
	mkdir -p lib/protoc-gen-openapiv2/options
	mkdir -p lib/validate

	# Google API
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o lib/google/api/annotations.proto
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o lib/google/api/http.proto
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto -o lib/google/api/httpbody.proto

	# Google Protobuf (core descriptors)
	curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/empty.proto -o lib/google/protobuf/empty.proto
	curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/timestamp.proto -o lib/google/protobuf/timestamp.proto
	curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto -o lib/google/protobuf/descriptor.proto

	# AnyKey types
	curl -sSL https://raw.githubusercontent.com/xakepp35/pkg/refs/heads/main/src/pkg/types/empty.proto -o lib/pkg/types/empty.proto
	curl -sSL https://raw.githubusercontent.com/xakepp35/pkg/refs/heads/main/src/pkg/types/struct.proto -o lib/pkg/types/struct.proto
	curl -sSL https://raw.githubusercontent.com/xakepp35/pkg/refs/heads/main/src/pkg/types/time.proto -o lib/pkg/types/time.proto

	# OpenAPI v2 options
	curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/annotations.proto -o lib/protoc-gen-openapiv2/options/annotations.proto
	curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/openapiv2.proto -o lib/protoc-gen-openapiv2/options/openapiv2.proto

	# Validate
	curl -sSL https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto -o lib/validate/validate.proto
.PHONY: proto_download

install_protoc:
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/$(PROTOC_VERSION)/$(PROTOC_ZIP)

	unzip -o $(PROTOC_ZIP) -d /usr/local bin/protoc
	unzip -o $(PROTOC_ZIP) -d /usr/local 'include/*'

	rm -f $(PROTOC_ZIP)

	protoc --version
.PHONY: install_protoc

prepare: install_tools install_protoc proto_download
.PHONY: prepare

generate: proto_gen
.PHONY: proto_download

versions:
	@echo "================= Versions ================="
	@echo "protoc          : $(shell protoc --version)"
	@echo "go              : $(shell go version)"
.PHONY: versions
