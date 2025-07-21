package main

import (
	"github.com/xakepp35/pkg/xerrors"
	"google.golang.org/protobuf/compiler/protogen"
)

func genService(g *protogen.GeneratedFile, service *protogen.Service) error {
	if len(service.Methods) == 0 {
		return nil
	}

	genServiceRouter(g, service)

	g.P("func Register", service.GoName, "FiberRoutes(app *", fiberImport.Ident("App"), ", server ", service.GoName, "Server, interceptor grpc.UnaryServerInterceptor) {")

	g.P("if app == nil {")
	g.P("panic(\"register fiber router filed: fiber Server \\\"app\\\" is nil\")")
	g.P("}")
	g.P()
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
		err := genMethod(g, method)
		if err != nil {
			return xerrors.Err(err).Msg("method generation").Str("method", method.GoName).Err()
		}
	}

	return nil
}
