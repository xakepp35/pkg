package main

import (
	"fmt"
	"regexp"
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

	g.P("ctx, cancel := ", contextImport.Ident("WithCancel"), "(c.UserContext())")
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
			g.P("if err := Unmarshal(c.Body(), &req); err != nil {")
			g.P("\treturn HandleUnmarshalError(c, err)")
			g.P("}")
			g.P()
		}

		err := genReadReqFromQueryOrParams(g, method.Input, httpMethod, httpPath)
		if err != nil {
			return err
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
			g.P("\treturn HandleValidationError(c, err)")
			g.P("}")
			g.P()
		}
	}

	return nil
}

func genReadReqFromQueryOrParams(g *protogen.GeneratedFile, message *protogen.Message, method, path string) error {
	// 1. Check: all path fields must not be optional
	for _, field := range message.Fields {
		protoName := field.Desc.TextName()
		if strings.Contains(path, ":"+protoName) && field.Desc.HasPresence() {
			return xerrors.Err(nil).Str("field", field.GoName).Msg("optional field in path params is not allowed")
		}
	}

	for _, field := range message.Fields {
		fieldName := field.GoName
		protoName := field.Desc.TextName()

		var inPath bool

		// determine the source (Query or Param)
		var accessor string
		if strings.Contains(path, ":"+protoName) {
			accessor = `c.Params("` + protoName + `")`
			inPath = true
		} else {
			accessor = `c.Query("` + protoName + `")`
		}

		if !inPath && method != "Get" {
			continue
		}

		// determine the parser
		var parserFunc string
		var parserType string
		switch field.Desc.Kind() {
		case protoreflect.StringKind:
			parserFunc = "String"
			parserType = "string"
			if field.Desc.IsList() {
				accessor = fmt.Sprintf(`%s(%s, ",")`, g.QualifiedGoIdent(stringsImport.Ident("Split")), accessor)
			}
			g.P("req.", fieldName, " = ", accessor)
			g.P()
			continue
		case protoreflect.BoolKind:
			parserFunc = "Bool"
			parserType = "bool"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			parserFunc = "Int32"
			parserType = "int32"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			parserFunc = "Int64"
			parserType = "int64"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			parserFunc = "Uint32"
			parserType = "uint32"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			parserFunc = "Uint64"
			parserType = "uint64"
		case protoreflect.FloatKind:
			parserFunc = "Float32"
			parserType = "float32"
		case protoreflect.DoubleKind:
			parserFunc = "Float64"
			parserType = "float64"
		case protoreflect.BytesKind:
			parserFunc = "Bytes"
			parserType = "[]byte"
		default:
			g.P("// unsupported type for ", fieldName)
			continue
		}

		parseExpression := fmt.Sprintf("Parse%s(%s)", parserFunc, accessor)

		switch {
		case field.Desc.HasPresence():
			// Для optional query: если пусто — пропускаем парсинг
			if !inPath {
				g.P("if v := ", accessor, "; v != \"\" {")
				g.P("  req.", fieldName, ", err = ", g.QualifiedGoIdent(parsersImport.Ident("FirstArgPtr")), "(", fmt.Sprintf("Parse%s(v)", parserFunc), ")")
				g.P("  if err != nil {")
				g.P(`    return HandleGRPCStatusError(c, `,
					errorsBuilderImport.Ident("Err"), `(err).Str("field", "`, protoName, `").MsgProto(`, protoCodesImport.Ident("InvalidArgument"), `, "parse query/params field failed"))`)
				g.P("  }")
				g.P("}")
				g.P()
				continue
			}
			// Для path optional — мы уже выше вернули ошибку, сюда не попадём
			parseExpression = fmt.Sprintf("%s(%s)", g.QualifiedGoIdent(parsersImport.Ident("FirstArgPtr")), parseExpression)
		case field.Desc.IsList():
			if inPath {
				return xerrors.Err(nil).Str("field", fieldName).Msg("repeated field in params")
			}
			parseExpression = fmt.Sprintf("%s[%s](%s, Parse%s)", g.QualifiedGoIdent(parsersImport.Ident("ParseRepeated")), parserType, accessor, parserFunc)
		}

		g.P("req.", fieldName, ", err = ", parseExpression)
		g.P("if err != nil {")
		g.P(`  return HandleGRPCStatusError(c, `,
			errorsBuilderImport.Ident("Err"), `(err).Str("field", "`, protoName, `").MsgProto(`, protoCodesImport.Ident("InvalidArgument"), `, "parse query/params field failed"))`)
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

	g.P("	if err != nil { return HandleGRPCStatusError(c, err) }\n")

	httpMethod, _ := httpMethodParamsFromGrpcMethod(method)

	if httpMethod == "Get" && method.Output.GoIdent.GoName == "HttpBody" {
		// Parse HTTP body response
		g.P("    httpResp, ok := resp.(*", httpbodyImport.Ident("HttpBody"), ")")
		g.P("    if !ok || httpResp == nil {")
		g.P("        return HandleGRPCStatusError(c, ",
			errorsBuilderImport.Ident("Err"), "(nil).MsgProto(", protoCodesImport.Ident("Internal"), ", \"invalid http response\"))")
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

// httpMethodParamsFromGrpcMethod detects the HTTP method from the google.api.http annotation
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

	path = regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(path, `:$1`)

	return methodType, path
}
