// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: proto/webss.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// WebssClient is the client API for Webss service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WebssClient interface {
	DeleteFiles(ctx context.Context, in *Files, opts ...grpc.CallOption) (*WEmpty, error)
	SetBucketLife(ctx context.Context, in *Life, opts ...grpc.CallOption) (*WEmpty, error)
	DelBucketLife(ctx context.Context, in *Lifes, opts ...grpc.CallOption) (*WEmpty, error)
	SetFileLife(ctx context.Context, in *Tag, opts ...grpc.CallOption) (*WEmpty, error)
	SignFileUrl(ctx context.Context, in *Sign, opts ...grpc.CallOption) (*SignUrl, error)
	SignFileUrls(ctx context.Context, in *Signs, opts ...grpc.CallOption) (*SignUrls, error)
	OriginalUrl(ctx context.Context, in *FName, opts ...grpc.CallOption) (*SignUrl, error)
	OriginalUrls(ctx context.Context, in *FNames, opts ...grpc.CallOption) (*SignUrls, error)
	GetFileInfo(ctx context.Context, in *File, opts ...grpc.CallOption) (*Info, error)
}

type webssClient struct {
	cc grpc.ClientConnInterface
}

func NewWebssClient(cc grpc.ClientConnInterface) WebssClient {
	return &webssClient{cc}
}

func (c *webssClient) DeleteFiles(ctx context.Context, in *Files, opts ...grpc.CallOption) (*WEmpty, error) {
	out := new(WEmpty)
	err := c.cc.Invoke(ctx, "/proto.Webss/DeleteFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) SetBucketLife(ctx context.Context, in *Life, opts ...grpc.CallOption) (*WEmpty, error) {
	out := new(WEmpty)
	err := c.cc.Invoke(ctx, "/proto.Webss/SetBucketLife", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) DelBucketLife(ctx context.Context, in *Lifes, opts ...grpc.CallOption) (*WEmpty, error) {
	out := new(WEmpty)
	err := c.cc.Invoke(ctx, "/proto.Webss/DelBucketLife", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) SetFileLife(ctx context.Context, in *Tag, opts ...grpc.CallOption) (*WEmpty, error) {
	out := new(WEmpty)
	err := c.cc.Invoke(ctx, "/proto.Webss/SetFileLife", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) SignFileUrl(ctx context.Context, in *Sign, opts ...grpc.CallOption) (*SignUrl, error) {
	out := new(SignUrl)
	err := c.cc.Invoke(ctx, "/proto.Webss/SignFileUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) SignFileUrls(ctx context.Context, in *Signs, opts ...grpc.CallOption) (*SignUrls, error) {
	out := new(SignUrls)
	err := c.cc.Invoke(ctx, "/proto.Webss/SignFileUrls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) OriginalUrl(ctx context.Context, in *FName, opts ...grpc.CallOption) (*SignUrl, error) {
	out := new(SignUrl)
	err := c.cc.Invoke(ctx, "/proto.Webss/OriginalUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) OriginalUrls(ctx context.Context, in *FNames, opts ...grpc.CallOption) (*SignUrls, error) {
	out := new(SignUrls)
	err := c.cc.Invoke(ctx, "/proto.Webss/OriginalUrls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webssClient) GetFileInfo(ctx context.Context, in *File, opts ...grpc.CallOption) (*Info, error) {
	out := new(Info)
	err := c.cc.Invoke(ctx, "/proto.Webss/GetFileInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WebssServer is the server API for Webss service.
// All implementations must embed UnimplementedWebssServer
// for forward compatibility
type WebssServer interface {
	DeleteFiles(context.Context, *Files) (*WEmpty, error)
	SetBucketLife(context.Context, *Life) (*WEmpty, error)
	DelBucketLife(context.Context, *Lifes) (*WEmpty, error)
	SetFileLife(context.Context, *Tag) (*WEmpty, error)
	SignFileUrl(context.Context, *Sign) (*SignUrl, error)
	SignFileUrls(context.Context, *Signs) (*SignUrls, error)
	OriginalUrl(context.Context, *FName) (*SignUrl, error)
	OriginalUrls(context.Context, *FNames) (*SignUrls, error)
	GetFileInfo(context.Context, *File) (*Info, error)
	mustEmbedUnimplementedWebssServer()
}

// UnimplementedWebssServer must be embedded to have forward compatible implementations.
type UnimplementedWebssServer struct {
}

func (UnimplementedWebssServer) DeleteFiles(context.Context, *Files) (*WEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFiles not implemented")
}
func (UnimplementedWebssServer) SetBucketLife(context.Context, *Life) (*WEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetBucketLife not implemented")
}
func (UnimplementedWebssServer) DelBucketLife(context.Context, *Lifes) (*WEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelBucketLife not implemented")
}
func (UnimplementedWebssServer) SetFileLife(context.Context, *Tag) (*WEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetFileLife not implemented")
}
func (UnimplementedWebssServer) SignFileUrl(context.Context, *Sign) (*SignUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignFileUrl not implemented")
}
func (UnimplementedWebssServer) SignFileUrls(context.Context, *Signs) (*SignUrls, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignFileUrls not implemented")
}
func (UnimplementedWebssServer) OriginalUrl(context.Context, *FName) (*SignUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OriginalUrl not implemented")
}
func (UnimplementedWebssServer) OriginalUrls(context.Context, *FNames) (*SignUrls, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OriginalUrls not implemented")
}
func (UnimplementedWebssServer) GetFileInfo(context.Context, *File) (*Info, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileInfo not implemented")
}
func (UnimplementedWebssServer) mustEmbedUnimplementedWebssServer() {}

// UnsafeWebssServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WebssServer will
// result in compilation errors.
type UnsafeWebssServer interface {
	mustEmbedUnimplementedWebssServer()
}

func RegisterWebssServer(s grpc.ServiceRegistrar, srv WebssServer) {
	s.RegisterService(&Webss_ServiceDesc, srv)
}

func _Webss_DeleteFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Files)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).DeleteFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/DeleteFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).DeleteFiles(ctx, req.(*Files))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_SetBucketLife_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Life)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).SetBucketLife(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/SetBucketLife",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).SetBucketLife(ctx, req.(*Life))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_DelBucketLife_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Lifes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).DelBucketLife(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/DelBucketLife",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).DelBucketLife(ctx, req.(*Lifes))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_SetFileLife_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tag)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).SetFileLife(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/SetFileLife",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).SetFileLife(ctx, req.(*Tag))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_SignFileUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Sign)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).SignFileUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/SignFileUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).SignFileUrl(ctx, req.(*Sign))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_SignFileUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Signs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).SignFileUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/SignFileUrls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).SignFileUrls(ctx, req.(*Signs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_OriginalUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).OriginalUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/OriginalUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).OriginalUrl(ctx, req.(*FName))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_OriginalUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FNames)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).OriginalUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/OriginalUrls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).OriginalUrls(ctx, req.(*FNames))
	}
	return interceptor(ctx, in, info, handler)
}

func _Webss_GetFileInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(File)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebssServer).GetFileInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Webss/GetFileInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebssServer).GetFileInfo(ctx, req.(*File))
	}
	return interceptor(ctx, in, info, handler)
}

// Webss_ServiceDesc is the grpc.ServiceDesc for Webss service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Webss_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Webss",
	HandlerType: (*WebssServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteFiles",
			Handler:    _Webss_DeleteFiles_Handler,
		},
		{
			MethodName: "SetBucketLife",
			Handler:    _Webss_SetBucketLife_Handler,
		},
		{
			MethodName: "DelBucketLife",
			Handler:    _Webss_DelBucketLife_Handler,
		},
		{
			MethodName: "SetFileLife",
			Handler:    _Webss_SetFileLife_Handler,
		},
		{
			MethodName: "SignFileUrl",
			Handler:    _Webss_SignFileUrl_Handler,
		},
		{
			MethodName: "SignFileUrls",
			Handler:    _Webss_SignFileUrls_Handler,
		},
		{
			MethodName: "OriginalUrl",
			Handler:    _Webss_OriginalUrl_Handler,
		},
		{
			MethodName: "OriginalUrls",
			Handler:    _Webss_OriginalUrls_Handler,
		},
		{
			MethodName: "GetFileInfo",
			Handler:    _Webss_GetFileInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/webss.proto",
}
