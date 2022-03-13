package main

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
	"math/rand"
	"time"

	svcsvr "github.com/bhojpur/service/pkg/engine"
)

type noiseData struct {
	Noise float32 `json:"noise"` // Noise value
	Time  int64   `json:"time"`  // Timestamp (ms)
	From  string  `json:"from"`  // Source IP
}

func main() {
	// connect source to the Bhojpur Service-Processor.
	source := svcsvr.NewSource("bhojpur-source", svcsvr.WithProcessorAddr("localhost:9140"))
	err := source.Connect()
	if err != nil {
		log.Printf("❌ Emit the data to Bhojpur Service-Processor failure with err: %v", err)
		return
	}

	defer source.Close()

	source.SetDataTag(0x33)
	// generate mock data and send it to Bhojpur Service-Processor in every 100 ms.
	generateAndSendData(source)
}

func generateAndSendData(stream svcsvr.Source) {
	for {
		// generate random data.
		data := noiseData{
			Noise: rand.New(rand.NewSource(time.Now().UnixNano())).Float32() * 200,
			Time:  time.Now().UnixNano() / int64(time.Millisecond),
			From:  "localhost",
		}

		// encode data via JSON codec.
		sendingBuf, _ := json.Marshal(data)

		// send data via QUIC stream.
		_, err := stream.Write(sendingBuf)
		if err != nil {
			log.Printf("❌ Emit %v to Bhojpur Service-Processor failure with err: %v", data, err)
		} else {
			log.Printf("✅ Emit %v to Bhojpur Service-Processor", data)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
