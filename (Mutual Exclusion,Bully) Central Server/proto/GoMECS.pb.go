// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: proto/GoMECS.proto

package GoMECS

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_GoMECS_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_proto_GoMECS_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_proto_GoMECS_proto_rawDescGZIP(), []int{0}
}

func (x *Query) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Query   int32 `protobuf:"varint,2,opt,name=query,proto3" json:"query,omitempty"`
	Lamport int64 `protobuf:"varint,3,opt,name=lamport,proto3" json:"lamport,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_GoMECS_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_proto_GoMECS_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_proto_GoMECS_proto_rawDescGZIP(), []int{1}
}

func (x *Task) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Task) GetQuery() int32 {
	if x != nil {
		return x.Query
	}
	return 0
}

func (x *Task) GetLamport() int64 {
	if x != nil {
		return x.Lamport
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_GoMECS_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_GoMECS_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_proto_GoMECS_proto_rawDescGZIP(), []int{2}
}

var File_proto_GoMECS_proto protoreflect.FileDescriptor

var file_proto_GoMECS_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x22, 0x17, 0x0a, 0x05,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x46, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x07, 0x0a,
	0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x84, 0x01, 0x0a, 0x04, 0x4d, 0x45, 0x43, 0x53, 0x12,
	0x29, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0d, 0x2e, 0x47, 0x6f, 0x4d,
	0x45, 0x43, 0x53, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0x0d, 0x2e, 0x47, 0x6f, 0x4d, 0x45,
	0x43, 0x53, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x07, 0x52, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x0d, 0x2e, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x2e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x1a, 0x0d, 0x2e, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x26, 0x0a, 0x05, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x12, 0x0c,
	0x2e, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x1a, 0x0d, 0x2e, 0x47,
	0x6f, 0x4d, 0x45, 0x43, 0x53, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x27, 0x5a,
	0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x6f, 0x4d, 0x45,
	0x52, 0x41, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x3b,
	0x47, 0x6f, 0x4d, 0x45, 0x43, 0x53, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_GoMECS_proto_rawDescOnce sync.Once
	file_proto_GoMECS_proto_rawDescData = file_proto_GoMECS_proto_rawDesc
)

func file_proto_GoMECS_proto_rawDescGZIP() []byte {
	file_proto_GoMECS_proto_rawDescOnce.Do(func() {
		file_proto_GoMECS_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_GoMECS_proto_rawDescData)
	})
	return file_proto_GoMECS_proto_rawDescData
}

var file_proto_GoMECS_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_GoMECS_proto_goTypes = []interface{}{
	(*Query)(nil), // 0: GoMECS.Query
	(*Task)(nil),  // 1: GoMECS.Task
	(*Empty)(nil), // 2: GoMECS.Empty
}
var file_proto_GoMECS_proto_depIdxs = []int32{
	0, // 0: GoMECS.MECS.Request:input_type -> GoMECS.Query
	0, // 1: GoMECS.MECS.Release:input_type -> GoMECS.Query
	1, // 2: GoMECS.MECS.Elect:input_type -> GoMECS.Task
	2, // 3: GoMECS.MECS.Request:output_type -> GoMECS.Empty
	2, // 4: GoMECS.MECS.Release:output_type -> GoMECS.Empty
	2, // 5: GoMECS.MECS.Elect:output_type -> GoMECS.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_GoMECS_proto_init() }
func file_proto_GoMECS_proto_init() {
	if File_proto_GoMECS_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_GoMECS_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Query); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_GoMECS_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_GoMECS_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_GoMECS_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_GoMECS_proto_goTypes,
		DependencyIndexes: file_proto_GoMECS_proto_depIdxs,
		MessageInfos:      file_proto_GoMECS_proto_msgTypes,
	}.Build()
	File_proto_GoMECS_proto = out.File
	file_proto_GoMECS_proto_rawDesc = nil
	file_proto_GoMECS_proto_goTypes = nil
	file_proto_GoMECS_proto_depIdxs = nil
}
