package main

import "google.golang.org/protobuf/compiler/protogen"

var (
	fasthttpImport       = protogen.GoImportPath("github.com/valyala/fasthttp")
	fasthttpRouterImport = protogen.GoImportPath("github.com/fasthttp/router")
	contextImport        = protogen.GoImportPath("context")
	grpcMetadataImport   = protogen.GoImportPath("google.golang.org/grpc/metadata")
	grpcImport           = protogen.GoImportPath("google.golang.org/grpc")
	errorHandlersImport  = protogen.GoImportPath(defaultFlagErrorHandlersPackage)
	jsonUnmarshalImport  = protogen.GoImportPath(defaultJsonUnmarshalPackage)
)
