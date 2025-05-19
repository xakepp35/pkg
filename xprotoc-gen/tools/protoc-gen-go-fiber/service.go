package main

import "google.golang.org/protobuf/compiler/protogen"

func genService(g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	genServiceRouter(g, service)

	g.P("func Register", service.GoName, "FiberRoutes(app *", fiberImport.Ident("App"), ", server ", service.GoName, "Server, interceptor grpc.UnaryServerInterceptor) {")

	g.P("if server == nil {")
	g.P("panic(\"register fiber router filed: server ", service.GoName, "Server is nil\")")
	g.P("}")
	g.P()

	genServiceRouterDeclaration(g, service)

	for _, method := range service.Methods {
		genFiberMethodRote(g, method)
	}
	g.P("}")
	g.P()

	for _, method := range service.Methods {
		genMethod(g, method)
	}
}
