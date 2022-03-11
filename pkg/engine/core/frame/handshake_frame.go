package frame

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

import (
	"github.com/bhojpur/service/pkg/engine/codec"
)

// HandshakeFrame is a Bhojpur Service encoded.
type HandshakeFrame struct {
	// Name is client name
	Name string
	// ClientType represents client type (source or sfn)
	ClientType byte
	// ObserveDataTags are the client data tag list.
	ObserveDataTags []byte
	// auth
	authType    byte
	authPayload []byte
	// app id
	appID string
}

// NewHandshakeFrame creates a new HandshakeFrame.
func NewHandshakeFrame(name string, clientType byte, observeDataTags []byte, appID string, authType byte, authPayload []byte) *HandshakeFrame {
	return &HandshakeFrame{
		Name:            name,
		ClientType:      clientType,
		ObserveDataTags: observeDataTags,
		appID:           appID,
		authType:        authType,
		authPayload:     authPayload,
	}
}

// Type gets the type of Frame.
func (h *HandshakeFrame) Type() Type {
	return TagOfHandshakeFrame
}

// Encode to Bhojpur Service encoding.
func (h *HandshakeFrame) Encode() []byte {
	// name
	nameBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeName)))
	nameBlock.SetStringValue(h.Name)
	// type
	typeBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeType)))
	typeBlock.SetBytesValue([]byte{h.ClientType})
	// observe data tags
	observeDataTagsBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeObserveDataTags)))
	observeDataTagsBlock.SetBytesValue(h.ObserveDataTags)
	// app id
	appIDBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeAppID)))
	appIDBlock.SetStringValue(h.appID)
	// auth
	authTypeBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeAuthType)))
	authTypeBlock.SetBytesValue([]byte{h.authType})
	authPayloadBlock := codec.NewPrimitivePacketEncoder(int(byte(TagOfHandshakeAuthPayload)))
	authPayloadBlock.SetBytesValue(h.authPayload)
	// handshake frame
	handshake := codec.NewNodePacketEncoder(int(byte(h.Type())))
	handshake.AddPrimitivePacket(nameBlock)
	handshake.AddPrimitivePacket(typeBlock)
	handshake.AddPrimitivePacket(observeDataTagsBlock)
	handshake.AddPrimitivePacket(appIDBlock)
	handshake.AddPrimitivePacket(authTypeBlock)
	handshake.AddPrimitivePacket(authPayloadBlock)

	return handshake.Encode()
}

// DecodeToHandshakeFrame decodes Bhojpur Service encoded bytes to HandshakeFrame.
func DecodeToHandshakeFrame(buf []byte) (*HandshakeFrame, error) {
	//node := codec.NodePacket{}
	node, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}

	handshake := &HandshakeFrame{}

	// name
	nameBlock := node.PrimitivePackets[int(byte(TagOfHandshakeName))]
	name, err := nameBlock.ToUTF8String()
	if err != nil {
		return nil, err
	}
	handshake.Name = name

	// type
	typeBlock := node.PrimitivePackets[int(byte(TagOfHandshakeType))]
	clientType := typeBlock.ToBytes()
	handshake.ClientType = clientType[0]

	// observe data tag
	observeDataTagsBlock := node.PrimitivePackets[int(byte(TagOfHandshakeObserveDataTags))]
	handshake.ObserveDataTags = observeDataTagsBlock.ToBytes()

	// app id
	appIDBlock := node.PrimitivePackets[int(byte(TagOfHandshakeAppID))]
	appID, err := appIDBlock.ToUTF8String()
	if err != nil {
		return nil, err
	}
	handshake.appID = appID

	// auth
	authTypeBlock := node.PrimitivePackets[int(byte(TagOfHandshakeAuthType))]
	authType := authTypeBlock.ToBytes()
	handshake.authType = authType[0]

	// auth
	authPayloadBlock := node.PrimitivePackets[int(byte(TagOfHandshakeAuthPayload))]
	authPayload := authPayloadBlock.ToBytes()
	handshake.authPayload = authPayload

	return handshake, nil
}

func (h *HandshakeFrame) AuthType() byte {
	return h.authType
}

func (h *HandshakeFrame) AuthPayload() []byte {
	return h.authPayload
}

func (h *HandshakeFrame) AppID() string {
	return h.appID
}
