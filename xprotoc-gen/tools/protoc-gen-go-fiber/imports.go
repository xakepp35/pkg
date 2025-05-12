package main

import "google.golang.org/protobuf/compiler/protogen"

var (
	fiberImport         = protogen.GoImportPath("github.com/gofiber/fiber/v2")
	contextImport       = protogen.GoImportPath("context")
	grpcMetadataImport  = protogen.GoImportPath("google.golang.org/grpc/metadata")
	grpcImport          = protogen.GoImportPath("google.golang.org/grpc")
	errorHandlersImport = protogen.GoImportPath(defaultFlagErrorHandlersPackage)
	jsonUnmarshalImport = protogen.GoImportPath(defaultJsonUnmarshalPackage)
)
