package main

import "google.golang.org/protobuf/compiler/protogen"

var (
	fiberImport         = protogen.GoImportPath("github.com/gofiber/fiber/v2")
	contextImport       = protogen.GoImportPath("context")
	stringsImport       = protogen.GoImportPath("strings")
	grpcMetadataImport  = protogen.GoImportPath("google.golang.org/grpc/metadata")
	grpcImport          = protogen.GoImportPath("google.golang.org/grpc")
	errorHandlersImport = protogen.GoImportPath(defaultFlagErrorHandlersPackage)
	errorsBuilderImport = protogen.GoImportPath("github.com/xakepp35/pkg/xerrors")
	protoCodesImport    = protogen.GoImportPath("google.golang.org/grpc/codes")
	parsersImport       = protogen.GoImportPath(defaultParsersPackage)
	jsonUnmarshalImport = protogen.GoImportPath(defaultJsonUnmarshalPackage)
	httpbodyImport      = protogen.GoImportPath("google.golang.org/genproto/googleapis/api/httpbody")
)
