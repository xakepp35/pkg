package main

import (
	"fmt"
	"slices"
	"unicode"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func genMethod(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("func (r *", serviceRouterStructName(method.Parent), ")", genRouteMethodName(method), `(c *`, fasthttpImport.Ident("RequestCtx"), `) error {`)

	g.P("ctx, cancel := ", contextImport.Ident("WithCancel"), "(c)")
	g.P("defer cancel()\n")

	g.P("md := ", grpcMetadataImport.Ident("New"), "(nil)")
	g.P("c.Request.Header.VisitAll(func(key, value []byte) {")
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
		g.P("if err := ", jsonUnmarshalImport.Ident("Unmarshal"), "(c.PostBody(), &req); err != nil {")
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
	g.P("var (")
	g.P("resp any")
	g.P("err  error")
	g.P(")")

	g.P("if r.interceptor != nil {")
	g.P("handler := func(c context.Context, req any) (any, error) {")
	g.P("return r.server.", method.GoName, "(c, req.(*", method.Input.GoIdent, "))")
	g.P("}")
	g.P("info := &", grpcImport.Ident("UnaryServerInfo"), "{")
	g.P("Server: r.server,")
	g.P(fmt.Sprintf(`FullMethod: %s_%s_FullMethodName,`, method.Parent.GoName, method.GoName))
	g.P("}")
	g.P("resp, err = r.interceptor(c, &req, info, handler)")
	g.P("} else {")
	g.P("resp, err = r.server.", method.GoName, "(c, &req)")
	g.P("}")

	g.P("if err != nil {")
	g.P(errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, err)")
	g.P("return err")
	g.P("}")
	g.P()

	g.P("data, mErr := ", jsonUnmarshalImport.Ident("Marshal"), "(resp)")
	g.P("if mErr != nil {")
	g.P(errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, mErr)")
	g.P("return mErr")
	g.P("}")
	g.P("c.SetContentType(\"application/json\")")
	g.P("c.SetStatusCode(", fasthttpImport.Ident("StatusOK"), ")")
	g.P("c.Write(data)")
}

func genRouteMethodName(method *protogen.Method) string {
	return fmt.Sprintf("__%s_%s_Route", method.Parent.GoName, method.GoName)
}

func genFastHTTPMethodRoute(g *protogen.GeneratedFile, method *protogen.Method) {
	opts := method.Desc.Options().(*descriptorpb.MethodOptions)

	methodType, httpPath := grpcOptionToMethodAndPathString(opts)
	if httpPath == "/" {
		httpPath = fmt.Sprintf(`%s_%s_FullMethodName`, method.Parent.GoName, method.GoName)
	} else {
		httpPath = `"` + httpPath + `"`
	}

	g.P("	r.", methodType, `(`, httpPath, `, router.`, genRouteMethodName(method), `)`)
}

func grpcOptionToMethodAndPathString(opts *descriptorpb.MethodOptions) (string, string) {
	ext := proto.GetExtension(opts, annotations.E_Http)
	var methodType, path string

	if httpRule, ok := ext.(*annotations.HttpRule); ok {
		switch pattern := httpRule.Pattern.(type) {
		case *annotations.HttpRule_Get:
			methodType, path = "Get", pattern.Get
		case *annotations.HttpRule_Post:
			methodType, path = "Post", pattern.Post
		case *annotations.HttpRule_Put:
			methodType, path = "Put", pattern.Put
		case *annotations.HttpRule_Patch:
			methodType, path = "Patch", pattern.Patch
		case *annotations.HttpRule_Delete:
			methodType, path = "Delete", pattern.Delete
		default:
			methodType, path = "Post", "/"
		}
	}
	return methodType, path
}
