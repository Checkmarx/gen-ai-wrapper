//protoc redirect_prompt.proto -Iapi/client -Ithird_party --go_out=plugins=grpc:./internal/

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.5
// source: redirect_prompt.proto

package redirect_prompt

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RedirectPromptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant    string `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	RequestId string `protobuf:"bytes,2,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Origin    string `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	Feature   string `protobuf:"bytes,4,opt,name=feature,proto3" json:"feature,omitempty"`
	Content   []byte `protobuf:"bytes,10,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *RedirectPromptRequest) Reset() {
	*x = RedirectPromptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redirect_prompt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedirectPromptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectPromptRequest) ProtoMessage() {}

func (x *RedirectPromptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_redirect_prompt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectPromptRequest.ProtoReflect.Descriptor instead.
func (*RedirectPromptRequest) Descriptor() ([]byte, []int) {
	return file_redirect_prompt_proto_rawDescGZIP(), []int{0}
}

func (x *RedirectPromptRequest) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

func (x *RedirectPromptRequest) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *RedirectPromptRequest) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *RedirectPromptRequest) GetFeature() string {
	if x != nil {
		return x.Feature
	}
	return ""
}

func (x *RedirectPromptRequest) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type RedirectPromptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant         string `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	RequestId      string `protobuf:"bytes,2,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Origin         string `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	GenAiErrorCode int32  `protobuf:"varint,8,opt,name=gen_ai_error_code,json=genAiErrorCode,proto3" json:"gen_ai_error_code,omitempty"`
	Content        []byte `protobuf:"bytes,10,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *RedirectPromptResponse) Reset() {
	*x = RedirectPromptResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redirect_prompt_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedirectPromptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectPromptResponse) ProtoMessage() {}

func (x *RedirectPromptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_redirect_prompt_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectPromptResponse.ProtoReflect.Descriptor instead.
func (*RedirectPromptResponse) Descriptor() ([]byte, []int) {
	return file_redirect_prompt_proto_rawDescGZIP(), []int{1}
}

func (x *RedirectPromptResponse) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

func (x *RedirectPromptResponse) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *RedirectPromptResponse) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *RedirectPromptResponse) GetGenAiErrorCode() int32 {
	if x != nil {
		return x.GenAiErrorCode
	}
	return 0
}

func (x *RedirectPromptResponse) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_redirect_prompt_proto protoreflect.FileDescriptor

var file_redirect_prompt_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x15, 0x52, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12,
	0x18, 0x0a, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x22, 0xac, 0x01, 0x0a, 0x16, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x29, 0x0a,
	0x11, 0x67, 0x65, 0x6e, 0x5f, 0x61, 0x69, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x67, 0x65, 0x6e, 0x41, 0x69, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x32, 0x75, 0x0a, 0x0e, 0x41, 0x69, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x0e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x26, 0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27,
	0x2e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x61, 0x70, 0x69,
	0x2f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_redirect_prompt_proto_rawDescOnce sync.Once
	file_redirect_prompt_proto_rawDescData = file_redirect_prompt_proto_rawDesc
)

func file_redirect_prompt_proto_rawDescGZIP() []byte {
	file_redirect_prompt_proto_rawDescOnce.Do(func() {
		file_redirect_prompt_proto_rawDescData = protoimpl.X.CompressGZIP(file_redirect_prompt_proto_rawDescData)
	})
	return file_redirect_prompt_proto_rawDescData
}

var file_redirect_prompt_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_redirect_prompt_proto_goTypes = []interface{}{
	(*RedirectPromptRequest)(nil),  // 0: redirect_prompt.RedirectPromptRequest
	(*RedirectPromptResponse)(nil), // 1: redirect_prompt.RedirectPromptResponse
}
var file_redirect_prompt_proto_depIdxs = []int32{
	0, // 0: redirect_prompt.AiProxyService.RedirectPrompt:input_type -> redirect_prompt.RedirectPromptRequest
	1, // 1: redirect_prompt.AiProxyService.RedirectPrompt:output_type -> redirect_prompt.RedirectPromptResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_redirect_prompt_proto_init() }
func file_redirect_prompt_proto_init() {
	if File_redirect_prompt_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_redirect_prompt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedirectPromptRequest); i {
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
		file_redirect_prompt_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedirectPromptResponse); i {
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
			RawDescriptor: file_redirect_prompt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_redirect_prompt_proto_goTypes,
		DependencyIndexes: file_redirect_prompt_proto_depIdxs,
		MessageInfos:      file_redirect_prompt_proto_msgTypes,
	}.Build()
	File_redirect_prompt_proto = out.File
	file_redirect_prompt_proto_rawDesc = nil
	file_redirect_prompt_proto_goTypes = nil
	file_redirect_prompt_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AiProxyServiceClient is the client API for AiProxyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AiProxyServiceClient interface {
	RedirectPrompt(ctx context.Context, in *RedirectPromptRequest, opts ...grpc.CallOption) (*RedirectPromptResponse, error)
}

type aiProxyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAiProxyServiceClient(cc grpc.ClientConnInterface) AiProxyServiceClient {
	return &aiProxyServiceClient{cc}
}

func (c *aiProxyServiceClient) RedirectPrompt(ctx context.Context, in *RedirectPromptRequest, opts ...grpc.CallOption) (*RedirectPromptResponse, error) {
	out := new(RedirectPromptResponse)
	err := c.cc.Invoke(ctx, "/redirect_prompt.AiProxyService/RedirectPrompt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiProxyServiceServer is the server API for AiProxyService service.
type AiProxyServiceServer interface {
	RedirectPrompt(context.Context, *RedirectPromptRequest) (*RedirectPromptResponse, error)
}

// UnimplementedAiProxyServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAiProxyServiceServer struct {
}

func (*UnimplementedAiProxyServiceServer) RedirectPrompt(context.Context, *RedirectPromptRequest) (*RedirectPromptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RedirectPrompt not implemented")
}

func RegisterAiProxyServiceServer(s *grpc.Server, srv AiProxyServiceServer) {
	s.RegisterService(&_AiProxyService_serviceDesc, srv)
}

func _AiProxyService_RedirectPrompt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RedirectPromptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiProxyServiceServer).RedirectPrompt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redirect_prompt.AiProxyService/RedirectPrompt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiProxyServiceServer).RedirectPrompt(ctx, req.(*RedirectPromptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AiProxyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "redirect_prompt.AiProxyService",
	HandlerType: (*AiProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RedirectPrompt",
			Handler:    _AiProxyService_RedirectPrompt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "redirect_prompt.proto",
}
