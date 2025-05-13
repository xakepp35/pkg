package main

import (
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	defaultFlagErrorHandlersPackage = "github.com/petara94/protoc-gen-go-fiber/utils"
	defaultJsonUnmarshalPackage     = "encoding/json"
)

var (
	flags                         = flag.NewFlagSet("protoc-gen-go-fiber", flag.ExitOnError)
	flagErrorHandlersPackage      = flags.String("error_handlers_package", defaultFlagErrorHandlersPackage, "package with error handlers funcs")
	flagJsonUnmarshalPackage      = flags.String("json_unmarshal_package", defaultJsonUnmarshalPackage, "package with json unmarshalers")
	flagGrpcErrorHandleFunc       = flags.String("grpc_error_handle_func", "HandleGRPCStatusError", "func name for handle grpc error")
	flagUnmarshalErrorHandleFunc  = flags.String("unmarshal_error_handle_func", "HandleUnmarshalError", "func name for handle unmarshal error")
	flagValidationErrorHandleFunc = flags.String("validation_error_handle_func", "HandleValidationError", "func name for handle validation error")
)

func flagInit() {
	errorHandlersImport = protogen.GoImportPath(*flagErrorHandlersPackage)
	jsonUnmarshalImport = protogen.GoImportPath(*flagJsonUnmarshalPackage)
}
