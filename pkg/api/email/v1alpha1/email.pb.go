// Code generated by protoc-gen-go. DO NOT EDIT.
// source: email.proto

package email

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EmailRequest struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmailRequest) Reset()         { *m = EmailRequest{} }
func (m *EmailRequest) String() string { return proto.CompactTextString(m) }
func (*EmailRequest) ProtoMessage()    {}
func (*EmailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_email_facc9f55703e9d70, []int{0}
}
func (m *EmailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmailRequest.Unmarshal(m, b)
}
func (m *EmailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmailRequest.Marshal(b, m, deterministic)
}
func (dst *EmailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmailRequest.Merge(dst, src)
}
func (m *EmailRequest) XXX_Size() int {
	return xxx_messageInfo_EmailRequest.Size(m)
}
func (m *EmailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmailRequest proto.InternalMessageInfo

func (m *EmailRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type EmailResponse struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmailResponse) Reset()         { *m = EmailResponse{} }
func (m *EmailResponse) String() string { return proto.CompactTextString(m) }
func (*EmailResponse) ProtoMessage()    {}
func (*EmailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_email_facc9f55703e9d70, []int{1}
}
func (m *EmailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmailResponse.Unmarshal(m, b)
}
func (m *EmailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmailResponse.Marshal(b, m, deterministic)
}
func (dst *EmailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmailResponse.Merge(dst, src)
}
func (m *EmailResponse) XXX_Size() int {
	return xxx_messageInfo_EmailResponse.Size(m)
}
func (m *EmailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmailResponse proto.InternalMessageInfo

func (m *EmailResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*EmailRequest)(nil), "korepta.rafal.email.v1alpha1.EmailRequest")
	proto.RegisterType((*EmailResponse)(nil), "korepta.rafal.email.v1alpha1.EmailResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EmailServiceClient is the client API for EmailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EmailServiceClient interface {
	// SendMail
	SendMail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailServiceClient struct {
	cc *grpc.ClientConn
}

func NewEmailServiceClient(cc *grpc.ClientConn) EmailServiceClient {
	return &emailServiceClient{cc}
}

func (c *emailServiceClient) SendMail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/korepta.rafal.email.v1alpha1.EmailService/SendMail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailServiceServer is the server API for EmailService service.
type EmailServiceServer interface {
	// SendMail
	SendMail(context.Context, *EmailRequest) (*EmailResponse, error)
}

func RegisterEmailServiceServer(s *grpc.Server, srv EmailServiceServer) {
	s.RegisterService(&_EmailService_serviceDesc, srv)
}

func _EmailService_SendMail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendMail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/korepta.rafal.email.v1alpha1.EmailService/SendMail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendMail(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EmailService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "korepta.rafal.email.v1alpha1.EmailService",
	HandlerType: (*EmailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMail",
			Handler:    _EmailService_SendMail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "email.proto",
}

func init() { proto.RegisterFile("email.proto", fileDescriptor_email_facc9f55703e9d70) }

var fileDescriptor_email_facc9f55703e9d70 = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0xcd, 0x4d, 0xcc,
	0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0xc9, 0xce, 0x2f, 0x4a, 0x2d, 0x28, 0x49,
	0xd4, 0x2b, 0x4a, 0x4c, 0x4b, 0xcc, 0xd1, 0x83, 0x48, 0x95, 0x19, 0x26, 0xe6, 0x14, 0x64, 0x24,
	0x1a, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x27, 0x16, 0x64, 0xea, 0x27, 0xe6,
	0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x15, 0x43, 0xf4, 0x2a, 0x69, 0x70, 0xf1, 0xb8,
	0x82, 0xd4, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x70, 0xb1, 0xe7, 0xa6, 0x16,
	0x17, 0x27, 0xa6, 0xa7, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a, 0xaa, 0x5c,
	0xbc, 0x50, 0x95, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x22, 0x5c, 0xac, 0xa9, 0x45, 0x45,
	0xf9, 0x45, 0x50, 0x85, 0x10, 0x8e, 0x51, 0x3f, 0x23, 0xd4, 0xc4, 0xe0, 0xd4, 0xa2, 0xb2, 0xcc,
	0xe4, 0x54, 0xa1, 0x7a, 0x2e, 0x8e, 0xe0, 0xd4, 0xbc, 0x14, 0xdf, 0xc4, 0xcc, 0x1c, 0x21, 0x2d,
	0x3d, 0x7c, 0x4e, 0xd5, 0x43, 0x76, 0x89, 0x94, 0x36, 0x51, 0x6a, 0x21, 0x6e, 0x51, 0x92, 0x6a,
	0xba, 0xfc, 0x64, 0x32, 0x93, 0x88, 0x12, 0xbf, 0x3e, 0x4c, 0x81, 0x3e, 0x58, 0xbd, 0x15, 0xa3,
	0x96, 0x93, 0x29, 0x97, 0x42, 0x72, 0x7e, 0x2e, 0x5e, 0xd3, 0x9c, 0x38, 0xc0, 0xc6, 0x39, 0x06,
	0x78, 0x06, 0x30, 0x46, 0xb1, 0x82, 0xe5, 0x92, 0xd8, 0xc0, 0x01, 0x64, 0x0c, 0x08, 0x00, 0x00,
	0xff, 0xff, 0x25, 0xa7, 0xd6, 0x2d, 0x6b, 0x01, 0x00, 0x00,
}