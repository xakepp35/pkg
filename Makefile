GOBIN=$(PWD)/bin
PATH := $(PATH):$(GOBIN)

generate_proto:
	protoc \
		-I lib \
		-I types \
    	--go_out=types --go_opt=paths=source_relative \
        $(shell find types -name "*.proto") 
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
