package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	"github.com/xakepp35/pkg/xerrors"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func genMethod(g *protogen.GeneratedFile, method *protogen.Method) error {
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

	err := genMethodReqPart(g, method)
	if err != nil {
		return err
	}

	genMethodExecPart(g, method)

	g.P("	}")
	g.P()

	return nil
}

func genMethodReqPart(g *protogen.GeneratedFile, method *protogen.Method) error {
	g.P("var (")
	g.P("req ", method.Input.GoIdent)
	g.P("resp any")
	g.P("err error")
	g.P(")")
	g.P()

	hasExportedField := slices.ContainsFunc(method.Input.Fields, func(f *protogen.Field) bool {
		return unicode.IsUpper(rune(f.GoName[0]))
	})

	// use marshaller if we need
	if hasExportedField {
		httpMethod, httpPath := httpMethodParamsFromGrpcMethod(method)

		if httpMethod != "Get" {
			g.P("if err := ", jsonUnmarshalImport.Ident("Unmarshal"), "(c.Body(), &req); err != nil {")
			g.P("	return ", errorHandlersImport.Ident(*flagUnmarshalErrorHandleFunc), "(c, err)")
			g.P("}")
			g.P()
		} else {
			err := genReadReqFromQueryOrParams(g, method.Input, httpPath)
			if err != nil {
				return err
			}
		}

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

	return nil
}

func genReadReqFromQueryOrParams(g *protogen.GeneratedFile, message *protogen.Message, path string) error {
	for _, field := range message.Fields {
		fieldName := field.GoName
		protoName := field.Desc.TextName()

		var inPath bool

		// определение источника (Query или Param)
		var accessor string
		if strings.Contains(path, ":"+protoName) {
			accessor = `c.Params("` + protoName + `")`
			inPath = true
		} else {
			accessor = `c.Query("` + protoName + `")`
		}

		// определение парсера
		var parserFunc string
		switch field.Desc.Kind() {
		case protoreflect.StringKind:
			parserFunc = "ParseString"

			if field.Desc.IsList() {
				accessor = fmt.Sprintf(`%s(%s, ",")`, g.QualifiedGoIdent(stringsImport.Ident("Split")), accessor)
			}

			g.P("req.", fieldName, " = ", accessor)
			g.P()
			continue
		case protoreflect.BoolKind:
			parserFunc = "ParseBool"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			parserFunc = "ParseInt32"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			parserFunc = "ParseInt64"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			parserFunc = "ParseUint32"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			parserFunc = "ParseUint64"
		case protoreflect.FloatKind:
			parserFunc = "ParseFloat32"
		case protoreflect.DoubleKind:
			parserFunc = "ParseFloat64"
		case protoreflect.BytesKind:
			parserFunc = "ParseBytes"
		default:
			g.P("// unsupported type for ", fieldName)
			continue
		}

		parseExpression := fmt.Sprintf("%s(%s)", g.QualifiedGoIdent(parsersImport.Ident(parserFunc)), accessor)

		switch {
		case field.Desc.HasPresence():
			parseExpression = fmt.Sprintf("%s(%s)", g.QualifiedGoIdent(parsersImport.Ident("FirstArgPtr")), parseExpression)
		case field.Desc.IsList():
			if inPath {
				return xerrors.Err(nil).Msg("repeated field in params").Str("field", fieldName).Err()
			}
			parseExpression = fmt.Sprintf("%s(%s, %s)",
				g.QualifiedGoIdent(parsersImport.Ident("ParseRepeated")),
				accessor, g.QualifiedGoIdent(parsersImport.Ident(parserFunc)),
			)
		}

		// генерация кода парсинга
		g.P("req.", fieldName, ", err = ", parseExpression)
		g.P("if err != nil {")
		g.P(`  return `, errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), `(c, `,
			errorsBuilderImport.Ident("Err"),
			`(err).Msg("parse query/params field failed").`)
		g.P(`Str("field", "`, protoName, `").`)
		g.P(`ProtoErr(`, protoCodesImport.Ident("InvalidArgument"), `))`)
		g.P("}")
		g.P()
	}

	return nil
}

func genMethodExecPart(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("if r.interceptor != nil {")
	g.P("handler := func(ctx context.Context, req any) (any, error) {")
	g.P("return r.server.", method.GoName, "(ctx, req.(*", method.Input.GoIdent, "))")
	g.P("}")
	g.P("info := &", grpcImport.Ident("UnaryServerInfo"), "{")
	g.P("Server: r.server,")
	g.P(fmt.Sprintf(`FullMethod: %s_%s_FullMethodName,`, method.Parent.GoName, method.GoName))
	g.P("}")
	g.P()
	g.P("resp, err = r.interceptor(ctx, &req, info, handler)")
	g.P("} else {")
	g.P("resp, err = r.server.", method.GoName, "(ctx, &req)")
	g.P("}")

	g.P("	if err != nil { return ", errorHandlersImport.Ident(*flagGrpcErrorHandleFunc), "(c, err) }\n")

	httpMethod, _ := httpMethodParamsFromGrpcMethod(method)

	if httpMethod == "Get" && method.Output.GoIdent.GoName == "HttpBody" {
		// Parse HTTP body response
		g.P("    httpResp, ok := resp.(*", httpbodyImport.Ident("HttpBody"), ")")
		g.P("    if !ok || httpResp == nil {")
		g.P("        return ", errorHandlersImport.Ident(*flagGrpcErrorHandleFunc),
			"(c, ",
			errorsBuilderImport.Ident("Err"), "(nil).",
			"Msg(\"invalid http response\").",
			"ProtoErr(", protoCodesImport.Ident("Internal"), "))",
		)
		g.P("    }")
		g.P("    c.Set(", fiberImport.Ident("HeaderContentType"), ", httpResp.GetContentType())")
		g.P("    return c.Status(", fiberImport.Ident("StatusOK"), ").Send(httpResp.GetData())")
	} else {
		g.P("	return c.JSON(resp)")
	}
}

func genRouteMethodName(method *protogen.Method) string {
	return fmt.Sprintf("__%s_%s_Route", method.Parent.GoName, method.GoName)
}

func genFiberMethodRote(g *protogen.GeneratedFile, method *protogen.Method) {
	methodType, httpPath := httpMethodParamsFromGrpcMethod(method)
	if httpPath == "/" {
		httpPath = fmt.Sprintf(`%s_%s_FullMethodName`, method.Parent.GoName, method.GoName)
	} else {
		httpPath = `"` + httpPath + `"`
	}

	g.P("	app.", methodType, `(`, httpPath, `, router.`, genRouteMethodName(method), `)`)
}

// httpMethodParamsFromGrpcMethod узнает метод из аннотации google.api.http
func httpMethodParamsFromGrpcMethod(method *protogen.Method) (string, string) {
	opts := method.Desc.Options().(*descriptorpb.MethodOptions)

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
