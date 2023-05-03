// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: orbis/transport/v1alpha1/transport.proto

package transportv1alpha1

import (
	pb "github.com/libp2p/go-libp2p/core/crypto/pb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetHostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetHostRequest) Reset() {
	*x = GetHostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetHostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHostRequest) ProtoMessage() {}

func (x *GetHostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHostRequest.ProtoReflect.Descriptor instead.
func (*GetHostRequest) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{0}
}

type GetHostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *GetHostResponse) Reset() {
	*x = GetHostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetHostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHostResponse) ProtoMessage() {}

func (x *GetHostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHostResponse.ProtoReflect.Descriptor instead.
func (*GetHostResponse) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{1}
}

func (x *GetHostResponse) GetNode() *Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node    *Node    `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Message *Message `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SendRequest) Reset() {
	*x = SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRequest) ProtoMessage() {}

func (x *SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRequest.ProtoReflect.Descriptor instead.
func (*SendRequest) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{2}
}

func (x *SendRequest) GetNode() *Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *SendRequest) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type SendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendResponse) Reset() {
	*x = SendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendResponse) ProtoMessage() {}

func (x *SendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendResponse.ProtoReflect.Descriptor instead.
func (*SendResponse) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{3}
}

type ConnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"` // multiaddress
}

func (x *ConnectRequest) Reset() {
	*x = ConnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectRequest) ProtoMessage() {}

func (x *ConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectRequest.ProtoReflect.Descriptor instead.
func (*ConnectRequest) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{4}
}

func (x *ConnectRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ConnectRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type ConnectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectResponse) Reset() {
	*x = ConnectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectResponse) ProtoMessage() {}

func (x *ConnectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectResponse.ProtoReflect.Descriptor instead.
func (*ConnectResponse) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{5}
}

type GossipRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Message *Message `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *GossipRequest) Reset() {
	*x = GossipRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipRequest) ProtoMessage() {}

func (x *GossipRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GossipRequest.ProtoReflect.Descriptor instead.
func (*GossipRequest) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{6}
}

func (x *GossipRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *GossipRequest) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type GossipResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GossipResponse) Reset() {
	*x = GossipResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipResponse) ProtoMessage() {}

func (x *GossipResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GossipResponse.ProtoReflect.Descriptor instead.
func (*GossipResponse) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{7}
}

type NewMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Gossip  bool   `protobuf:"varint,2,opt,name=gossip,proto3" json:"gossip,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Type    string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	RingId  string `protobuf:"bytes,5,opt,name=ring_id,json=ringId,proto3" json:"ring_id,omitempty"`
}

func (x *NewMessageRequest) Reset() {
	*x = NewMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewMessageRequest) ProtoMessage() {}

func (x *NewMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewMessageRequest.ProtoReflect.Descriptor instead.
func (*NewMessageRequest) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{8}
}

func (x *NewMessageRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *NewMessageRequest) GetGossip() bool {
	if x != nil {
		return x.Gossip
	}
	return false
}

func (x *NewMessageRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *NewMessageRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *NewMessageRequest) GetRingId() string {
	if x != nil {
		return x.RingId
	}
	return ""
}

type NewMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *NewMessageResponse) Reset() {
	*x = NewMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewMessageResponse) ProtoMessage() {}

func (x *NewMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewMessageResponse.ProtoReflect.Descriptor instead.
func (*NewMessageResponse) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{9}
}

func (x *NewMessageResponse) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address   string        `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"` // multiaddress
	PublicKey *pb.PublicKey `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{10}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Node) GetPublicKey() *pb.PublicKey {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp  int64  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                      // unix time
	Id         string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`                                     // message id
	Type       string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`                                 // message type
	Payload    []byte `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`                           // generic payload
	Gossip     bool   `protobuf:"varint,5,opt,name=gossip,proto3" json:"gossip,omitempty"`                            // gossip message over pubsub
	NodeId     string `protobuf:"bytes,6,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`               // author node id (peer.ID)
	NodePubKey []byte `protobuf:"bytes,7,opt,name=node_pub_key,json=nodePubKey,proto3" json:"node_pub_key,omitempty"` // authoring node pubkey (32bytes)
	RingId     string `protobuf:"bytes,8,opt,name=ring_id,json=ringId,proto3" json:"ring_id,omitempty"`               //
	Signature  []byte `protobuf:"bytes,9,opt,name=signature,proto3" json:"signature,omitempty"`                       // signature of the message (including payload)
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_orbis_transport_v1alpha1_transport_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP(), []int{11}
}

func (x *Message) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Message) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Message) GetGossip() bool {
	if x != nil {
		return x.Gossip
	}
	return false
}

func (x *Message) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *Message) GetNodePubKey() []byte {
	if x != nil {
		return x.NodePubKey
	}
	return nil
}

func (x *Message) GetRingId() string {
	if x != nil {
		return x.RingId
	}
	return ""
}

func (x *Message) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_orbis_transport_v1alpha1_transport_proto protoreflect.FileDescriptor

var file_orbis_transport_v1alpha1_transport_proto_rawDesc = []byte{
	0x0a, 0x28, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x6f, 0x72, 0x62, 0x69,
	0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1d, 0x6c, 0x69, 0x62, 0x70, 0x32, 0x70, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x45, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x7e, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x3b, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6f, 0x72,
	0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3a, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x0a, 0x0d, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x3b, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x82, 0x01, 0x0a,
	0x11, 0x4e, 0x65, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x69, 0x6e, 0x67,
	0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x69, 0x6e, 0x67, 0x49,
	0x64, 0x22, 0x51, 0x0a, 0x12, 0x4e, 0x65, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73,
	0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x6c, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3a, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x62,
	0x70, 0x32, 0x70, 0x2e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x22, 0xef, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x67, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x70, 0x75, 0x62, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x0a,
	0x07, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x32, 0xf9, 0x03, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x60, 0x0a, 0x07, 0x47, 0x65, 0x74,
	0x48, 0x6f, 0x73, 0x74, 0x12, 0x28, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29,
	0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x60, 0x0a, 0x07, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x28, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x29, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x57, 0x0a,
	0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x25, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x6f,
	0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5d, 0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70,
	0x12, 0x27, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x6f, 0x72, 0x62, 0x69,
	0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x69, 0x0a, 0x0a, 0x4e, 0x65, 0x77, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x2b, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4e,
	0x65, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2c, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4e, 0x65, 0x77, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x88, 0x02, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x42, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6f, 0x72,
	0x62, 0x69, 0x73, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6f, 0x72, 0x62, 0x69, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x4f, 0x54,
	0x58, 0xaa, 0x02, 0x18, 0x4f, 0x72, 0x62, 0x69, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x18, 0x4f,
	0x72, 0x62, 0x69, 0x73, 0x5c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x5c, 0x56,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x24, 0x4f, 0x72, 0x62, 0x69, 0x73, 0x5c,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x1a, 0x4f, 0x72, 0x62, 0x69, 0x73, 0x3a, 0x3a, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_orbis_transport_v1alpha1_transport_proto_rawDescOnce sync.Once
	file_orbis_transport_v1alpha1_transport_proto_rawDescData = file_orbis_transport_v1alpha1_transport_proto_rawDesc
)

func file_orbis_transport_v1alpha1_transport_proto_rawDescGZIP() []byte {
	file_orbis_transport_v1alpha1_transport_proto_rawDescOnce.Do(func() {
		file_orbis_transport_v1alpha1_transport_proto_rawDescData = protoimpl.X.CompressGZIP(file_orbis_transport_v1alpha1_transport_proto_rawDescData)
	})
	return file_orbis_transport_v1alpha1_transport_proto_rawDescData
}

var file_orbis_transport_v1alpha1_transport_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_orbis_transport_v1alpha1_transport_proto_goTypes = []interface{}{
	(*GetHostRequest)(nil),     // 0: orbis.transport.v1alpha1.GetHostRequest
	(*GetHostResponse)(nil),    // 1: orbis.transport.v1alpha1.GetHostResponse
	(*SendRequest)(nil),        // 2: orbis.transport.v1alpha1.SendRequest
	(*SendResponse)(nil),       // 3: orbis.transport.v1alpha1.SendResponse
	(*ConnectRequest)(nil),     // 4: orbis.transport.v1alpha1.ConnectRequest
	(*ConnectResponse)(nil),    // 5: orbis.transport.v1alpha1.ConnectResponse
	(*GossipRequest)(nil),      // 6: orbis.transport.v1alpha1.GossipRequest
	(*GossipResponse)(nil),     // 7: orbis.transport.v1alpha1.GossipResponse
	(*NewMessageRequest)(nil),  // 8: orbis.transport.v1alpha1.NewMessageRequest
	(*NewMessageResponse)(nil), // 9: orbis.transport.v1alpha1.NewMessageResponse
	(*Node)(nil),               // 10: orbis.transport.v1alpha1.Node
	(*Message)(nil),            // 11: orbis.transport.v1alpha1.Message
	(*pb.PublicKey)(nil),       // 12: libp2p.crypto.v1.PublicKey
}
var file_orbis_transport_v1alpha1_transport_proto_depIdxs = []int32{
	10, // 0: orbis.transport.v1alpha1.GetHostResponse.node:type_name -> orbis.transport.v1alpha1.Node
	10, // 1: orbis.transport.v1alpha1.SendRequest.node:type_name -> orbis.transport.v1alpha1.Node
	11, // 2: orbis.transport.v1alpha1.SendRequest.message:type_name -> orbis.transport.v1alpha1.Message
	11, // 3: orbis.transport.v1alpha1.GossipRequest.message:type_name -> orbis.transport.v1alpha1.Message
	11, // 4: orbis.transport.v1alpha1.NewMessageResponse.message:type_name -> orbis.transport.v1alpha1.Message
	12, // 5: orbis.transport.v1alpha1.Node.public_key:type_name -> libp2p.crypto.v1.PublicKey
	0,  // 6: orbis.transport.v1alpha1.TransportService.GetHost:input_type -> orbis.transport.v1alpha1.GetHostRequest
	4,  // 7: orbis.transport.v1alpha1.TransportService.Connect:input_type -> orbis.transport.v1alpha1.ConnectRequest
	2,  // 8: orbis.transport.v1alpha1.TransportService.Send:input_type -> orbis.transport.v1alpha1.SendRequest
	6,  // 9: orbis.transport.v1alpha1.TransportService.Gossip:input_type -> orbis.transport.v1alpha1.GossipRequest
	8,  // 10: orbis.transport.v1alpha1.TransportService.NewMessage:input_type -> orbis.transport.v1alpha1.NewMessageRequest
	1,  // 11: orbis.transport.v1alpha1.TransportService.GetHost:output_type -> orbis.transport.v1alpha1.GetHostResponse
	5,  // 12: orbis.transport.v1alpha1.TransportService.Connect:output_type -> orbis.transport.v1alpha1.ConnectResponse
	3,  // 13: orbis.transport.v1alpha1.TransportService.Send:output_type -> orbis.transport.v1alpha1.SendResponse
	7,  // 14: orbis.transport.v1alpha1.TransportService.Gossip:output_type -> orbis.transport.v1alpha1.GossipResponse
	9,  // 15: orbis.transport.v1alpha1.TransportService.NewMessage:output_type -> orbis.transport.v1alpha1.NewMessageResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_orbis_transport_v1alpha1_transport_proto_init() }
func file_orbis_transport_v1alpha1_transport_proto_init() {
	if File_orbis_transport_v1alpha1_transport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetHostRequest); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetHostResponse); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRequest); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendResponse); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectRequest); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectResponse); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GossipRequest); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GossipResponse); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewMessageRequest); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewMessageResponse); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_orbis_transport_v1alpha1_transport_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_orbis_transport_v1alpha1_transport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orbis_transport_v1alpha1_transport_proto_goTypes,
		DependencyIndexes: file_orbis_transport_v1alpha1_transport_proto_depIdxs,
		MessageInfos:      file_orbis_transport_v1alpha1_transport_proto_msgTypes,
	}.Build()
	File_orbis_transport_v1alpha1_transport_proto = out.File
	file_orbis_transport_v1alpha1_transport_proto_rawDesc = nil
	file_orbis_transport_v1alpha1_transport_proto_goTypes = nil
	file_orbis_transport_v1alpha1_transport_proto_depIdxs = nil
}