# Bhojpur Service - Serverless Framework

It is a streaming `serverless framework` for building low-latency Edge Computing applications. It is
built on top of `QUIC` Transport Protocol and `Functional Reactive Programming` interface. It makes
real-time data processing reliable, secure, and easy.

## Key Features

|     | **Features**|
| --- | ----------------------------------------------------------------------------------|
| ⚡️  | **Low-latency** Guaranteed by implementing atop QUIC [QUIC](https://datatracker.ietf.org/wg/quic/documents/) |
| 🔐  | **Security** TLS v1.3 on every data packet by design |
| 📱  | **5G/WiFi-6** Reliable networking in Celluar/Wireless |
| 🌎  | **Distributed Cloud Computing** EdgeMesh native architecture makes your services close to end users |
| 📸  | **Event-First** Architecture leverages serverless services to be event driven and elastic  |
| 🦖  | **Streaming Serverless** Write only a few lines of code to build applications and microservices |
| 🚀  | **Codec** a faster than real-time codec |
| 📨  | **Reactive** the core engine powered by a stream processing framework |

## 🚀 Getting Started

### Pre-requisites

Firstly, you must install [Go](https://golang.org/doc/install)

### 1. Install Bhojpur Service CLI

#### Build from Source Code

```bash
$ go install github.com/bhojpur/service@latest
```

#### Verify if Bhojpur Service CLI is installed successfully

```bash
$ svcutl -v

Bhojpur Service CLI version: v1.0.0
```

### 2. Create your Stream Function

```bash
$ svcutl init my-app-demo

⌛  Initializing the Bhojpur Service stream function...
✅  Congratulations! You have initialized the stream function successfully.
ℹ️  You can enjoy the Bhojpur Service stream function using following command: 
ℹ️   	DEV: 	svcutl dev -n Noise my-app-demo/app.go
ℹ️   	PROD: 	Firstly, run source application (e.g., go run internal/source/main.go)
	Secondly, svcutl run -n Noise -u localhost:9140 my-app-demo/app.go

$ cd my-app-demo
```

The Bhojpur Service `CLI` will automatically create the `app.go`. For example:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rx "github.com/bhojpur/service/pkg/reactive"
)

// NoiseData represents the structure of data
type NoiseData struct {
	Noise float32 `json:"noise"` // Noise value
	Time  int64   `json:"time"`  // Timestamp (ms)
	From  string  `json:"from"`  // Source IP
}

var echo = func(_ context.Context, i interface{}) (interface{}, error) {
	value := i.(*NoiseData)
	value.Noise = value.Noise / 10
	rightNow := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(fmt.Sprintf("[%s] %d > value: %f ⚡️=%dms", value.From, value.Time, value.Noise, rightNow-value.Time))
	return value.Noise, nil
}

// Handler will handle data in Reactive Stream way
func Handler(rxstream rx.Stream) rx.Stream {
	stream := rxstream.
		Unmarshal(json.Unmarshal, func() interface{} { return &NoiseData{} }).
		Debounce(50).
		Map(echo).
		StdOut()

	return stream
}

func DataTags() []byte {
	return []byte{0x33}
}
```

### 3. Run the Service Processor

Create a sample `workflow.yaml` file. Please note the stream function name (i.e., `Noise`).

```yaml
name: Service
host: localhost
port: 9140
functions:
  - name: Noise
  - name: MockDB
```

Now, run the following command in a new Terminal window.

```bash
$ svcutl serve --config workflow.yaml

Using config file: workflow.yaml
ℹ️   Running Bhojpur Service-Processor...
2018-03-26 21:07:18.548	[bhojpur:processor] Listening SIGUSR1, SIGUSR2, SIGTERM/SIGINT...
2018-03-26 21:07:18.551	[bhojpur:server] ✅ [Service] Bhojpur Service listening on: 127.0.0.1:9140, MODE: DEVELOPMENT, QUIC: [v1 draft-29], AUTH: [None]
2018-03-26 21:08:33.471	[bhojpur:server] ❤️  <Stream Function> [::Noise](127.0.0.1:55304) is connected!
2018-03-26 21:08:53.318	[bhojpur:server] ❤️  <Source> [::bhojpur-source](127.0.0.1:60306) is connected!
2018-03-26 21:09:35.006	ERROR	[bhojpur:server]  [ERR] timeout: no recent network activity
2018-03-26 21:09:35.006	[bhojpur:server] 💔 [::bhojpur-source](127.0.0.1:60306) close the Client connection
```

### 4. Build and Run Stream Function

Run `svcutl dev` or `svcutl run` command from the terminal. You will see the following messages:

```bash
$ svcutl run -n Noise -u localhost:9140 my-app-demo/app.go

ℹ️  Bhojpur Service stream function filename: my-app-demo/app.go
⌛  Create Bhojpur Service stream function instance...
ℹ️  Starting the Bhojpur Service stream function instance with Name: Noise. Host: localhost. Port: 9140.
⌛  Bhojpur Service stream function building...
✅  Success! Bhojpur Service stream function build.
ℹ️     Bhojpur Service stream function is running...
ℹ️     Run Go serverless: /Users/bhojpur/my-app-demo/sl.basm
2018-03-26 21:08:33.465	[core:client] use credential: [None]
2018-03-26 21:08:33.470	[core:client] ❤️  [Noise]([::]:55304) is connected to Bhojpur Service-Processor localhost:9140
2018-03-26 21:08:33.470	Reactive Stream handler is running...
[localhost] 1637028164050 > value: 6.575044 ⚡️=9ms
[StdOut]:  6.5750437
[localhost] 1637028164151 > value: 10.076103 ⚡️=5ms
[StdOut]:  10.076103
```

Congratulations! You have done your first Bhojpur Service stream function. Please note that `Noise` name
should be available as a function name in `workflow.yaml` file.

### 5. Run the Data Source

You must run the data feed. For example

```bash
$ go run ./internal/source/main.go

2018-03-26 21:00:27.286	[bhojpur:client] use credential: [None]
2018-03-26 21:00:27.290	[bhojpur:client] ❤️  [bhojpur-source]([::]:61746) is connected to Bhojpur Service-Processor localhost:9140
2018-03-26 21:00:27.290	✅ Emit {128.39642 1647185427290 localhost} to Bhojpur Service-Processor
2018-03-26 21:00:27.792	✅ Emit {58.476276 1647185427792 localhost} to Bhojpur Service-Processor
```

Now, check the Bhojpur Service-Processor terminal window for activities going on.

## 🧩 Interoperability

### Input Data/Event Sources

+ Connect EMQ X Broker to Bhojpur Service
+ Connect MQTT to Bhojpur Service

### Stream Functions

+ Write a Stream Function with WebAssembly

### Output Connectors

+ Connect to Graph database to store post-processed result the serverless way
+ Connect to Time Series database to store post-processed result
+ Connect to Big Data platform to store post-processed result

## 🎯 Focus on computing out-of Data Center

- IoT/IIoT/AIoT
- Latency-sensitive applications.
- Networking situation with packet loss or high latency.
- Handling continuous high frequency generated data with stream-processing.
- Building Complex systems with Streaming-Serverless architecture.

## 🌟 Why Bhojpur Service?

- Based on the QUIC (i.e. Quick UDP Internet Connection) protocol for data transmission, which uses
the user datagram protocol (UDP) as its basis instead of the transmission control protocol (TCP).
It significantly improves the overall stability and throughput of data transmission, especially
cellular networks, such as: 5G Mobile.
- A self-developed `codec` that optimizes decoding performance.
- Based on Stream Computing paradigm, which improves speed and accuracy when dealing with data handling
and analysis; simplifies the complexity of stream-oriented programming.
- Secure-by-default from transport protocol perspective.
