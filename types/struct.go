package types

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type Document = map[string]any

// NewStruct constructs a Struct from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue.
func NewStruct(v Document) *Struct {
	return &Struct{
		Fields: MapAsFields(v),
	}
}

func MapAsFields(v Document) map[string]*Value {
	x := make(map[string]*Value, len(v))
	for k, v := range v {
		// if !utf8.ValidString(k) {
		// 	return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", k)
		// }
		// var err error
		x[k] = NewValue(v)
		// if err != nil {
		// 	return nil, err
		// }
	}
	return x
}

// AsMap converts x to a general-purpose Go map.
// The map values are converted by calling Value.AsInterface.
func (x *Struct) AsMap() Document {
	return FieldsAsMap(x.GetFields())
}

func FieldsAsMap(req map[string]*Value) Document {
	if len(req) == 0 {
		return nil
	}
	res := make(Document, len(req))
	for k, v := range req {
		res[k] = v.AsInterface()
	}
	return res
}

func (x *Struct) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Fields)
}

func (x *Struct) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &x.Fields)
}

func (x *Struct) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(x.Fields)
}

func (x *Struct) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.Decode(&x.Fields)
}

// NewValue constructs a Value from a general-purpose Go interface.
//
//	╔════════════════════════╤════════════════════════════════════════════╗
//	║ Go type                │ Conversion                                 ║
//	╠════════════════════════╪════════════════════════════════════════════╣
//	║ nil                    │ stored as NullValue                        ║
//	║ bool                   │ stored as BoolValue                        ║
//	║ int, int32, int64      │ stored as NumberValue                      ║
//	║ uint, uint32, uint64   │ stored as NumberValue                      ║
//	║ float32, float64       │ stored as NumberValue                      ║
//	║ string                 │ stored as StringValue; must be valid UTF-8 ║
//	║ []byte                 │ stored as StringValue; base64-encoded      ║
//	║ Document               │ stored as StructValue                      ║
//	║ []any                  │ stored as ListValue                        ║
//	╚════════════════════════╧════════════════════════════════════════════╝
//
// When converting an int64 or uint64 to a NumberValue, numeric precision loss
// is possible since they are stored as a float64.
func NewValue(v any) *Value {
	return &Value{Kind: NewValueKind(v)}
}

func NewValueKind(v any) isValue_Kind {
	switch v := v.(type) {
	case nil:
		return NewNullValueKind()
	case bool:
		return NewBoolValueKind(v)
	case int:
		return NewIntValueKind(int64(v))
	case int8:
		return NewIntValueKind(int64(v))
	case int16:
		return NewIntValueKind(int64(v))
	case int32:
		return NewIntValueKind(int64(v))
	case int64:
		return NewIntValueKind(int64(v))
	case uint:
		return NewIntValueKind(int64(v))
	case uint8:
		return NewIntValueKind(int64(v))
	case uint16:
		return NewIntValueKind(int64(v))
	case uint32:
		return NewIntValueKind(int64(v))
	case uint64:
		return NewIntValueKind(int64(v))
	// case uintptr: // special case, panics
	// 	return NewIntValueKind(int64(v))
	case float32:
		return NewNumberValueKind(float64(v))
	case float64:
		return NewNumberValueKind(v)
	case string:
		// if !utf8.ValidString(v) {
		// 	return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", v)
		// }
		return NewStringValueKind(v)
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewStringValueKind(s)
	case Document:
		v2 := NewStruct(v)
		return NewStructValueKind(v2)
	case []any:
		v2 := NewList(v)
		return NewListValueKind(v2)
	case time.Time:
		v2 := NewTime(v)
		return NewTimeValueKind(v2)
	case *Time:
		return NewTimeValueKind(v)
	default:
		panic(protoimpl.X.NewError("framework.structpb.NewValue(%T): invalid type", v))
		// return nil
		// , protoimpl.X.NewError("invalid type: %T", v)
	}
}

// NewNullValue constructs a new null Value.
func NewNullValue() *Value {
	return &Value{Kind: NewNullValueKind()}
}

// NewNullValueKind constructs a new null Value.
func NewNullValueKind() *Value_NullValue {
	return &Value_NullValue{NullValue: Null_NULL_VALUE}
}

// NewBoolValue constructs a new boolean Value.
func NewBoolValue(v bool) *Value {
	return &Value{Kind: NewBoolValueKind(v)}
}

// NewBoolValueKind constructs a new boolean Value.
func NewBoolValueKind(v bool) *Value_BoolValue {
	return &Value_BoolValue{BoolValue: v}
}

// NewIntValue constructs a new int Value.
func NewIntValue(v int64) *Value {
	return &Value{Kind: NewIntValueKind(v)}
}

// NewIntValueKind constructs a new int Value.
func NewIntValueKind(v int64) *Value_IntValue {
	return &Value_IntValue{IntValue: v}
}

// NewNumberValue constructs a new number (floating point) Value.
func NewNumberValue(v float64) *Value {
	return &Value{Kind: NewNumberValueKind(v)}
}

// NewNumberValueKind constructs a new number (floating point) Value.
func NewNumberValueKind(v float64) *Value_NumberValue {
	return &Value_NumberValue{NumberValue: v}
}

// NewStringValue constructs a new string Value.
func NewStringValue(v string) *Value {
	return &Value{Kind: &Value_StringValue{StringValue: v}}
}

// NewStringValueKind constructs a new string Value.
func NewStringValueKind(v string) *Value_StringValue {
	return &Value_StringValue{StringValue: v}
}

// NewStructValue constructs a new struct Value.
func NewStructValue(v *Struct) *Value {
	return &Value{Kind: &Value_StructValue{StructValue: v}}
}

// NewStructValueKind constructs a new struct Value.
func NewStructValueKind(v *Struct) *Value_StructValue {
	return &Value_StructValue{StructValue: v}
}

// NewListValue constructs a new list Value.
func NewListValue(v *List) *Value {
	return &Value{Kind: NewListValueKind(v)}
}

// NewListValueKind constructs a new list Value.
func NewListValueKind(v *List) *Value_ListValue {
	return &Value_ListValue{ListValue: v}
}

// NewTimeValue constructs a new Time Value.
func NewTimeValue(v *Time) *Value {
	return &Value{Kind: NewTimeValueKind(v)}
}

// NewTimeValueKind constructs a new Time Value.
func NewTimeValueKind(v *Time) *Value_TimeValue {
	return &Value_TimeValue{TimeValue: v}
}

func NewStringList(v ...string) *List {
	list := make([]*Value, len(v))
	for i := range v {
		list[i] = NewStringValue(v[i])
	}
	return &List{List: list}
}

func NewStringListValue(v ...string) *Value {
	return NewListValue(NewStringList(v...))
}

func (x *Value) ToReflectField(value reflect.Value) {
	valueFace := value.Interface()
	switch valueFace.(type) {
	case string:
		value.SetString(x.GetStringValue())
	case bool:
		value.SetBool(x.GetBoolValue())
	case float64:
		value.SetFloat(x.GetNumberValue())
	case int:
		value.SetInt(int64(x.GetNumberValue()))
	case int64:
		value.SetInt(int64(x.GetNumberValue()))
	default:
		err := status.Errorf(codes.InvalidArgument, "Value.ToReflect(%T) not implemented", valueFace)
		panic(err.Error())
		// TODO.. extend.
	}
}

func (x *Value) Clear() {
	switch x.Kind.(type) {
	case *Value_NullValue:
		x.Kind = NewNullValueKind()
	case *Value_NumberValue:
		x.Kind = NewNumberValueKind(0)
	case *Value_StringValue:
		x.Kind = NewStringValueKind("")
	case *Value_BoolValue:
		x.Kind = NewBoolValueKind(false)
	case *Value_StructValue:
		x.Kind = NewStructValueKind(NewStruct(nil))
	case *Value_ListValue:
		x.Kind = NewListValueKind(NewList(nil))
	case *Value_TimeValue:
		x.Kind = NewTimeValueKind(NewTime(time.Time{}))
	case *Value_IntValue:
		x.Kind = NewIntValueKind(0)
	}
}

func (s *Value) TruncateSublist(maxLength int) {
	subList := s.GetListValue()
	if subList != nil && len(subList.List) > maxLength {
		subList.List = subList.List[:maxLength]
	}
}

func (x *Value) AsBool() bool {
	if x == nil {
		return false
	}
	switch t := x.Kind.(type) {
	case *Value_IntValue:
		return t.IntValue != 0
	case *Value_NumberValue:
		return t.NumberValue != 0
	case *Value_NullValue:
		return false
	case *Value_StringValue:
		switch strings.ToLower(t.StringValue) {
		case "1", "+", "true":
			return true
		case "0", "-", "false":
			return false
		default:
			return false // TODO: panic?
		}
	case *Value_BoolValue:
		return t.BoolValue
	case *Value_StructValue:
		return false
	case *Value_ListValue:
		return false
	case *Value_TimeValue:
		return t.TimeValue != nil
	default:
		return false
	}
}

func (x *Value) AsUint64() uint64 {
	return uint64(x.AsInt64())
}

func (x *Value) AsInt64() int64 {
	if x == nil {
		return 0
	}
	switch t := x.Kind.(type) {
	case *Value_IntValue:
		return t.IntValue
	case *Value_NumberValue:
		return int64(t.NumberValue)
	case *Value_NullValue:
		return 0
	case *Value_StringValue:
		res, _ := strconv.ParseInt(t.StringValue, 10, 64)
		return res
	case *Value_BoolValue:
		if t.BoolValue {
			return 1
		}
		return 0
	case *Value_StructValue:
		return 0
	case *Value_ListValue:
		return 0
	case *Value_TimeValue:
		return t.TimeValue.Seconds
	default:
		return 0
	}
}

func (x *Value) AsFloat64() float64 {
	if x == nil {
		return 0
	}
	switch t := x.Kind.(type) {
	case *Value_NullValue:
		return 0
	case *Value_NumberValue:
		return t.NumberValue
	case *Value_StringValue:
		res, _ := strconv.ParseFloat(t.StringValue, 64)
		return res
	case *Value_BoolValue:
		if t.BoolValue {
			return 1
		}
		return 0
	case *Value_StructValue:
		return 0
	case *Value_ListValue:
		return 0
	case *Value_TimeValue:
		return t.TimeValue.Unix()
	case *Value_IntValue:
		return float64(t.IntValue)
	default:
		return 0
	}
}

func (x *Value) AsString() string {
	if x == nil {
		return ""
	}
	switch t := x.Kind.(type) {
	case *Value_NullValue:
		return ""
	case *Value_NumberValue:
		return strconv.FormatFloat(t.NumberValue, 'f', 0, 64)
	case *Value_StringValue:
		return t.StringValue
	case *Value_BoolValue:
		if t.BoolValue {
			return "true"
		}
		return "false"
	case *Value_StructValue:
		return ""
	case *Value_ListValue:
		return ""
	case *Value_TimeValue:
		return t.TimeValue.AsTime().Format(time.RFC3339Nano)
	case *Value_IntValue:
		return strconv.FormatInt(t.IntValue, 10)
	default:
		return ""
	}
}

// AsInterface converts x to a general-purpose Go interface.
//
// Calling Value.MarshalJSON and "encoding/json".Marshal on this output produce
// semantically equivalent JSON (assuming no errors occur).
//
// Floating-point values (i.e., "NaN", "Infinity", and "-Infinity") are
// converted as strings to remain compatible with MarshalJSON.
func (x *Value) AsInterface() any {
	if x == nil {
		return nil
	}
	switch v := x.GetKind().(type) {
	case *Value_IntValue:
		return v.IntValue
	case *Value_TimeValue:
		return v.TimeValue
	case *Value_NumberValue:
		if v != nil {
			switch {
			case math.IsNaN(v.NumberValue):
				return "NaN"
			case math.IsInf(v.NumberValue, +1):
				return "Infinity"
			case math.IsInf(v.NumberValue, -1):
				return "-Infinity"
			default:
				return v.NumberValue
			}
		}
	case *Value_StringValue:
		if v != nil {
			return v.StringValue
		}
	case *Value_BoolValue:
		if v != nil {
			return v.BoolValue
		}
	case *Value_StructValue:
		if v != nil {
			return v.StructValue.AsMap()
		}
	case *Value_ListValue:
		if v != nil {
			return v.ListValue.AsSlice()
		}
	}
	return nil
}

// copied from golang libs
// func (bits floatEncoder) encode(e *encodeState, v reflect.Value, opts encOpts) {
// /snap/go/current/src/encoding/json/encode.go:574
func DetectNumberFmt(f float64, bits int) byte {
	// const bits = 64
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	return fmt
	// b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
}

func (x *Value) MarshalJSON() ([]byte, error) {
	switch tp := x.Kind.(type) {
	case *Value_NullValue:
		return []byte("null"), nil
		// return json.Marshal(tp.NullValue)
	case *Value_NumberValue:
		fmt := DetectNumberFmt(tp.NumberValue, 64)
		return []byte(strconv.FormatFloat(tp.NumberValue, fmt, -1, 64)), nil
		// return []byte(strconv.FormatFloat(tp.NumberValue, 'g', 20, 64)), nil
		// return []byte(strconv.FormatFloat(tp.NumberValue, 'g', -1, 64)), nil
		// return json.Marshal(tp.NumberValue)
	case *Value_StringValue:
		return json.Marshal(tp.StringValue)
		// fails on escaped characters. todo port this func:
		// /snap/go/current/src/encoding/json/encode.go:1030
		// func (e *encodeState) string(s string, escapeHTML bool) {
		// return []byte(`"` + tp.StringValue + `"`), nil
	case *Value_IntValue:
		return []byte(strconv.FormatInt(tp.IntValue, 10)), nil
		// return json.Marshal(tp.IntValue)
	case *Value_BoolValue:
		return []byte(strconv.FormatBool(tp.BoolValue)), nil
		// return json.Marshal(tp.BoolValue)
	case *Value_StructValue:
		return tp.StructValue.MarshalJSON()
		// return json.Marshal(tp.StructValue)
	case *Value_ListValue:
		return tp.ListValue.MarshalJSON()
		// return json.Marshal(tp.ListValue)
	case *Value_TimeValue:
		return tp.TimeValue.MarshalJSON()
		// return protojson.Marshal(tp.TimeValue)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Value.MarshalJSON(): unknown type %T", x.Kind)
	}
}

func (x *Value) UnmarshalJSON(b []byte) error {
	var src any
	if err := json.Unmarshal(b, &src); err != nil {
		return err
	}
	x.Kind = NewValueKind(src)
	return nil
	// panic("not implemented!")
	// return status.Errorf(codes.Unimplemented, "structpb.Value.UnmarshalJSON(): not implemented")
	// panic("structpb.Value.MarshalJSON(): unknown type %T", x.Kind)
	// json.
	// return protojson.Unmarshal(b, x)
}

func (x *Value) EncodeMsgpack(enc *msgpack.Encoder) error {
	switch tp := x.Kind.(type) {
	case *Value_NullValue:
		return enc.EncodeNil()
	case *Value_NumberValue:
		return enc.EncodeFloat64(tp.NumberValue)
	case *Value_StringValue:
		return enc.EncodeString(tp.StringValue)
	case *Value_IntValue:
		return enc.EncodeInt(tp.IntValue)
	case *Value_BoolValue:
		return enc.EncodeBool(tp.BoolValue)
	case *Value_StructValue:
		return tp.StructValue.EncodeMsgpack(enc)
	case *Value_ListValue:
		return tp.ListValue.EncodeMsgpack(enc)
	case *Value_TimeValue:
		return tp.TimeValue.EncodeMsgpack(enc)
	default:
		return status.Errorf(codes.InvalidArgument, "Value.EncodeMsgpack(): unknown type %T", x.Kind)
	}
}

func (x *Value) DecodeMsgpack(dec *msgpack.Decoder) error {
	var vKind any
	err := dec.Decode(&vKind)
	if err != nil {
		return err
	}

	x.Kind = NewValueKind(vKind)
	return nil
}

// NewList constructs a ListValue from a general-purpose Go slice.
// The slice elements are converted using NewValue.
func NewList(v []any) *List {
	x := make([]*Value, len(v))
	for i, v := range v {
		x[i] = NewValue(v)
		// var err error
		// x.Values[i], err = NewValue(v)
		// if err != nil {
		// 	return nil, err
		// }
	}
	return &List{List: x}
}

func (x *List) ClearItems(clearItems ...bool) {
	n := len(x.List)
	numItems := len(clearItems)
	if n > numItems {
		n = numItems
	}
	for i := 0; i < n; i++ {
		if clearItems[i] {
			x.List[i].Clear()
		}
	}
}

// AsSlice converts x to a general-purpose Go slice.
// The slice elements are converted by calling Value.AsInterface.
func (x *List) AsSlice() []any {
	vs := make([]any, len(x.GetList()))
	for i, v := range x.GetList() {
		vs[i] = v.AsInterface()
	}
	return vs
}

func (x *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.List)
	// return protojson.Marshal(x)
}

func (x *List) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &x.List)
	// return protojson.Unmarshal(b, x)
}

func (x *List) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(x.List)
}

func (x *List) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.Decode(&x.List)
}

// func (x *NullValue) MarshalJSON() ([]byte, error) {
// 	return []byte("null"), nil
// 	//json.Marshal(nil)
// 	// return protojson.Marshal(x)
// }

//	func (x *NullValue) UnmarshalJSON(b []byte) error {
//		return nil // ,json.Unmarshal(b, &x.Values)
//		// return protojson.Unmarshal(b, x)
//	}
type NonExistingValueKind struct{}

func (x *NonExistingValueKind) isValue_Kind() {}
