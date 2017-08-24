// Code generated by protoc-gen-go.
// source: ping.proto
// DO NOT EDIT!

/*
Package ping is a generated protocol buffer package.

It is generated from these files:
	ping.proto

It has these top-level messages:
	Connect
	Session
	ConnACK
	Submit
	SubmitACK
	Deliver
	DeliverACK
*/
package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConnAckCode int32

const (
	ConnAckCode_ACCEPTED           ConnAckCode = 0
	ConnAckCode_CLIENTED_REJECTED  ConnAckCode = 1
	ConnAckCode_SERVER_UNAVAILABLE ConnAckCode = 2
	ConnAckCode_TOKEN_EXPIRY       ConnAckCode = 3
	ConnAckCode_LOGINED            ConnAckCode = 4
	ConnAckCode_UPGRADE            ConnAckCode = 5
)

var ConnAckCode_name = map[int32]string{
	0: "ACCEPTED",
	1: "CLIENTED_REJECTED",
	2: "SERVER_UNAVAILABLE",
	3: "TOKEN_EXPIRY",
	4: "LOGINED",
	5: "UPGRADE",
}
var ConnAckCode_value = map[string]int32{
	"ACCEPTED":           0,
	"CLIENTED_REJECTED":  1,
	"SERVER_UNAVAILABLE": 2,
	"TOKEN_EXPIRY":       3,
	"LOGINED":            4,
	"UPGRADE":            5,
}

func (x ConnAckCode) String() string {
	return proto.EnumName(ConnAckCode_name, int32(x))
}
func (ConnAckCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Type int32

const (
	Type_Type_UNDEFINED Type = 0
	Type_SINGLE         Type = 1
	Type_GROUP          Type = 2
	Type_SYSTEM         Type = 3
	Type_ASSISTANT      Type = 4
	Type_APPLICATION    Type = 5
)

var Type_name = map[int32]string{
	0: "Type_UNDEFINED",
	1: "SINGLE",
	2: "GROUP",
	3: "SYSTEM",
	4: "ASSISTANT",
	5: "APPLICATION",
}
var Type_value = map[string]int32{
	"Type_UNDEFINED": 0,
	"SINGLE":         1,
	"GROUP":          2,
	"SYSTEM":         3,
	"ASSISTANT":      4,
	"APPLICATION":    5,
}

func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Connect struct {
	KeepAlive int32  `protobuf:"varint,1,opt,name=keepAlive" json:"keepAlive,omitempty"`
	Token     []byte `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	LoginType bool   `protobuf:"varint,3,opt,name=loginType" json:"loginType,omitempty"`
	ClientID  []byte `protobuf:"bytes,4,opt,name=clientID,proto3" json:"clientID,omitempty"`
}

func (m *Connect) Reset()                    { *m = Connect{} }
func (m *Connect) String() string            { return proto.CompactTextString(m) }
func (*Connect) ProtoMessage()               {}
func (*Connect) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Connect) GetKeepAlive() int32 {
	if m != nil {
		return m.KeepAlive
	}
	return 0
}

func (m *Connect) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *Connect) GetLoginType() bool {
	if m != nil {
		return m.LoginType
	}
	return false
}

func (m *Connect) GetClientID() []byte {
	if m != nil {
		return m.ClientID
	}
	return nil
}

type Session struct {
	Id        int32              `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Sessionid []byte             `protobuf:"bytes,2,opt,name=sessionid,proto3" json:"sessionid,omitempty"`
	MsgSize   int32              `protobuf:"varint,3,opt,name=msgSize" json:"msgSize,omitempty"`
	Delivers  map[int32]*Deliver `protobuf:"bytes,4,rep,name=delivers" json:"delivers,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Session) Reset()                    { *m = Session{} }
func (m *Session) String() string            { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()               {}
func (*Session) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Session) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Session) GetSessionid() []byte {
	if m != nil {
		return m.Sessionid
	}
	return nil
}

func (m *Session) GetMsgSize() int32 {
	if m != nil {
		return m.MsgSize
	}
	return 0
}

func (m *Session) GetDelivers() map[int32]*Deliver {
	if m != nil {
		return m.Delivers
	}
	return nil
}

type ConnACK struct {
	AckCode  ConnAckCode        `protobuf:"varint,1,opt,name=ackCode,enum=ConnAckCode" json:"ackCode,omitempty"`
	Sessions map[int32]*Session `protobuf:"bytes,2,rep,name=sessions" json:"sessions,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ConnACK) Reset()                    { *m = ConnACK{} }
func (m *ConnACK) String() string            { return proto.CompactTextString(m) }
func (*ConnACK) ProtoMessage()               {}
func (*ConnACK) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ConnACK) GetAckCode() ConnAckCode {
	if m != nil {
		return m.AckCode
	}
	return ConnAckCode_ACCEPTED
}

func (m *ConnACK) GetSessions() map[int32]*Session {
	if m != nil {
		return m.Sessions
	}
	return nil
}

// 上行
type Submit struct {
	To      []byte `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	Type    Type   `protobuf:"varint,2,opt,name=type,enum=Type" json:"type,omitempty"`
	Notify  bool   `protobuf:"varint,3,opt,name=notify" json:"notify,omitempty"`
	Offline bool   `protobuf:"varint,4,opt,name=offline" json:"offline,omitempty"`
	History bool   `protobuf:"varint,5,opt,name=history" json:"history,omitempty"`
	Payload []byte `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
	Id      uint32 `protobuf:"varint,7,opt,name=id" json:"id,omitempty"`
}

func (m *Submit) Reset()                    { *m = Submit{} }
func (m *Submit) String() string            { return proto.CompactTextString(m) }
func (*Submit) ProtoMessage()               {}
func (*Submit) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Submit) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Submit) GetType() Type {
	if m != nil {
		return m.Type
	}
	return Type_Type_UNDEFINED
}

func (m *Submit) GetNotify() bool {
	if m != nil {
		return m.Notify
	}
	return false
}

func (m *Submit) GetOffline() bool {
	if m != nil {
		return m.Offline
	}
	return false
}

func (m *Submit) GetHistory() bool {
	if m != nil {
		return m.History
	}
	return false
}

func (m *Submit) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Submit) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type SubmitACK struct {
	Id uint32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *SubmitACK) Reset()                    { *m = SubmitACK{} }
func (m *SubmitACK) String() string            { return proto.CompactTextString(m) }
func (*SubmitACK) ProtoMessage()               {}
func (*SubmitACK) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SubmitACK) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Deliver struct {
	SessionID []byte `protobuf:"bytes,1,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	From      []byte `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To        []byte `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Timestamp int64  `protobuf:"varint,4,opt,name=timestamp" json:"timestamp,omitempty"`
	Type      Type   `protobuf:"varint,5,opt,name=type,enum=Type" json:"type,omitempty"`
	Id        uint64 `protobuf:"varint,6,opt,name=id" json:"id,omitempty"`
	Payload   []byte `protobuf:"bytes,7,opt,name=payload,proto3" json:"payload,omitempty"`
	Upgrade   bool   `protobuf:"varint,8,opt,name=upgrade" json:"upgrade,omitempty"`
}

func (m *Deliver) Reset()                    { *m = Deliver{} }
func (m *Deliver) String() string            { return proto.CompactTextString(m) }
func (*Deliver) ProtoMessage()               {}
func (*Deliver) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Deliver) GetSessionID() []byte {
	if m != nil {
		return m.SessionID
	}
	return nil
}

func (m *Deliver) GetFrom() []byte {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Deliver) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Deliver) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Deliver) GetType() Type {
	if m != nil {
		return m.Type
	}
	return Type_Type_UNDEFINED
}

func (m *Deliver) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Deliver) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Deliver) GetUpgrade() bool {
	if m != nil {
		return m.Upgrade
	}
	return false
}

type DeliverACK struct {
	SessionID []byte `protobuf:"bytes,1,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	Id        uint64 `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
}

func (m *DeliverACK) Reset()                    { *m = DeliverACK{} }
func (m *DeliverACK) String() string            { return proto.CompactTextString(m) }
func (*DeliverACK) ProtoMessage()               {}
func (*DeliverACK) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *DeliverACK) GetSessionID() []byte {
	if m != nil {
		return m.SessionID
	}
	return nil
}

func (m *DeliverACK) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*Connect)(nil), "Connect")
	proto.RegisterType((*Session)(nil), "Session")
	proto.RegisterType((*ConnACK)(nil), "ConnACK")
	proto.RegisterType((*Submit)(nil), "Submit")
	proto.RegisterType((*SubmitACK)(nil), "SubmitACK")
	proto.RegisterType((*Deliver)(nil), "Deliver")
	proto.RegisterType((*DeliverACK)(nil), "DeliverACK")
	proto.RegisterEnum("ConnAckCode", ConnAckCode_name, ConnAckCode_value)
	proto.RegisterEnum("Type", Type_name, Type_value)
}

func init() { proto.RegisterFile("ping.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 660 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xc1, 0x6e, 0xda, 0x4a,
	0x14, 0x8d, 0x8d, 0x8d, 0xe1, 0x42, 0xf2, 0xfc, 0x46, 0xef, 0x45, 0x7e, 0x79, 0x51, 0x85, 0x58,
	0x54, 0x28, 0x0b, 0x16, 0x74, 0x53, 0x65, 0xe7, 0xda, 0x53, 0xe4, 0x86, 0x1a, 0x34, 0x36, 0x51,
	0xb3, 0xa8, 0x10, 0x89, 0x07, 0x3a, 0x02, 0x3c, 0x16, 0x76, 0x52, 0xd1, 0x3f, 0xaa, 0xd4, 0x0f,
	0xe9, 0xa2, 0x1f, 0x55, 0xcd, 0x78, 0x0c, 0x44, 0x95, 0xaa, 0xae, 0xe2, 0x73, 0xce, 0xcc, 0xdc,
	0x73, 0xee, 0xbd, 0x01, 0x20, 0x63, 0xe9, 0xb2, 0x9f, 0x6d, 0x79, 0xc1, 0xbb, 0x9f, 0xc1, 0xf2,
	0x78, 0x9a, 0xd2, 0x87, 0x02, 0x5d, 0x42, 0x73, 0x45, 0x69, 0xe6, 0xae, 0xd9, 0x13, 0x75, 0xb4,
	0x8e, 0xd6, 0x33, 0xc9, 0x81, 0x40, 0xff, 0x80, 0x59, 0xf0, 0x15, 0x4d, 0x1d, 0xbd, 0xa3, 0xf5,
	0xda, 0xa4, 0x04, 0xe2, 0xce, 0x9a, 0x2f, 0x59, 0x1a, 0xef, 0x32, 0xea, 0xd4, 0x3a, 0x5a, 0xaf,
	0x41, 0x0e, 0x04, 0xba, 0x80, 0xc6, 0xc3, 0x9a, 0xd1, 0xb4, 0x08, 0x7c, 0xc7, 0x90, 0xd7, 0xf6,
	0xb8, 0xfb, 0x43, 0x03, 0x2b, 0xa2, 0x79, 0xce, 0x78, 0x8a, 0xce, 0x40, 0x67, 0x89, 0x2a, 0xa9,
	0xb3, 0x44, 0xbc, 0x9a, 0x97, 0x12, 0x4b, 0x54, 0xbd, 0x03, 0x81, 0x1c, 0xb0, 0x36, 0xf9, 0x32,
	0x62, 0x5f, 0xca, 0x8a, 0x26, 0xa9, 0x20, 0x1a, 0x40, 0x23, 0xa1, 0xc2, 0xed, 0x36, 0x77, 0x8c,
	0x4e, 0xad, 0xd7, 0x1a, 0x9c, 0xf7, 0x55, 0x8d, 0xbe, 0xaf, 0x04, 0x9c, 0x16, 0xdb, 0x1d, 0xd9,
	0x9f, 0xbb, 0xc0, 0x70, 0xfa, 0x4c, 0x42, 0x36, 0xd4, 0x56, 0x74, 0xa7, 0xdc, 0x88, 0x4f, 0xf4,
	0x02, 0xcc, 0xa7, 0xf9, 0xfa, 0x91, 0x4a, 0x2b, 0xad, 0x41, 0xa3, 0x7a, 0x8b, 0x94, 0xf4, 0xb5,
	0xfe, 0x5a, 0xeb, 0x7e, 0xd3, 0xca, 0x46, 0xba, 0xde, 0x0d, 0x7a, 0x09, 0xd6, 0xfc, 0x61, 0xe5,
	0xf1, 0xa4, 0x6c, 0xe3, 0xd9, 0xa0, 0xdd, 0x97, 0x52, 0xc9, 0x91, 0x4a, 0x14, 0x76, 0x55, 0xaa,
	0xdc, 0xd1, 0x95, 0x5d, 0xf5, 0x46, 0x65, 0xbb, 0xb2, 0x5b, 0x9d, 0x13, 0x76, 0x9f, 0x49, 0x7f,
	0x62, 0x57, 0x5d, 0x38, 0xb6, 0xfb, 0x55, 0x83, 0x7a, 0xf4, 0x78, 0xbf, 0x61, 0x85, 0x68, 0x7e,
	0xc1, 0xe5, 0xfd, 0x36, 0xd1, 0x0b, 0x8e, 0xfe, 0x03, 0xa3, 0x10, 0xd3, 0xd4, 0xa5, 0x75, 0xb3,
	0x2f, 0x26, 0x49, 0x24, 0x85, 0xce, 0xa1, 0x9e, 0xf2, 0x82, 0x2d, 0x76, 0x6a, 0xd4, 0x0a, 0x89,
	0x89, 0xf0, 0xc5, 0x62, 0xcd, 0x52, 0x2a, 0xc7, 0xdc, 0x20, 0x15, 0x14, 0xca, 0x27, 0x96, 0x17,
	0x7c, 0xbb, 0x73, 0xcc, 0x52, 0x51, 0x50, 0x28, 0xd9, 0x7c, 0xb7, 0xe6, 0xf3, 0xc4, 0xa9, 0xcb,
	0xda, 0x15, 0x54, 0xdb, 0x60, 0x75, 0xb4, 0xde, 0xa9, 0xd8, 0x86, 0xee, 0xff, 0xd0, 0x2c, 0xad,
	0x8a, 0xde, 0x1e, 0x56, 0xa5, 0x14, 0xbf, 0x6b, 0x60, 0xa9, 0x71, 0x1c, 0xad, 0x4d, 0xe0, 0xab,
	0x40, 0x07, 0x02, 0x21, 0x30, 0x16, 0x5b, 0xbe, 0x51, 0xfb, 0x24, 0xbf, 0x55, 0xf6, 0xda, 0x3e,
	0xfb, 0x25, 0x34, 0x0b, 0xb6, 0xa1, 0x79, 0x31, 0xdf, 0x64, 0x32, 0x4a, 0x8d, 0x1c, 0x88, 0x7d,
	0x67, 0xcc, 0x5f, 0x3b, 0x53, 0xda, 0x12, 0x41, 0x0c, 0xb9, 0xc1, 0x47, 0xe9, 0xac, 0xe7, 0xe9,
	0x1c, 0xb0, 0x1e, 0xb3, 0xe5, 0x76, 0x9e, 0x50, 0xa7, 0x51, 0x76, 0x44, 0xc1, 0xee, 0x35, 0x80,
	0x4a, 0x22, 0x82, 0xfe, 0x3e, 0x4c, 0x59, 0x4f, 0xaf, 0xea, 0x5d, 0x3d, 0x41, 0xeb, 0x68, 0xc5,
	0x50, 0x1b, 0x1a, 0xae, 0xe7, 0xe1, 0x49, 0x8c, 0x7d, 0xfb, 0x04, 0xfd, 0x0b, 0x7f, 0x7b, 0xa3,
	0x00, 0x87, 0x31, 0xf6, 0x67, 0x04, 0xbf, 0xc3, 0x9e, 0xa0, 0x35, 0x74, 0x0e, 0x28, 0xc2, 0xe4,
	0x16, 0x93, 0xd9, 0x34, 0x74, 0x6f, 0xdd, 0x60, 0xe4, 0xbe, 0x19, 0x61, 0x5b, 0x47, 0x36, 0xb4,
	0xe3, 0xf1, 0x0d, 0x0e, 0x67, 0xf8, 0xc3, 0x24, 0x20, 0x77, 0x76, 0x0d, 0xb5, 0xc0, 0x1a, 0x8d,
	0x87, 0x41, 0x88, 0x7d, 0xdb, 0x10, 0x60, 0x3a, 0x19, 0x12, 0xd7, 0xc7, 0xb6, 0x79, 0xf5, 0x11,
	0x0c, 0xf9, 0x9f, 0x8e, 0xe0, 0x4c, 0xfc, 0x9d, 0x4d, 0x43, 0x1f, 0xbf, 0x95, 0x07, 0x4f, 0x10,
	0x40, 0x3d, 0x0a, 0xc2, 0xe1, 0x08, 0xdb, 0x1a, 0x6a, 0x82, 0x39, 0x24, 0xe3, 0xe9, 0xc4, 0xd6,
	0x25, 0x7d, 0x17, 0xc5, 0xf8, 0xbd, 0x5d, 0x43, 0xa7, 0xd0, 0x74, 0xa3, 0x28, 0x88, 0x62, 0x37,
	0x8c, 0x6d, 0x03, 0xfd, 0x05, 0x2d, 0x77, 0x32, 0x19, 0x05, 0x9e, 0x1b, 0x07, 0xe3, 0xd0, 0x36,
	0xef, 0xeb, 0xf2, 0x47, 0xea, 0xd5, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x83, 0xd8, 0x4a, 0xfd,
	0xb2, 0x04, 0x00, 0x00,
}