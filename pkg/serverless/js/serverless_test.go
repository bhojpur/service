package js

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
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestByteSlice(t *testing.T) {
	vm := goja.New()
	// length: 16
	// [116 104 105 115 32 105 115 32 97 32 115 116 114 105 110 103]
	data := []byte("this is a string")
	t.Logf("data[%d]=%v\n", len(data), data)
	buf := vm.NewArrayBuffer(data)
	vm.Set("buf", buf)
	_, err := vm.RunString(`
	var a = new Uint8Array(buf);
	if (a.length !== 16 || a[0] !== 116 || a[1] !== 104 || a[15] !== 103) {
		throw new Error(a);
	}
	a[0]=84
	`)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data[%d]=%v\n", len(data), data)
	t.Logf("data=>string=%s\n", data)
	//
	ret, err := vm.RunString(`
	var b = Uint8Array.of(0xCC, 0xDD);
	b.buffer;
	`)
	if err != nil {
		t.Fatal(err)
	}
	buf1 := ret.Export().(goja.ArrayBuffer)
	data1 := buf1.Bytes()
	if len(data1) != 2 || data1[0] != 0xCC || data1[1] != 0xDD {
		t.Fatal(data1)
	}
}
func TestGoHandler(t *testing.T) {
}

func TestJSHandler(t *testing.T) {
	source := `
function handler(data) {
	log.Println(">> JS Handler Begin");
	var uint8buf = new Uint8Array(data); 
	// go method
	// var decodedString = svcsvr.ArrayBufferToString(data);
	// log.Printf(">> data.string: %s", decodedString);
	// js method
	decodedString = decode(uint8buf);
	log.Printf(">> data.string: %s", decodedString);
	// JSON
	var jsonData=JSON.parse(decodedString);
	log.Printf(">> data.JSON: %v", jsonData);
	log.Printf(">> data.Hex: %v", ab2hex(data));
	log.Println(">> JS Handler End");
	return {id:0x35, data:data}
}

// function DataID() []byte {
function dataID() {
    // return []byte {0x34}
    // return Uint8Array.of(0x34)
    // return [0x34]
	return [52,53]
}

function decode(uint8buf){
	 var encodedString = String.fromCharCode.apply(null,uint8buf);
	 decodedString = decodeURIComponent(escape(encodedString));
	 return decodedString;
}

// Uint8Array to Hex
function ab2hex(buffer) {
	let hexArr = Array.prototype.map.call(
		new Uint8Array(buffer),
		function (bit) { return ('00' + bit.toString(16)).slice(-2) }
	)
	return hexArr.join('')
}
`
	vm := goja.New()
	vm.Set("log", log.Default())
	vm.Set("bhojpur", BhojpurService{})
	var err error
	_, err = vm.RunString(source)
	if err != nil {
		t.Fatal(err)
	}
	// handler
	// data := vm.NewArrayBuffer([]byte("hello js"))
	// vm.Set("data", data)
	// var handlerFn func(goja.ArrayBuffer) goja.ArrayBuffer
	var handlerFn func(goja.ArrayBuffer) map[string]interface{}
	err = vm.ExportTo(vm.Get("handler"), &handlerFn)
	if err != nil {
		t.Fatal(err)
	}
	nd := NoiseData{
		Noise:   5.8,
		Time:    time.Now().UnixMilli(),
		From:    "Go->JS",
		Created: time.Now(),
	}
	t.Logf("source.Struct: %+v\n", nd)
	msg, _ := json.Marshal(nd)
	updata := vm.NewArrayBuffer(msg)
	t.Logf("source.Hex: %x\n", msg)
	data := handlerFn(updata)
	// t.Logf("JS Handler result: id=%v, data=%s\n", "id", data.Bytes())
	buf := data["data"].(goja.ArrayBuffer).Bytes()
	t.Logf("<< JS Handler result: id=%v, data=%s\n", data["id"], buf)
	var noise NoiseData
	json.Unmarshal(buf, &noise)
	t.Logf("<< JS Handler result: %+v\n", noise)
	// data id
	var dataIDFn func() []byte
	err = vm.ExportTo(vm.Get("dataID"), &dataIDFn)
	if err != nil {
		t.Fatal(err)
	}
	dataIDs := dataIDFn()
	t.Logf("<< JS DataID result: dataIDs=%v\n", dataIDs)

}

type NoiseData struct {
	Noise   float32   `json:"noise"` // Noise value
	Time    int64     `json:"time"`  // Timestamp (ms)
	From    string    `json:"from"`  // Source IP
	Created time.Time `json:created`
}

type BhojpurService struct{}

func (y BhojpurService) BytesToString(buf []byte) string {
	return string(buf)
}
func (y BhojpurService) ArrayBufferToString(buf goja.ArrayBuffer) string {
	return string(buf.Bytes())
}
