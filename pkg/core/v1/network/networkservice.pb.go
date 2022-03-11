// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: pkg/core/v1/network/networkservice.proto

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

// It contains the core Network Service API definitions for external consumption
// via gRPC protobufs.

package network

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NetworkServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Connection           *Connection  `protobuf:"bytes,1,opt,name=connection,proto3" json:"connection,omitempty"`
	MechanismPreferences []*Mechanism `protobuf:"bytes,2,rep,name=mechanism_preferences,json=mechanismPreferences,proto3" json:"mechanism_preferences,omitempty"`
}

func (x *NetworkServiceRequest) Reset() {
	*x = NetworkServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_network_networkservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkServiceRequest) ProtoMessage() {}

func (x *NetworkServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_network_networkservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkServiceRequest.ProtoReflect.Descriptor instead.
func (*NetworkServiceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_network_networkservice_proto_rawDescGZIP(), []int{0}
}

func (x *NetworkServiceRequest) GetConnection() *Connection {
	if x != nil {
		return x.Connection
	}
	return nil
}

func (x *NetworkServiceRequest) GetMechanismPreferences() []*Mechanism {
	if x != nil {
		return x.MechanismPreferences
	}
	return nil
}

var File_pkg_core_v1_network_networkservice_proto protoreflect.FileDescriptor

var file_pkg_core_v1_network_networkservice_proto_rawDesc = []byte{
	0x0a, 0x28, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x76, 0x31, 0x2e, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x15, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x31, 0x2e, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4a, 0x0a, 0x15, 0x6d,
	0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73, 0x6d, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x31, 0x2e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x4d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73,
	0x6d, 0x52, 0x14, 0x6d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73, 0x6d, 0x50, 0x72, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x32, 0x8f, 0x01, 0x0a, 0x0e, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x2e, 0x76, 0x31, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x76, 0x31, 0x2e, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x37, 0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x16, 0x2e, 0x76, 0x31, 0x2e, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x3b, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_core_v1_network_networkservice_proto_rawDescOnce sync.Once
	file_pkg_core_v1_network_networkservice_proto_rawDescData = file_pkg_core_v1_network_networkservice_proto_rawDesc
)

func file_pkg_core_v1_network_networkservice_proto_rawDescGZIP() []byte {
	file_pkg_core_v1_network_networkservice_proto_rawDescOnce.Do(func() {
		file_pkg_core_v1_network_networkservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_core_v1_network_networkservice_proto_rawDescData)
	})
	return file_pkg_core_v1_network_networkservice_proto_rawDescData
}

var file_pkg_core_v1_network_networkservice_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_core_v1_network_networkservice_proto_goTypes = []interface{}{
	(*NetworkServiceRequest)(nil), // 0: v1.network.NetworkServiceRequest
	(*Connection)(nil),            // 1: v1.network.Connection
	(*Mechanism)(nil),             // 2: v1.network.Mechanism
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_pkg_core_v1_network_networkservice_proto_depIdxs = []int32{
	1, // 0: v1.network.NetworkServiceRequest.connection:type_name -> v1.network.Connection
	2, // 1: v1.network.NetworkServiceRequest.mechanism_preferences:type_name -> v1.network.Mechanism
	0, // 2: v1.network.NetworkService.Request:input_type -> v1.network.NetworkServiceRequest
	1, // 3: v1.network.NetworkService.Close:input_type -> v1.network.Connection
	1, // 4: v1.network.NetworkService.Request:output_type -> v1.network.Connection
	3, // 5: v1.network.NetworkService.Close:output_type -> google.protobuf.Empty
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_core_v1_network_networkservice_proto_init() }
func file_pkg_core_v1_network_networkservice_proto_init() {
	if File_pkg_core_v1_network_networkservice_proto != nil {
		return
	}
	file_pkg_core_v1_network_connection_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pkg_core_v1_network_networkservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkServiceRequest); i {
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
			RawDescriptor: file_pkg_core_v1_network_networkservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_core_v1_network_networkservice_proto_goTypes,
		DependencyIndexes: file_pkg_core_v1_network_networkservice_proto_depIdxs,
		MessageInfos:      file_pkg_core_v1_network_networkservice_proto_msgTypes,
	}.Build()
	File_pkg_core_v1_network_networkservice_proto = out.File
	file_pkg_core_v1_network_networkservice_proto_rawDesc = nil
	file_pkg_core_v1_network_networkservice_proto_goTypes = nil
	file_pkg_core_v1_network_networkservice_proto_depIdxs = nil
}
