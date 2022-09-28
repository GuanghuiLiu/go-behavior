// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.9.1
// source: s2s.proto

package __

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

//
type Base struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From    string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To      string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	EventID uint64 `protobuf:"varint,3,opt,name=eventID,proto3" json:"eventID,omitempty"`
}

func (x *Base) Reset() {
	*x = Base{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s2s_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Base) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Base) ProtoMessage() {}

func (x *Base) ProtoReflect() protoreflect.Message {
	mi := &file_s2s_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Base.ProtoReflect.Descriptor instead.
func (*Base) Descriptor() ([]byte, []int) {
	return file_s2s_proto_rawDescGZIP(), []int{0}
}

func (x *Base) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Base) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Base) GetEventID() uint64 {
	if x != nil {
		return x.EventID
	}
	return 0
}

type CommonS2S struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base    *Base  `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Key     uint64 `protobuf:"varint,2,opt,name=key,proto3" json:"key,omitempty"`
	Message []byte `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CommonS2S) Reset() {
	*x = CommonS2S{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s2s_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonS2S) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonS2S) ProtoMessage() {}

func (x *CommonS2S) ProtoReflect() protoreflect.Message {
	mi := &file_s2s_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonS2S.ProtoReflect.Descriptor instead.
func (*CommonS2S) Descriptor() ([]byte, []int) {
	return file_s2s_proto_rawDescGZIP(), []int{1}
}

func (x *CommonS2S) GetBase() *Base {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *CommonS2S) GetKey() uint64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *CommonS2S) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

type SyncResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base    *Base  `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Code    uint32 `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Result  []byte `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SyncResult) Reset() {
	*x = SyncResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s2s_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncResult) ProtoMessage() {}

func (x *SyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_s2s_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncResult.ProtoReflect.Descriptor instead.
func (*SyncResult) Descriptor() ([]byte, []int) {
	return file_s2s_proto_rawDescGZIP(), []int{2}
}

func (x *SyncResult) GetBase() *Base {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *SyncResult) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *SyncResult) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *SyncResult) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

type StarModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base *Base  `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Uid  string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *StarModel) Reset() {
	*x = StarModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s2s_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StarModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StarModel) ProtoMessage() {}

func (x *StarModel) ProtoReflect() protoreflect.Message {
	mi := &file_s2s_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StarModel.ProtoReflect.Descriptor instead.
func (*StarModel) Descriptor() ([]byte, []int) {
	return file_s2s_proto_rawDescGZIP(), []int{3}
}

func (x *StarModel) GetBase() *Base {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *StarModel) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

type StopModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base *Base `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
}

func (x *StopModel) Reset() {
	*x = StopModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s2s_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopModel) ProtoMessage() {}

func (x *StopModel) ProtoReflect() protoreflect.Message {
	mi := &file_s2s_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopModel.ProtoReflect.Descriptor instead.
func (*StopModel) Descriptor() ([]byte, []int) {
	return file_s2s_proto_rawDescGZIP(), []int{4}
}

func (x *StopModel) GetBase() *Base {
	if x != nil {
		return x.Base
	}
	return nil
}

var File_s2s_proto protoreflect.FileDescriptor

var file_s2s_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x32, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22,
	0x44, 0x0a, 0x04, 0x42, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x49, 0x44, 0x22, 0x55, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53,
	0x32, 0x53, 0x12, 0x1c, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x70, 0x0a, 0x0a,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1c, 0x0a, 0x04, 0x62, 0x61,
	0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x42, 0x61,
	0x73, 0x65, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x3b,
	0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x04, 0x62,
	0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x42,
	0x61, 0x73, 0x65, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x29, 0x0a, 0x09, 0x53,
	0x74, 0x6f, 0x70, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x65,
	0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2f, 0x2e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_s2s_proto_rawDescOnce sync.Once
	file_s2s_proto_rawDescData = file_s2s_proto_rawDesc
)

func file_s2s_proto_rawDescGZIP() []byte {
	file_s2s_proto_rawDescOnce.Do(func() {
		file_s2s_proto_rawDescData = protoimpl.X.CompressGZIP(file_s2s_proto_rawDescData)
	})
	return file_s2s_proto_rawDescData
}

var file_s2s_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_s2s_proto_goTypes = []interface{}{
	(*Base)(nil),       // 0: pb.Base
	(*CommonS2S)(nil),  // 1: pb.CommonS2S
	(*SyncResult)(nil), // 2: pb.SyncResult
	(*StarModel)(nil),  // 3: pb.StarModel
	(*StopModel)(nil),  // 4: pb.StopModel
}
var file_s2s_proto_depIdxs = []int32{
	0, // 0: pb.CommonS2S.base:type_name -> pb.Base
	0, // 1: pb.SyncResult.base:type_name -> pb.Base
	0, // 2: pb.StarModel.base:type_name -> pb.Base
	0, // 3: pb.StopModel.base:type_name -> pb.Base
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_s2s_proto_init() }
func file_s2s_proto_init() {
	if File_s2s_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_s2s_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Base); i {
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
		file_s2s_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonS2S); i {
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
		file_s2s_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncResult); i {
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
		file_s2s_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StarModel); i {
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
		file_s2s_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopModel); i {
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
			RawDescriptor: file_s2s_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_s2s_proto_goTypes,
		DependencyIndexes: file_s2s_proto_depIdxs,
		MessageInfos:      file_s2s_proto_msgTypes,
	}.Build()
	File_s2s_proto = out.File
	file_s2s_proto_rawDesc = nil
	file_s2s_proto_goTypes = nil
	file_s2s_proto_depIdxs = nil
}
