package main

import (
	"fmt"
	"github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"slices"
	"unicode"
)

func genMethod(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("func (r *", serviceRouterStructName(method.Parent), ")", genRouteMethodName(method), `(c *`, fiberImport.Ident("Ctx"), `) error {`)

	g.P("ctx, cancel := ", contextImport.Ident("WithCancel"), "(c.Context())")
	g.P("defer cancel()\n")

	g.P("md := ", grpcMetadataImport.Ident("New"), "(nil)")
	g.P("c.Request().Header.VisitAll(func(key, value []byte) {")
	g.P("md.Append(string(key), string(value))")
	g.P("})")
	g.P()

	g.P("ctx = metadata.NewIncomingContext(ctx, md)")
	g.P()

	genMethodReqPart(g, method)

	genMethodExecPart(g, method)

	g.P("	}")

}

func genMethodReqPart(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("var req ", method.Input.GoIdent)
	g.P()

	hasExportedField := slices.ContainsFunc(method.Input.Fields, func(f *protogen.Field) bool {
		return unicode.IsUpper(rune(f.GoName[0]))
	})

	// use marshaller if we need
	if hasExportedField {
		g.P("if err := ", jsonUnmarshalImport.Ident("Unmarshal"), "(c.Body(), &req); err != nil {")
		g.P("	return ", errorHandlersImport.Ident(*flagUnmarshalErrorHandleFunc), "(c, err)")
		g.P("}")
		g.P()

		hasValidation := false
		for _, field := range method.Input.Fields {
			if proto.HasExtension(field.Desc.Options(), validate.E_Rules) {
				hasValidation = true
				break
			}
		}

		if hasValidation {
			g.P("if err := req.Validate(); err != nil {")
			g.P("	return ", errorHandlersImport.Ident(*flagValidationErrorHandleFunc), "(c, err)")
			g.P("}")
			g.P()
		}
	}
}

func genMethodExecPart(g *protogen.GeneratedFile, method *protogen.Method) {
	if !*flagDisableGrpcInterceptor {
		g.P("handler := func(ctx context.Context, req any) (any, error) {")
		g.P("return r.server.", method.GoName, "(ctx, req.(*", method.Input.GoIdent, "))")
		g.P("}")
		g.P("info := &", grpcImport.Ident("UnaryServerInfo"), "{")
		g.P("Server: r.server,")
		g.P(fmt.Sprintf(`FullMethod: %s_%s_FullMethodName,`, method.Parent.GoName, method.GoName))
		g.P("}")
		g.P()
		g.P("resp, err := r.interceptor(ctx, &req, info, handler)")
	} else {
		g.P("resp, err := r.server.", method.GoName, "(ctx, &req)")
	}

	g.P("	if err != nil { return ", errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, err) }\n")

	g.P("	return c.JSON(resp)")
}

func genRouteMethodName(method *protogen.Method) string {
	return fmt.Sprintf("__%s_%s_Route", method.Parent.GoName, method.GoName)
}

func genFiberMethodRote(g *protogen.GeneratedFile, method *protogen.Method) {
	opts := method.Desc.Options().(*descriptorpb.MethodOptions)

	methodType, httpPath := grpcOptionToMethodAndPathString(opts)
	if httpPath == "/" {
		httpPath += string(method.Parent.Desc.FullName()) + "/" + string(method.Desc.Name())
	}

	g.P("	app.", methodType, `("`, httpPath, `", router.`, genRouteMethodName(method), `)`)
}

// grpcOptionToMethodAndPathString узнает метод из google.api.http
func grpcOptionToMethodAndPathString(opts *descriptorpb.MethodOptions) (string, string) {
	ext := proto.GetExtension(opts, annotations.E_Http)
	var methodType, path string

	if httpRule, ok := ext.(*annotations.HttpRule); ok {

		switch pattern := httpRule.Pattern.(type) {
		case *annotations.HttpRule_Get:
			methodType = "Get"
			path = pattern.Get
		case *annotations.HttpRule_Post:
			methodType = "Post"
			path = pattern.Post
		case *annotations.HttpRule_Put:
			methodType = "Put"
			path = pattern.Put
		case *annotations.HttpRule_Patch:
			methodType = "Patch"
			path = pattern.Patch
		case *annotations.HttpRule_Delete:
			methodType = "Delete"
			path = pattern.Delete
		default:
			// fallback
			methodType = "Post"
			path = "/"
		}
	}
	return methodType, path
}
