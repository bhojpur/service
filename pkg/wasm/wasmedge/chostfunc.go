package wasmedge

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

/*
#include <wasmedge/wasmedge.h>
// The gateway function
WasmEdge_Result wasmedgego_HostFuncInvoke(void *Func, void *Data,
                           WasmEdge_MemoryInstanceContext *MemCxt,
                           const WasmEdge_Value *Params, const uint32_t ParamLen,
                           WasmEdge_Value *Returns, const uint32_t ReturnLen)
{
    WasmEdge_Result wasmedgego_HostFuncInvokeImpl(
        void *Func, void *Data,
        WasmEdge_MemoryInstanceContext *MemCxt,
        const WasmEdge_Value *Params, const uint32_t ParamLen,
        WasmEdge_Value *Returns, const uint32_t ReturnLen);
    return wasmedgego_HostFuncInvokeImpl(Func, Data, MemCxt, Params, ParamLen, Returns, ReturnLen);
}
*/
import "C"
