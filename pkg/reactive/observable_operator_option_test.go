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

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func Test_Observable_Option_WithOnErrorStrategy_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			if i == 2 {
				return nil, errFoo
			}
			return i, nil
		}, WithErrorStrategy(ContinueOnError))
	Assert(context.Background(), t, obs, HasItems(1, 3), HasError(errFoo))
}

func Test_Observable_Option_WithOnErrorStrategy_Propagate(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, 1, 2, 3).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			if i == 1 {
				return nil, errFoo
			}
			return i, nil
		}).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			if i == 2 {
				return nil, errBar
			}
			return i, nil
		}, WithErrorStrategy(ContinueOnError))
	Assert(context.Background(), t, obs, HasItems(3), HasErrors(errFoo, errBar))
}

func Test_Observable_Option_SimpleCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := Just(1)(WithBufferedChannel(5)).Observe()
	assert.Equal(t, 5, cap(ch))
}

func Test_Observable_Option_ComposedCapacity(t *testing.T) {
	defer goleak.VerifyNone(t)
	obs1 := Just(1)().Map(func(_ context.Context, _ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(11))
	obs2 := obs1.Map(func(_ context.Context, _ interface{}) (interface{}, error) {
		return 1, nil
	}, WithBufferedChannel(12))

	assert.Equal(t, 11, cap(obs1.Observe()))
	assert.Equal(t, 12, cap(obs2.Observe()))
}

func Test_Observable_Option_ContextPropagation(t *testing.T) {
	defer goleak.VerifyNone(t)
	expectedCtx := context.Background()
	var gotCtx context.Context
	<-Just(1)().Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		gotCtx = ctx
		return i, nil
	}, WithContext(expectedCtx)).Run()
	assert.Equal(t, expectedCtx, gotCtx)
}

// FIXME
//func Test_Observable_Option_Serialize(t *testing.T) {
//	defer goleak.VerifyNone(t)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	idx := 1
//	<-testObservable(ctx, 1, 3, 2, 6, 4, 5).Map(func(_ context.Context, i interface{}) (interface{}, error) {
//		return i, nil
//	}, WithBufferedChannel(10), WithCPUPool(), WithContext(ctx), Serialize(func(i interface{}) int {
//		return i.(int)
//	})).DoOnNext(func(i interface{}) {
//		v := i.(int)
//		if v != idx {
//			assert.FailNow(t, "not sequential", "expected=%d, got=%d", idx, v)
//		}
//		idx++
//	})
//}

func Test_Observable_Option_Serialize_Range(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	idx := 0
	<-Range(0, 10000).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}, WithBufferedChannel(10), WithCPUPool(), WithContext(ctx), Serialize(func(i interface{}) int {
		return i.(int)
	})).DoOnNext(func(i interface{}) {
		v := i.(int)
		if v != idx {
			assert.FailNow(t, "not sequential", "expected=%d, got=%d", idx, v)
		}
		idx++
	})
}

func Test_Observable_Option_Serialize_SingleElement(t *testing.T) {
	defer goleak.VerifyNone(t)
	idx := 0
	<-Just(0)().Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}, WithBufferedChannel(10), WithCPUPool(), Serialize(func(i interface{}) int {
		return i.(int)
	})).DoOnNext(func(i interface{}) {
		v := i.(int)
		if v != idx {
			assert.FailNow(t, "not sequential", "expected=%d, got=%d", idx, v)
		}
		idx++
	})
}

func Test_Observable_Option_Serialize_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := testObservable(ctx, errFoo, 2, 3, 4).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}, WithBufferedChannel(10), WithCPUPool(), WithContext(ctx), Serialize(func(i interface{}) int {
		return i.(int)
	}))
	Assert(context.Background(), t, obs, IsEmpty(), HasError(errFoo))
}
