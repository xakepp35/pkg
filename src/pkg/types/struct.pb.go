// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
// https://developers.google.com/protocol-buffers/
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: pkg/types/struct.proto

// package google.protobuf;

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// `NullValue` is a singleton enumeration to represent the null value for the
// `Value` type union.
//
//	The JSON representation for `NullValue` is JSON `null`.
type Null int32

const (
	// Null value.
	Null_NULL_VALUE Null = 0
)

// Enum value maps for Null.
var (
	Null_name = map[int32]string{
		0: "NULL_VALUE",
	}
	Null_value = map[string]int32{
		"NULL_VALUE": 0,
	}
)

func (x Null) Enum() *Null {
	p := new(Null)
	*p = x
	return p
}

func (x Null) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Null) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_types_struct_proto_enumTypes[0].Descriptor()
}

func (Null) Type() protoreflect.EnumType {
	return &file_pkg_types_struct_proto_enumTypes[0]
}

func (x Null) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Null.Descriptor instead.
func (Null) EnumDescriptor() ([]byte, []int) {
	return file_pkg_types_struct_proto_rawDescGZIP(), []int{0}
}

// `Struct` represents a structured data value, consisting of fields
// which map to dynamically typed values. In some languages, `Struct`
// might be supported by a native representation. For example, in
// scripting languages like JS a struct is represented as an
// object. The details of that representation are described together
// with the proto support for the language.
//
// The JSON representation for `Struct` is JSON object.
type Struct struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unordered map of dynamically typed values.
	Fields        map[string]*Value `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Struct) Reset() {
	*x = Struct{}
	mi := &file_pkg_types_struct_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Struct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Struct) ProtoMessage() {}

func (x *Struct) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_types_struct_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Struct.ProtoReflect.Descriptor instead.
func (*Struct) Descriptor() ([]byte, []int) {
	return file_pkg_types_struct_proto_rawDescGZIP(), []int{0}
}

func (x *Struct) GetFields() map[string]*Value {
	if x != nil {
		return x.Fields
	}
	return nil
}

// `Value` represents a dynamically typed value which can be either
// null, a number, a string, a boolean, a recursive struct value, or a
// list of values. A producer of value is expected to set one of that
// variants, absence of any variant indicates an error.
//
// The JSON representation for `Value` is JSON value.
type Value struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The kind of value.
	//
	// Types that are valid to be assigned to Kind:
	//
	//	*Value_NullValue
	//	*Value_NumberValue
	//	*Value_StringValue
	//	*Value_BoolValue
	//	*Value_StructValue
	//	*Value_ListValue
	//	*Value_TimeValue
	//	*Value_IntValue
	Kind          isValue_Kind `protobuf_oneof:"kind"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Value) Reset() {
	*x = Value{}
	mi := &file_pkg_types_struct_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_types_struct_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_pkg_types_struct_proto_rawDescGZIP(), []int{1}
}

func (x *Value) GetKind() isValue_Kind {
	if x != nil {
		return x.Kind
	}
	return nil
}

func (x *Value) GetNullValue() Null {
	if x != nil {
		if x, ok := x.Kind.(*Value_NullValue); ok {
			return x.NullValue
		}
	}
	return Null_NULL_VALUE
}

func (x *Value) GetNumberValue() float64 {
	if x != nil {
		if x, ok := x.Kind.(*Value_NumberValue); ok {
			return x.NumberValue
		}
	}
	return 0
}

func (x *Value) GetStringValue() string {
	if x != nil {
		if x, ok := x.Kind.(*Value_StringValue); ok {
			return x.StringValue
		}
	}
	return ""
}

func (x *Value) GetBoolValue() bool {
	if x != nil {
		if x, ok := x.Kind.(*Value_BoolValue); ok {
			return x.BoolValue
		}
	}
	return false
}

func (x *Value) GetStructValue() *Struct {
	if x != nil {
		if x, ok := x.Kind.(*Value_StructValue); ok {
			return x.StructValue
		}
	}
	return nil
}

func (x *Value) GetListValue() *List {
	if x != nil {
		if x, ok := x.Kind.(*Value_ListValue); ok {
			return x.ListValue
		}
	}
	return nil
}

func (x *Value) GetTimeValue() *Time {
	if x != nil {
		if x, ok := x.Kind.(*Value_TimeValue); ok {
			return x.TimeValue
		}
	}
	return nil
}

func (x *Value) GetIntValue() int64 {
	if x != nil {
		if x, ok := x.Kind.(*Value_IntValue); ok {
			return x.IntValue
		}
	}
	return 0
}

type isValue_Kind interface {
	isValue_Kind()
}

type Value_NullValue struct {
	// Represents a null value.
	NullValue Null `protobuf:"varint,1,opt,name=null_value,json=nullValue,proto3,enum=pkg.types.Null,oneof"`
}

type Value_NumberValue struct {
	// Represents a double value.
	NumberValue float64 `protobuf:"fixed64,2,opt,name=number_value,json=numberValue,proto3,oneof"`
}

type Value_StringValue struct {
	// Represents a string value.
	StringValue string `protobuf:"bytes,3,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type Value_BoolValue struct {
	// Represents a boolean value.
	BoolValue bool `protobuf:"varint,4,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Value_StructValue struct {
	// Represents a structured value.
	StructValue *Struct `protobuf:"bytes,5,opt,name=struct_value,json=structValue,proto3,oneof"`
}

type Value_ListValue struct {
	// Represents a repeated `Value`.
	ListValue *List `protobuf:"bytes,6,opt,name=list_value,json=listValue,proto3,oneof"`
}

type Value_TimeValue struct {
	// Represents datetime stuff
	TimeValue *Time `protobuf:"bytes,7,opt,name=time_value,json=timeValue,proto3,oneof"`
}

type Value_IntValue struct {
	// Represents a int value.
	IntValue int64 `protobuf:"varint,8,opt,name=int_value,json=intValue,proto3,oneof"`
}

func (*Value_NullValue) isValue_Kind() {}

func (*Value_NumberValue) isValue_Kind() {}

func (*Value_StringValue) isValue_Kind() {}

func (*Value_BoolValue) isValue_Kind() {}

func (*Value_StructValue) isValue_Kind() {}

func (*Value_ListValue) isValue_Kind() {}

func (*Value_TimeValue) isValue_Kind() {}

func (*Value_IntValue) isValue_Kind() {}

// `ListValue` is a wrapper around a repeated field of values.
//
// The JSON representation for `ListValue` is JSON array.
type List struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Repeated field of dynamically typed values.
	List          []*Value `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *List) Reset() {
	*x = List{}
	mi := &file_pkg_types_struct_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *List) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*List) ProtoMessage() {}

func (x *List) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_types_struct_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use List.ProtoReflect.Descriptor instead.
func (*List) Descriptor() ([]byte, []int) {
	return file_pkg_types_struct_proto_rawDescGZIP(), []int{2}
}

func (x *List) GetList() []*Value {
	if x != nil {
		return x.List
	}
	return nil
}

var File_pkg_types_struct_proto protoreflect.FileDescriptor

const file_pkg_types_struct_proto_rawDesc = "" +
	"\n" +
	"\x16pkg/types/struct.proto\x12\tpkg.types\x1a\x14pkg/types/time.proto\"\x8c\x01\n" +
	"\x06Struct\x125\n" +
	"\x06fields\x18\x01 \x03(\v2\x1d.pkg.types.Struct.FieldsEntryR\x06fields\x1aK\n" +
	"\vFieldsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12&\n" +
	"\x05value\x18\x02 \x01(\v2\x10.pkg.types.ValueR\x05value:\x028\x01\"\xe7\x02\n" +
	"\x05Value\x120\n" +
	"\n" +
	"null_value\x18\x01 \x01(\x0e2\x0f.pkg.types.NullH\x00R\tnullValue\x12#\n" +
	"\fnumber_value\x18\x02 \x01(\x01H\x00R\vnumberValue\x12#\n" +
	"\fstring_value\x18\x03 \x01(\tH\x00R\vstringValue\x12\x1f\n" +
	"\n" +
	"bool_value\x18\x04 \x01(\bH\x00R\tboolValue\x126\n" +
	"\fstruct_value\x18\x05 \x01(\v2\x11.pkg.types.StructH\x00R\vstructValue\x120\n" +
	"\n" +
	"list_value\x18\x06 \x01(\v2\x0f.pkg.types.ListH\x00R\tlistValue\x120\n" +
	"\n" +
	"time_value\x18\a \x01(\v2\x0f.pkg.types.TimeH\x00R\ttimeValue\x12\x1d\n" +
	"\tint_value\x18\b \x01(\x03H\x00R\bintValueB\x06\n" +
	"\x04kind\",\n" +
	"\x04List\x12$\n" +
	"\x04list\x18\x01 \x03(\v2\x10.pkg.types.ValueR\x04list*\x16\n" +
	"\x04Null\x12\x0e\n" +
	"\n" +
	"NULL_VALUE\x10\x00B-Z+github.com/xakepp35/pkg/src/pkg/types;typesb\x06proto3"

var (
	file_pkg_types_struct_proto_rawDescOnce sync.Once
	file_pkg_types_struct_proto_rawDescData []byte
)

func file_pkg_types_struct_proto_rawDescGZIP() []byte {
	file_pkg_types_struct_proto_rawDescOnce.Do(func() {
		file_pkg_types_struct_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_types_struct_proto_rawDesc), len(file_pkg_types_struct_proto_rawDesc)))
	})
	return file_pkg_types_struct_proto_rawDescData
}

var file_pkg_types_struct_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_types_struct_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_types_struct_proto_goTypes = []any{
	(Null)(0),      // 0: pkg.types.Null
	(*Struct)(nil), // 1: pkg.types.Struct
	(*Value)(nil),  // 2: pkg.types.Value
	(*List)(nil),   // 3: pkg.types.List
	nil,            // 4: pkg.types.Struct.FieldsEntry
	(*Time)(nil),   // 5: pkg.types.Time
}
var file_pkg_types_struct_proto_depIdxs = []int32{
	4, // 0: pkg.types.Struct.fields:type_name -> pkg.types.Struct.FieldsEntry
	0, // 1: pkg.types.Value.null_value:type_name -> pkg.types.Null
	1, // 2: pkg.types.Value.struct_value:type_name -> pkg.types.Struct
	3, // 3: pkg.types.Value.list_value:type_name -> pkg.types.List
	5, // 4: pkg.types.Value.time_value:type_name -> pkg.types.Time
	2, // 5: pkg.types.List.list:type_name -> pkg.types.Value
	2, // 6: pkg.types.Struct.FieldsEntry.value:type_name -> pkg.types.Value
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_pkg_types_struct_proto_init() }
func file_pkg_types_struct_proto_init() {
	if File_pkg_types_struct_proto != nil {
		return
	}
	file_pkg_types_time_proto_init()
	file_pkg_types_struct_proto_msgTypes[1].OneofWrappers = []any{
		(*Value_NullValue)(nil),
		(*Value_NumberValue)(nil),
		(*Value_StringValue)(nil),
		(*Value_BoolValue)(nil),
		(*Value_StructValue)(nil),
		(*Value_ListValue)(nil),
		(*Value_TimeValue)(nil),
		(*Value_IntValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_types_struct_proto_rawDesc), len(file_pkg_types_struct_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_types_struct_proto_goTypes,
		DependencyIndexes: file_pkg_types_struct_proto_depIdxs,
		EnumInfos:         file_pkg_types_struct_proto_enumTypes,
		MessageInfos:      file_pkg_types_struct_proto_msgTypes,
	}.Build()
	File_pkg_types_struct_proto = out.File
	file_pkg_types_struct_proto_goTypes = nil
	file_pkg_types_struct_proto_depIdxs = nil
}
