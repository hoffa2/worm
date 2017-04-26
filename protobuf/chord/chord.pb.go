// Code generated by protoc-gen-go.
// source: protobuf/chord/chord.proto
// DO NOT EDIT!

/*
Package chord is a generated protocol buffer package.

It is generated from these files:
	protobuf/chord/chord.proto

It has these top-level messages:
	Empty
	NewTarget
	Node
	Alive
	ToNode
	FromNode
*/
package chord

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type NewTarget struct {
	From int32 `protobuf:"varint,1,opt,name=from" json:"from,omitempty"`
	To   int32 `protobuf:"varint,2,opt,name=to" json:"to,omitempty"`
}

func (m *NewTarget) Reset()                    { *m = NewTarget{} }
func (m *NewTarget) String() string            { return proto.CompactTextString(m) }
func (*NewTarget) ProtoMessage()               {}
func (*NewTarget) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *NewTarget) GetFrom() int32 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *NewTarget) GetTo() int32 {
	if m != nil {
		return m.To
	}
	return 0
}

type Node struct {
	ID        string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	IpAddress string `protobuf:"bytes,2,opt,name=ip_address,json=ipAddress" json:"ip_address,omitempty"`
	RpcPort   string `protobuf:"bytes,3,opt,name=rpc_port,json=rpcPort" json:"rpc_port,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Node) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Node) GetIpAddress() string {
	if m != nil {
		return m.IpAddress
	}
	return ""
}

func (m *Node) GetRpcPort() string {
	if m != nil {
		return m.RpcPort
	}
	return ""
}

type Alive struct {
	IsAlive bool    `protobuf:"varint,1,opt,name=is_alive,json=isAlive" json:"is_alive,omitempty"`
	Target  int32   `protobuf:"varint,2,opt,name=target" json:"target,omitempty"`
	Nodes   []*Node `protobuf:"bytes,3,rep,name=nodes" json:"nodes,omitempty"`
}

func (m *Alive) Reset()                    { *m = Alive{} }
func (m *Alive) String() string            { return proto.CompactTextString(m) }
func (*Alive) ProtoMessage()               {}
func (*Alive) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Alive) GetIsAlive() bool {
	if m != nil {
		return m.IsAlive
	}
	return false
}

func (m *Alive) GetTarget() int32 {
	if m != nil {
		return m.Target
	}
	return 0
}

func (m *Alive) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type ToNode struct {
	// Types that are valid to be assigned to Msg:
	//	*ToNode_Shutdown
	Msg isToNode_Msg `protobuf_oneof:"msg"`
}

func (m *ToNode) Reset()                    { *m = ToNode{} }
func (m *ToNode) String() string            { return proto.CompactTextString(m) }
func (*ToNode) ProtoMessage()               {}
func (*ToNode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isToNode_Msg interface {
	isToNode_Msg()
}

type ToNode_Shutdown struct {
	Shutdown bool `protobuf:"varint,1,opt,name=shutdown,oneof"`
}

func (*ToNode_Shutdown) isToNode_Msg() {}

func (m *ToNode) GetMsg() isToNode_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *ToNode) GetShutdown() bool {
	if x, ok := m.GetMsg().(*ToNode_Shutdown); ok {
		return x.Shutdown
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ToNode) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ToNode_OneofMarshaler, _ToNode_OneofUnmarshaler, _ToNode_OneofSizer, []interface{}{
		(*ToNode_Shutdown)(nil),
	}
}

func _ToNode_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ToNode)
	// msg
	switch x := m.Msg.(type) {
	case *ToNode_Shutdown:
		t := uint64(0)
		if x.Shutdown {
			t = 1
		}
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case nil:
	default:
		return fmt.Errorf("ToNode.Msg has unexpected type %T", x)
	}
	return nil
}

func _ToNode_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ToNode)
	switch tag {
	case 1: // msg.shutdown
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Msg = &ToNode_Shutdown{x != 0}
		return true, err
	default:
		return false, nil
	}
}

func _ToNode_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ToNode)
	// msg
	switch x := m.Msg.(type) {
	case *ToNode_Shutdown:
		n += proto.SizeVarint(1<<3 | proto.WireVarint)
		n += 1
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FromNode struct {
	// Types that are valid to be assigned to Msg:
	//	*FromNode_Ok
	Msg isFromNode_Msg `protobuf_oneof:"msg"`
}

func (m *FromNode) Reset()                    { *m = FromNode{} }
func (m *FromNode) String() string            { return proto.CompactTextString(m) }
func (*FromNode) ProtoMessage()               {}
func (*FromNode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type isFromNode_Msg interface {
	isFromNode_Msg()
}

type FromNode_Ok struct {
	Ok bool `protobuf:"varint,1,opt,name=ok,oneof"`
}

func (*FromNode_Ok) isFromNode_Msg() {}

func (m *FromNode) GetMsg() isFromNode_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *FromNode) GetOk() bool {
	if x, ok := m.GetMsg().(*FromNode_Ok); ok {
		return x.Ok
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FromNode) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FromNode_OneofMarshaler, _FromNode_OneofUnmarshaler, _FromNode_OneofSizer, []interface{}{
		(*FromNode_Ok)(nil),
	}
}

func _FromNode_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FromNode)
	// msg
	switch x := m.Msg.(type) {
	case *FromNode_Ok:
		t := uint64(0)
		if x.Ok {
			t = 1
		}
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case nil:
	default:
		return fmt.Errorf("FromNode.Msg has unexpected type %T", x)
	}
	return nil
}

func _FromNode_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FromNode)
	switch tag {
	case 1: // msg.ok
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Msg = &FromNode_Ok{x != 0}
		return true, err
	default:
		return false, nil
	}
}

func _FromNode_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FromNode)
	// msg
	switch x := m.Msg.(type) {
	case *FromNode_Ok:
		n += proto.SizeVarint(1<<3 | proto.WireVarint)
		n += 1
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*NewTarget)(nil), "NewTarget")
	proto.RegisterType((*Node)(nil), "Node")
	proto.RegisterType((*Alive)(nil), "alive")
	proto.RegisterType((*ToNode)(nil), "ToNode")
	proto.RegisterType((*FromNode)(nil), "FromNode")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Chord service

type ChordClient interface {
	Alive(ctx context.Context, in *Alive, opts ...grpc.CallOption) (*Alive, error)
	Notify(ctx context.Context, in *NewTarget, opts ...grpc.CallOption) (*Alive, error)
	Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type chordClient struct {
	cc *grpc.ClientConn
}

func NewChordClient(cc *grpc.ClientConn) ChordClient {
	return &chordClient{cc}
}

func (c *chordClient) Alive(ctx context.Context, in *Alive, opts ...grpc.CallOption) (*Alive, error) {
	out := new(Alive)
	err := grpc.Invoke(ctx, "/chord/Alive", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Notify(ctx context.Context, in *NewTarget, opts ...grpc.CallOption) (*Alive, error) {
	out := new(Alive)
	err := grpc.Invoke(ctx, "/chord/Notify", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/chord/Shutdown", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Chord service

type ChordServer interface {
	Alive(context.Context, *Alive) (*Alive, error)
	Notify(context.Context, *NewTarget) (*Alive, error)
	Shutdown(context.Context, *Empty) (*Empty, error)
}

func RegisterChordServer(s *grpc.Server, srv ChordServer) {
	s.RegisterService(&_Chord_serviceDesc, srv)
}

func _Chord_Alive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Alive)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Alive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord/Alive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Alive(ctx, req.(*Alive))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Notify(ctx, req.(*NewTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chord/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Shutdown(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chord_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chord",
	HandlerType: (*ChordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Alive",
			Handler:    _Chord_Alive_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _Chord_Notify_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _Chord_Shutdown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/chord/chord.proto",
}

func init() { proto.RegisterFile("protobuf/chord/chord.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x91, 0x41, 0x6b, 0x32, 0x31,
	0x10, 0x86, 0x75, 0xd7, 0xac, 0xbb, 0xf3, 0xc1, 0x47, 0xc9, 0xa1, 0xac, 0xd6, 0x16, 0x49, 0x2f,
	0x5e, 0xba, 0x82, 0xfd, 0x05, 0x16, 0x5b, 0xea, 0x45, 0x64, 0x2b, 0xf4, 0x68, 0xd5, 0x8d, 0x1a,
	0xec, 0x3a, 0x21, 0x89, 0x15, 0xff, 0x7d, 0x71, 0x36, 0x2c, 0xbd, 0x64, 0x26, 0xef, 0xbc, 0x99,
	0x3c, 0xcc, 0x40, 0x57, 0x1b, 0x74, 0xb8, 0x3e, 0x6d, 0x87, 0x9b, 0x3d, 0x9a, 0xa2, 0x3a, 0x33,
	0x12, 0x45, 0x1b, 0xd8, 0x6b, 0xa9, 0xdd, 0x45, 0x0c, 0x21, 0x99, 0xc9, 0xf3, 0x62, 0x65, 0x76,
	0xd2, 0x71, 0x0e, 0xad, 0xad, 0xc1, 0x32, 0x6d, 0xf6, 0x9b, 0x03, 0x96, 0x53, 0xce, 0xff, 0x43,
	0xe0, 0x30, 0x0d, 0x48, 0x09, 0x1c, 0x8a, 0x39, 0xb4, 0x66, 0x58, 0xc8, 0xab, 0x3e, 0x9d, 0x90,
	0x33, 0xc9, 0x83, 0xe9, 0x84, 0xdf, 0x03, 0x28, 0xbd, 0x5c, 0x15, 0x85, 0x91, 0xd6, 0x92, 0x3f,
	0xc9, 0x13, 0xa5, 0xc7, 0x95, 0xc0, 0x3b, 0x10, 0x1b, 0xbd, 0x59, 0x6a, 0x34, 0x2e, 0x0d, 0xa9,
	0xd8, 0x36, 0x7a, 0x33, 0x47, 0xe3, 0xc4, 0x27, 0xb0, 0xd5, 0xb7, 0xfa, 0x91, 0x57, 0x8f, 0xb2,
	0x4b, 0xca, 0xa9, 0x71, 0x9c, 0xb7, 0x95, 0x1d, 0x53, 0xe9, 0x16, 0x22, 0x47, 0x8c, 0x9e, 0xc4,
	0xdf, 0xf8, 0x1d, 0xb0, 0x23, 0x16, 0xd2, 0xa6, 0x61, 0x3f, 0x1c, 0xfc, 0x1b, 0xb1, 0xec, 0xca,
	0x96, 0x57, 0x9a, 0x78, 0x82, 0x68, 0x81, 0x04, 0xdb, 0x83, 0xd8, 0xee, 0x4f, 0xae, 0xc0, 0xf3,
	0xb1, 0xea, 0xfc, 0xde, 0xc8, 0x6b, 0xe5, 0x85, 0x41, 0x58, 0xda, 0x9d, 0x78, 0x84, 0xf8, 0xcd,
	0x60, 0x49, 0x0f, 0x6e, 0x20, 0xc0, 0x43, 0x6d, 0x0d, 0xf0, 0xe0, 0x4d, 0xa3, 0x2f, 0x60, 0x34,
	0x47, 0xde, 0x01, 0x56, 0xa1, 0x45, 0x19, 0x11, 0x77, 0x7d, 0x14, 0x0d, 0xfe, 0x00, 0xd1, 0x0c,
	0x9d, 0xda, 0x5e, 0x38, 0x64, 0xf5, 0x70, 0xff, 0xd4, 0x7b, 0x10, 0x7f, 0xf8, 0xbf, 0x79, 0x94,
	0xd1, 0x1e, 0xba, 0x3e, 0x8a, 0xc6, 0x3a, 0xa2, 0x0d, 0x3d, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff,
	0x9d, 0x83, 0xcf, 0xee, 0xbf, 0x01, 0x00, 0x00,
}
