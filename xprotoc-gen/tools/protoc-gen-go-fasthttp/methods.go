package main

import (
	"fmt"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func genMethod(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("func (r *", serviceRouterStructName(method.Parent), ")", genRouteMethodName(method), `(c *`, fasthttpImport.Ident("RequestCtx"), `) {`)

	g.P("ctx, cancel := ", contextImport.Ident("WithCancel"), "(context.Background())")
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
	g.P()
}

func genMethodReqPart(g *protogen.GeneratedFile, method *protogen.Method) {
	readAll := protogen.GoIdent{GoName: "ReadAll", GoImportPath: "io"}

	g.P("var req ", method.Input.GoIdent)
	g.P("var err error")
	g.P()

	stringsImportToUpper := protogen.GoIdent{
		GoImportPath: "strings",
		GoName:       "ToUpper",
	}
	stringsImportHasPrefix := protogen.GoIdent{
		GoImportPath: "strings",
		GoName:       "HasPrefix",
	}
	fmtImport := protogen.GoIdent{
		GoImportPath: "fmt",
		GoName:       "Sprint",
	}

	g.P("httpMethod := ", g.QualifiedGoIdent(stringsImportToUpper), "(string(c.Method()))")
	g.P("if httpMethod == \"GET\" || httpMethod == \"DELETE\" {")
	// Сначала path-параметры
	for _, field := range method.Input.Fields {
		if field.Desc.Kind() == protoreflect.StringKind {
			g.P(fmt.Sprintf("    if vPath := c.UserValue(%q); vPath != nil {", field.Desc.Name()))
			g.P("        req.", field.GoName, " = ", g.QualifiedGoIdent(fmtImport), "(vPath)")
			g.P("    }")
		}
	}
	// Потом query-параметры (перезаписывают path)
	for _, field := range method.Input.Fields {
		if field.Desc.Kind() == protoreflect.StringKind {
			g.P(fmt.Sprintf("    if vQuery := c.QueryArgs().Peek(%q); len(vQuery) > 0 {", field.Desc.Name()))
			g.P(fmt.Sprintf("        req.%s = string(vQuery)", field.GoName))
			g.P("    }")
		}
	}
	g.P("} else {")

	var bytesField *protogen.Field
	for _, field := range method.Input.Fields {
		if field.Desc.Kind() == protoreflect.BytesKind {
			bytesField = field
			break
		}
	}

	if bytesField != nil {
		// имя поля в Go и в форме
		fieldNameGo := bytesField.GoName
		fieldNameProto := bytesField.Desc.Name()

		g.P("contentType := string(c.Request.Header.ContentType())")
		g.P("if ", g.QualifiedGoIdent(stringsImportHasPrefix), "(contentType, \"multipart/form-data\") {")
		g.P(fmt.Sprintf("    if fh, errForm := c.FormFile(%q); errForm == nil {", fieldNameProto))
		g.P("        if f, errOpen := fh.Open(); errOpen == nil {")
		g.P("            defer f.Close()")
		g.P("            if data, errRead := ", g.QualifiedGoIdent(readAll), "(f); errRead == nil {")
		g.P(fmt.Sprintf("                req.%s = data", fieldNameGo))
		g.P("            }")
		g.P("        }")
		g.P("    }")

		// Дополнительно парсим остальные string-поля как FormValue
		for _, field := range method.Input.Fields {
			if field.Desc.Kind() == protoreflect.StringKind {
				g.P(fmt.Sprintf("    if v := c.FormValue(%q); len(v) > 0 {", field.Desc.Name()))
				g.P(fmt.Sprintf("        req.%s = string(v)", field.GoName))
				g.P("    }")
			}
		}

		g.P("} else {")
		g.P("    if err = ", jsonUnmarshalImport.Ident("Unmarshal"), "(c.PostBody(), &req); err != nil {")
		g.P("        ", errorHandlersImport.Ident(*flagUnmarshalErrorHandleFunc), "(c, err)")
		g.P("        return")
		g.P("    }")
		g.P("}")
	} else {
		g.P("if err = ", jsonUnmarshalImport.Ident("Unmarshal"), "(c.PostBody(), &req); err != nil {")
		g.P("    ", errorHandlersImport.Ident(*flagUnmarshalErrorHandleFunc), "(c, err)")
		g.P("    return")
		g.P("}")
	}

	g.P("}")

	hasValidation := false
	for _, field := range method.Input.Fields {
		if proto.HasExtension(field.Desc.Options(), validate.E_Rules) {
			hasValidation = true
			break
		}
	}
	if hasValidation {
		g.P("if err = req.Validate(); err != nil {")
		g.P("    ", errorHandlersImport.Ident(*flagValidationErrorHandleFunc), "(c, err)")
		g.P("    return")
		g.P("}")
		g.P()
	}
}

func genMethodExecPart(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("var (")
	g.P("resp any")
	g.P(")")

	g.P("if r.interceptor != nil {")
	g.P("handler := func(c context.Context, req any) (any, error) {")
	g.P("return r.server.", method.GoName, "(c, req.(*", method.Input.GoIdent, "))")
	g.P("}")
	g.P("info := &", grpcImport.Ident("UnaryServerInfo"), "{")
	g.P("Server: r.server,")
	g.P(fmt.Sprintf(`FullMethod: %s_%s_FullMethodName,`, method.Parent.GoName, method.GoName))
	g.P("}")
	g.P("resp, err = r.interceptor(ctx, &req, info, handler)")
	g.P("} else {")
	g.P("resp, err = r.server.", method.GoName, "(ctx, &req)")
	g.P("}")

	g.P("if err != nil {")
	g.P(errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, err)")
	g.P("return")
	g.P("}")
	g.P()

	// --- Всегда проверяем HttpBody ---
	httpBodyImport := g.QualifiedGoIdent(
		protogen.GoIdent{
			GoImportPath: "google.golang.org/genproto/googleapis/api/httpbody",
			GoName:       "HttpBody",
		},
	)

	structpbImport := g.QualifiedGoIdent(
		protogen.GoIdent{
			GoImportPath: "google.golang.org/protobuf/types/known/structpb",
			GoName:       "Struct",
		},
	)

	g.P("if hb, ok := resp.(*", httpBodyImport, "); ok {")
	g.P("if ct := hb.GetContentType(); ct != \"\" { c.SetContentType(ct) }")
	g.P("for _, ext := range hb.GetExtensions() {")
	g.P("    var s ", structpbImport)
	g.P("    if err := ext.UnmarshalTo(&s); err == nil {")
	g.P("        for k, v := range s.Fields {")
	g.P("            if v.GetKind() != nil {")
	g.P("                c.Response.Header.Set(k, v.GetStringValue())")
	g.P("            }")
	g.P("        }")
	g.P("    }")
	g.P("}")
	g.P("c.SetStatusCode(", fasthttpImport.Ident("StatusOK"), ")")
	g.P("c.Write(hb.GetData())")
	g.P("return")
	g.P("}")

	g.P("if dataBytes, ok := resp.([]byte); ok {")
	g.P("    c.SetContentType(\"application/octet-stream\")")
	g.P("    c.SetStatusCode(", fasthttpImport.Ident("StatusOK"), ")")
	g.P("    c.Write(dataBytes)")
	g.P("    return")
	g.P("}")

	g.P("data, mErr := ", jsonUnmarshalImport.Ident("Marshal"), "(resp)")
	g.P("if mErr != nil {")
	g.P(errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, mErr)")
	g.P("return")
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

	g.P("	r.", methodType, `(`, httpPath, `, h.`, genRouteMethodName(method), `)`)
}

func grpcOptionToMethodAndPathString(opts *descriptorpb.MethodOptions) (string, string) {
	ext := proto.GetExtension(opts, annotations.E_Http)
	var methodType, path string

	if httpRule, ok := ext.(*annotations.HttpRule); ok {
		switch pattern := httpRule.Pattern.(type) {
		case *annotations.HttpRule_Get:
			methodType, path = "GET", pattern.Get
		case *annotations.HttpRule_Post:
			methodType, path = "POST", pattern.Post
		case *annotations.HttpRule_Put:
			methodType, path = "PUT", pattern.Put
		case *annotations.HttpRule_Patch:
			methodType, path = "PATCH", pattern.Patch
		case *annotations.HttpRule_Delete:
			methodType, path = "DELETE", pattern.Delete
		default:
			methodType, path = "POST", "/"
		}
	}
	return methodType, path
}
