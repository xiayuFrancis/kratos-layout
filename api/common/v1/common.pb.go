// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.7.1
// source: api/common/v1/common.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

// 通用响应格式
type Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 状态码，0 表示成功，非 0 表示错误
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// 响应消息
	Msg string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	// 扩展数据，使用 JSON 对象
	Ext           *structpb.Struct `protobuf:"bytes,3,opt,name=ext,proto3" json:"ext,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_api_common_v1_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_common_v1_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_api_common_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Response) GetExt() *structpb.Struct {
	if x != nil {
		return x.Ext
	}
	return nil
}

// 空请求
type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_api_common_v1_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_api_common_v1_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_common_v1_common_proto_rawDescGZIP(), []int{1}
}

var File_api_common_v1_common_proto protoreflect.FileDescriptor

const file_api_common_v1_common_proto_rawDesc = "" +
	"\n" +
	"\x1aapi/common/v1/common.proto\x12\rapi.common.v1\x1a\x1cgoogle/protobuf/struct.proto\"[\n" +
	"\bResponse\x12\x12\n" +
	"\x04code\x18\x01 \x01(\x05R\x04code\x12\x10\n" +
	"\x03msg\x18\x02 \x01(\tR\x03msg\x12)\n" +
	"\x03ext\x18\x03 \x01(\v2\x17.google.protobuf.StructR\x03ext\"\a\n" +
	"\x05EmptyB\x1dZ\x1bkratosdemo/api/common/v1;v1b\x06proto3"

var (
	file_api_common_v1_common_proto_rawDescOnce sync.Once
	file_api_common_v1_common_proto_rawDescData []byte
)

func file_api_common_v1_common_proto_rawDescGZIP() []byte {
	file_api_common_v1_common_proto_rawDescOnce.Do(func() {
		file_api_common_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_common_v1_common_proto_rawDesc), len(file_api_common_v1_common_proto_rawDesc)))
	})
	return file_api_common_v1_common_proto_rawDescData
}

var file_api_common_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_common_v1_common_proto_goTypes = []any{
	(*Response)(nil),        // 0: api.common.v1.Response
	(*Empty)(nil),           // 1: api.common.v1.Empty
	(*structpb.Struct)(nil), // 2: google.protobuf.Struct
}
var file_api_common_v1_common_proto_depIdxs = []int32{
	2, // 0: api.common.v1.Response.ext:type_name -> google.protobuf.Struct
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_common_v1_common_proto_init() }
func file_api_common_v1_common_proto_init() {
	if File_api_common_v1_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_common_v1_common_proto_rawDesc), len(file_api_common_v1_common_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_common_v1_common_proto_goTypes,
		DependencyIndexes: file_api_common_v1_common_proto_depIdxs,
		MessageInfos:      file_api_common_v1_common_proto_msgTypes,
	}.Build()
	File_api_common_v1_common_proto = out.File
	file_api_common_v1_common_proto_goTypes = nil
	file_api_common_v1_common_proto_depIdxs = nil
}
