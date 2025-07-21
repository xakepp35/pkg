GOBIN=$(PWD)/bin
PATH := $(PATH):$(GOBIN)

generate_proto:
	protoc \
		-I lib \
		-I src \
    	--go_out=src --go_opt=paths=source_relative \
        $(shell find src -name "*.proto")
.PHONY: generate_proto

generate_mocks:
	mockery
.PHONY: generate_mocks

generate: generate_proto generate_mocks
.PHONY: generate

prepare: install_tools
.PHONY: prepare

install_tools:
	GOBIN=$(GOBIN) go install github.com/vektra/mockery/v3@v3.0.2
	GOBIN=$(OUT) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
.PHONY: install_tools

FIBER_TOOL_VERSION:=0.0.7
tools_download:
	curl -L https://github.com/petara94/protoc-gen-go-fiber/archive/refs/tags/v$(FIBER_TOOL_VERSION).zip -o repo.zip
	rm -rf tools/protoc-gen-go-fiber
	cd xprotoc-gen/tools && unzip ../../repo.zip
	rm repo.zip

	rm -rf xprotoc-gen/tools/protoc-gen-go-fiber
	mv xprotoc-gen/tools/protoc-gen-go-fiber-$(FIBER_TOOL_VERSION) xprotoc-gen/tools/protoc-gen-go-fiber
	rm -rf xprotoc-gen/tools/protoc-gen-go-fiber/.git
	rm -rf xprotoc-gen/tools/protoc-gen-go-fiber/example
	rm -rf xprotoc-gen/utils
	mv xprotoc-gen/tools/protoc-gen-go-fiber/utils xprotoc-gen/

.PHONY: tools_download


INPUT_PROTO=proto
API_OUT=api-go
GO_OUT=$(PWD)/$(API_OUT)

XPROTOC_GEN_VERSION:=v0.0.1

docker_build:
	docker build -t xprotoc-gen:$(XPROTOC_GEN_VERSION) ./xprotoc-gen
.PHONY: docker_build

copy_lib_from_docker:
	docker create --name temp-xprotoc-gen xprotoc-gen
	rm -rf lib
	docker cp temp-xprotoc-gen:/app/lib lib
	docker rm temp-xprotoc-gen
.PHONY: copy_lib_from_docker

docker_generate:
	docker run --rm -v $(PWD)/$(INPUT_PROTO):/app/$(INPUT_PROTO) -v $(GO_OUT):/app/$(API_OUT) xprotoc-gen bash -c "make generate"
.PHONY: docker_generate
