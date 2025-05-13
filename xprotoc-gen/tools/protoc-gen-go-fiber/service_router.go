package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

func genServiceRouter(g *protogen.GeneratedFile, service *protogen.Service) {
	g.P("type ", serviceRouterStructName(service), " struct {")
	g.P("server ", service.GoName, "Server")
	g.P("interceptor ", grpcImport.Ident("UnaryServerInterceptor"))
	g.P("}")
	g.P()
}

func genServiceRouterDeclaration(g *protogen.GeneratedFile, service *protogen.Service) {
	g.P("router := &", serviceRouterStructName(service), "{")
	g.P("server: server,")
	g.P("interceptor: interceptor,")
	g.P("}")
	g.P()
}

func serviceRouterStructName(service *protogen.Service) string {
	return fmt.Sprintf("__%s_FiberRouter", service.GoName)
}
