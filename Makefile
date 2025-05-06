generate:
	protoc \
		-I types \
    	--go_out=types --go_opt=paths=source_relative \
        $(shell find types -name "*.proto") 
.PHONY: generate
