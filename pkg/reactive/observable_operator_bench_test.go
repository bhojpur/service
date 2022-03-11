package reactive

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
	"context"
	"testing"
	"time"
)

const (
	benchChannelCap            = 1000
	benchNumberOfElementsSmall = 1000
	ioPool                     = 32
)

func Benchmark_Range_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Range_Serialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool(), WithBufferedChannel(benchChannelCap)).
			Serialize(0, func(i interface{}) int {
				return i.(int)
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Range_OptionSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool(), WithBufferedChannel(benchChannelCap), Serialize(func(i interface{}) int {
				return i.(int)
			}))
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Reduce_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				if a, ok := acc.(int); ok {
					if b, ok := elem.(int); ok {
						return a + b, nil
					}
				} else {
					return elem.(int), nil
				}
				return 0, errFoo
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Reduce_Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Reduce(func(_ context.Context, acc, elem interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				if a, ok := acc.(int); ok {
					if b, ok := elem.(int); ok {
						return a + b, nil
					}
				} else {
					return elem.(int), nil
				}
				return 0, errFoo
			}, WithPool(ioPool))
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Map_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			})
		b.StartTimer()
		<-obs.Run()
	}
}

func Benchmark_Map_Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		obs := Range(0, benchNumberOfElementsSmall, WithBufferedChannel(benchChannelCap)).
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				// Simulate a blocking IO call
				time.Sleep(5 * time.Millisecond)
				return i, nil
			}, WithCPUPool())
		b.StartTimer()
		<-obs.Run()
	}
}
