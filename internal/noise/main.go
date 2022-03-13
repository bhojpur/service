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

	"github.com/bhojpur/service/internal/noise/env"
	"github.com/bhojpur/service/internal/noise/starter"
	"github.com/bhojpur/service/internal/noise/utils"
	svcsvr "github.com/bhojpur/service/pkg/engine"
)

var (
	processorAddr = env.GetString("BHOJPUR_SOURCE_MQTT_PROCESSOR_ADDR", "localhost:9140")
	brokerAddr    = env.GetString("BHOJPUR_SOURCE_MQTT_BROKER_ADDR", "0.0.0.0:1883")
	source        svcsvr.Source
)

type NoiseData struct {
	Noise float32 `json:"noise"` // Noise value
	Time  int64   `json:"time"`  // Timestamp (ms)
	From  string  `json:"from"`  // Source IP
}

func main() {
	// connect to Bhojpur Service-Processor.
	source = svcsvr.NewSource("bhojpur-source", svcsvr.WithProcessorAddr(processorAddr))
	err := source.Connect()
	if err != nil {
		log.Printf("[source] ❌ Connect to Bhojpur Service-Processor %s failure with err: %v", processorAddr, err)
		return
	}

	defer source.Close()

	// set the data tag.
	source.SetDataTag(0x33)

	// start a new MQTT Broker.
	starter.NewBrokerSimply(brokerAddr, "NOISE").
		Run(handler)
}

func handler(topic string, payload []byte) {
	log.Printf("receive: topic=%v, payload=%v\n", topic, string(payload))

	// get data from MQTT
	var raw map[string]int32
	err := json.Unmarshal(payload, &raw)
	if err != nil {
		log.Printf("Unmarshal payload error:%v", err)
	}

	noise := float32(raw["noise"])
	data := NoiseData{Noise: noise, Time: utils.Now(), From: utils.IpAddr()}
	sendingBuf, _ := json.Marshal(data)

	// send data to Bhojpur Service-Processor.
	_, err = source.Write(sendingBuf)
	if err != nil {
		log.Printf("source.Write error: %v, sendingBuf=%#x\n", err, sendingBuf)
	}

	log.Printf("write: sendingBuf=%v\n", utils.FormatBytes(sendingBuf))
}
