package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func genService(g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	// generate fasthttp service router struct
	genServiceRouter(g, service)

	// func Register<Service>FastHTTPRoutes(r *router.Router, server <Service>Server, interceptor grpc.UnaryServerInterceptor)
	g.P("func Register", service.GoName, "FastHTTPRoutes(",
		"r *", fasthttpRouterImport.Ident("Router"), ", ",
		"server ", service.GoName, "Server, ",
		"interceptor grpc.UnaryServerInterceptor) {")

	g.P("if server == nil {")
	g.P("panic(\"register fasthttp router failed: server ", service.GoName, "Server is nil\")")
	g.P("}")
	g.P()

	// creating router
	genServiceRouterDeclaration(g, service)

	// generating http handlers
	for _, m := range service.Methods {
		genFastHTTPMethodRoute(g, m)
	}

	g.P("}")
	g.P()

	for _, m := range service.Methods {
		genMethod(g, m)
	}
}
