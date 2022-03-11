package mark

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
	"fmt"

	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

// Tag represents the Tag of TLV
// MSB used to represent the packet type, 0x80 means a node packet, otherwise is a primitive packet
// Low 7 bits represent Sequence ID, like `key` in JSON format
type Tag struct {
	raw byte
}

// IsNode returns true is MSB is 1.
func (t *Tag) IsNode() bool {
	return t.raw&utils.MSB == utils.MSB
}

// SeqID get the sequence ID, as key in JSON format
func (t *Tag) SeqID() byte {
	//return t.raw & packetutils.DropMSB
	return t.raw & utils.DropMSBArrayFlag
}

// String return string description
func (t *Tag) String() string {
	return fmt.Sprintf("Tag: raw=%4b, SeqID=%v", t.raw, t.SeqID())
}

// NewTag create a NodePacket Tag field
func NewTag(b byte) *Tag {
	return &Tag{raw: b}
}

// IsSlice determine if the current node is a Slice
func (t *Tag) IsSlice() bool {
	return t.raw&utils.SliceFlag == utils.SliceFlag
}

// Raw return the original byte
func (t *Tag) Raw() byte {
	return t.raw
}
