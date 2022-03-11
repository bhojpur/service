// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: pkg/core/v1/keyed/keyedserver.proto

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package keyed

import (
	_ "github.com/gogo/protobuf/gogoproto"
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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         *uint64 `protobuf:"varint,1,opt,name=ID,proto3,oneof" json:"ID,omitempty"`
	Method     *string `protobuf:"bytes,2,opt,name=Method,proto3,oneof" json:"Method,omitempty"`
	Path       *string `protobuf:"bytes,3,opt,name=Path,proto3,oneof" json:"Path,omitempty"`
	Val        *string `protobuf:"bytes,4,opt,name=Val,proto3,oneof" json:"Val,omitempty"`
	Dir        *bool   `protobuf:"varint,5,opt,name=Dir,proto3,oneof" json:"Dir,omitempty"`
	PrevValue  *string `protobuf:"bytes,6,opt,name=PrevValue,proto3,oneof" json:"PrevValue,omitempty"`
	PrevIndex  *uint64 `protobuf:"varint,7,opt,name=PrevIndex,proto3,oneof" json:"PrevIndex,omitempty"`
	PrevExist  *bool   `protobuf:"varint,8,opt,name=PrevExist,proto3,oneof" json:"PrevExist,omitempty"`
	Expiration *int64  `protobuf:"varint,9,opt,name=Expiration,proto3,oneof" json:"Expiration,omitempty"`
	Wait       *bool   `protobuf:"varint,10,opt,name=Wait,proto3,oneof" json:"Wait,omitempty"`
	Since      *uint64 `protobuf:"varint,11,opt,name=Since,proto3,oneof" json:"Since,omitempty"`
	Recursive  *bool   `protobuf:"varint,12,opt,name=Recursive,proto3,oneof" json:"Recursive,omitempty"`
	Sorted     *bool   `protobuf:"varint,13,opt,name=Sorted,proto3,oneof" json:"Sorted,omitempty"`
	Quorum     *bool   `protobuf:"varint,14,opt,name=Quorum,proto3,oneof" json:"Quorum,omitempty"`
	Time       *int64  `protobuf:"varint,15,opt,name=Time,proto3,oneof" json:"Time,omitempty"`
	Stream     *bool   `protobuf:"varint,16,opt,name=Stream,proto3,oneof" json:"Stream,omitempty"`
	Refresh    *bool   `protobuf:"varint,17,opt,name=Refresh,proto3,oneof" json:"Refresh,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_keyed_keyedserver_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetID() uint64 {
	if x != nil && x.ID != nil {
		return *x.ID
	}
	return 0
}

func (x *Request) GetMethod() string {
	if x != nil && x.Method != nil {
		return *x.Method
	}
	return ""
}

func (x *Request) GetPath() string {
	if x != nil && x.Path != nil {
		return *x.Path
	}
	return ""
}

func (x *Request) GetVal() string {
	if x != nil && x.Val != nil {
		return *x.Val
	}
	return ""
}

func (x *Request) GetDir() bool {
	if x != nil && x.Dir != nil {
		return *x.Dir
	}
	return false
}

func (x *Request) GetPrevValue() string {
	if x != nil && x.PrevValue != nil {
		return *x.PrevValue
	}
	return ""
}

func (x *Request) GetPrevIndex() uint64 {
	if x != nil && x.PrevIndex != nil {
		return *x.PrevIndex
	}
	return 0
}

func (x *Request) GetPrevExist() bool {
	if x != nil && x.PrevExist != nil {
		return *x.PrevExist
	}
	return false
}

func (x *Request) GetExpiration() int64 {
	if x != nil && x.Expiration != nil {
		return *x.Expiration
	}
	return 0
}

func (x *Request) GetWait() bool {
	if x != nil && x.Wait != nil {
		return *x.Wait
	}
	return false
}

func (x *Request) GetSince() uint64 {
	if x != nil && x.Since != nil {
		return *x.Since
	}
	return 0
}

func (x *Request) GetRecursive() bool {
	if x != nil && x.Recursive != nil {
		return *x.Recursive
	}
	return false
}

func (x *Request) GetSorted() bool {
	if x != nil && x.Sorted != nil {
		return *x.Sorted
	}
	return false
}

func (x *Request) GetQuorum() bool {
	if x != nil && x.Quorum != nil {
		return *x.Quorum
	}
	return false
}

func (x *Request) GetTime() int64 {
	if x != nil && x.Time != nil {
		return *x.Time
	}
	return 0
}

func (x *Request) GetStream() bool {
	if x != nil && x.Stream != nil {
		return *x.Stream
	}
	return false
}

func (x *Request) GetRefresh() bool {
	if x != nil && x.Refresh != nil {
		return *x.Refresh
	}
	return false
}

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID    *uint64 `protobuf:"varint,1,opt,name=NodeID,proto3,oneof" json:"NodeID,omitempty"`
	ClusterID *uint64 `protobuf:"varint,2,opt,name=ClusterID,proto3,oneof" json:"ClusterID,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_keyed_keyedserver_proto_rawDescGZIP(), []int{1}
}

func (x *Metadata) GetNodeID() uint64 {
	if x != nil && x.NodeID != nil {
		return *x.NodeID
	}
	return 0
}

func (x *Metadata) GetClusterID() uint64 {
	if x != nil && x.ClusterID != nil {
		return *x.ClusterID
	}
	return 0
}

var File_pkg_core_v1_keyed_keyedserver_proto protoreflect.FileDescriptor

var file_pkg_core_v1_keyed_keyedserver_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x65,
	0x79, 0x65, 0x64, 0x2f, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x76, 0x31, 0x2e, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x1a,
	0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x67, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x97,
	0x06, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x02, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x00, 0x52, 0x02,
	0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x01, 0x52, 0x06, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x02, 0x52, 0x04,
	0x50, 0x61, 0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x03, 0x56, 0x61, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x03, 0x52, 0x03, 0x56, 0x61,
	0x6c, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x03, 0x44, 0x69, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x04, 0x52, 0x03, 0x44, 0x69, 0x72, 0x88, 0x01,
	0x01, 0x12, 0x27, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x76, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x05, 0x52, 0x09, 0x50, 0x72,
	0x65, 0x76, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x09, 0x50, 0x72,
	0x65, 0x76, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x42, 0x04, 0xc8,
	0xde, 0x1f, 0x00, 0x48, 0x06, 0x52, 0x09, 0x50, 0x72, 0x65, 0x76, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x76, 0x45, 0x78, 0x69, 0x73, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x01, 0x48, 0x07, 0x52, 0x09,
	0x50, 0x72, 0x65, 0x76, 0x45, 0x78, 0x69, 0x73, 0x74, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0a,
	0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x08, 0x52, 0x0a, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x04, 0x57, 0x61, 0x69, 0x74, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x09, 0x52, 0x04, 0x57,
	0x61, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x05, 0x53, 0x69, 0x6e, 0x63, 0x65, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x04, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x0a, 0x52, 0x05, 0x53,
	0x69, 0x6e, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x75, 0x72,
	0x73, 0x69, 0x76, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00,
	0x48, 0x0b, 0x52, 0x09, 0x52, 0x65, 0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x21, 0x0a, 0x06, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08,
	0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x0c, 0x52, 0x06, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x06, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x0d, 0x52, 0x06, 0x51, 0x75, 0x6f,
	0x72, 0x75, 0x6d, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x03, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x0e, 0x52, 0x04, 0x54, 0x69,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18,
	0x10, 0x20, 0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x0f, 0x52, 0x06, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x07, 0x52, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x01, 0x48,
	0x10, 0x52, 0x07, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a,
	0x03, 0x5f, 0x49, 0x44, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x50, 0x61, 0x74, 0x68, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x56, 0x61, 0x6c,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x44, 0x69, 0x72, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x50, 0x72, 0x65,
	0x76, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x50, 0x72, 0x65, 0x76, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x50, 0x72, 0x65, 0x76, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x45, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x57, 0x61, 0x69, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x53,
	0x69, 0x6e, 0x63, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x52, 0x65, 0x63, 0x75, 0x72, 0x73, 0x69,
	0x76, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x54, 0x69, 0x6d,
	0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x22, 0x6f, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x48, 0x00, 0x52, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00,
	0x48, 0x01, 0x52, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x44, 0x42, 0x44, 0x5a, 0x32, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x3b, 0x6b, 0x65, 0x79, 0x65, 0x64, 0xc8,
	0xe2, 0x1e, 0x01, 0xe0, 0xe2, 0x1e, 0x01, 0xd0, 0xe2, 0x1e, 0x01, 0xc8, 0xe1, 0x1e, 0x00, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_core_v1_keyed_keyedserver_proto_rawDescOnce sync.Once
	file_pkg_core_v1_keyed_keyedserver_proto_rawDescData = file_pkg_core_v1_keyed_keyedserver_proto_rawDesc
)

func file_pkg_core_v1_keyed_keyedserver_proto_rawDescGZIP() []byte {
	file_pkg_core_v1_keyed_keyedserver_proto_rawDescOnce.Do(func() {
		file_pkg_core_v1_keyed_keyedserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_core_v1_keyed_keyedserver_proto_rawDescData)
	})
	return file_pkg_core_v1_keyed_keyedserver_proto_rawDescData
}

var file_pkg_core_v1_keyed_keyedserver_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_core_v1_keyed_keyedserver_proto_goTypes = []interface{}{
	(*Request)(nil),  // 0: v1.keyed.Request
	(*Metadata)(nil), // 1: v1.keyed.Metadata
}
var file_pkg_core_v1_keyed_keyedserver_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_core_v1_keyed_keyedserver_proto_init() }
func file_pkg_core_v1_keyed_keyedserver_proto_init() {
	if File_pkg_core_v1_keyed_keyedserver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
	file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_pkg_core_v1_keyed_keyedserver_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_core_v1_keyed_keyedserver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_core_v1_keyed_keyedserver_proto_goTypes,
		DependencyIndexes: file_pkg_core_v1_keyed_keyedserver_proto_depIdxs,
		MessageInfos:      file_pkg_core_v1_keyed_keyedserver_proto_msgTypes,
	}.Build()
	File_pkg_core_v1_keyed_keyedserver_proto = out.File
	file_pkg_core_v1_keyed_keyedserver_proto_rawDesc = nil
	file_pkg_core_v1_keyed_keyedserver_proto_goTypes = nil
	file_pkg_core_v1_keyed_keyedserver_proto_depIdxs = nil
}
