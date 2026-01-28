// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || windows

package purego

import "unsafe"

// FastFuncUint1 calls a C function with signature: uint foo(void* arg).
type FastFuncUint1 = func(uintptr) uint

// NewFastFuncUint1 returns a zero-allocation callable for a C function
// with signature uint foo(void* arg).
func NewFastFuncUint1(sym uintptr) FastFuncUint1 {
	return func(a1 uintptr) uint {
		return uint(fastCall1(sym, a1))
	}
}

// FastFuncPtr1 calls a C function with signature: void* foo(void* arg).
type FastFuncPtr1 = func(uintptr) unsafe.Pointer

// NewFastFuncPtr1 returns a zero-allocation callable for a C function
// with signature void* foo(void* arg).
func NewFastFuncPtr1(sym uintptr) FastFuncPtr1 {
	return func(a1 uintptr) unsafe.Pointer {
		r := fastCall1(sym, a1)
		return *(*unsafe.Pointer)(unsafe.Pointer(&r))
	}
}

// FastFuncStatus3 calls a C function with signature: int foo(void* res, void* a, void* stream).
type FastFuncStatus3 = func(unsafe.Pointer, uintptr, uintptr) int32

// NewFastFuncStatus3 returns a zero-allocation callable for a C function
// with signature int foo(void* res, void* a, void* stream).
func NewFastFuncStatus3(sym uintptr) FastFuncStatus3 {
	return func(res unsafe.Pointer, a, stream uintptr) int32 {
		return int32(fastCall3(sym, uintptr(res), a, stream))
	}
}

// FastFuncStatus4 calls a C function with signature: int foo(void* res, void* a, void* b, void* stream).
type FastFuncStatus4 = func(unsafe.Pointer, uintptr, uintptr, uintptr) int32

// NewFastFuncStatus4 returns a zero-allocation callable for a C function
// with signature int foo(void* res, void* a, void* b, void* stream).
func NewFastFuncStatus4(sym uintptr) FastFuncStatus4 {
	return func(res unsafe.Pointer, a, b, stream uintptr) int32 {
		return int32(fastCall4(sym, uintptr(res), a, b, stream))
	}
}

// FastFuncStatusIntStream calls a C function with signature: int foo(void* res, void* arr, int axis, void* stream).
type FastFuncStatusIntStream = func(unsafe.Pointer, uintptr, int32, uintptr) int32

// NewFastFuncStatusIntStream returns a zero-allocation callable for a C function
// with signature int foo(void* res, void* arr, int axis, void* stream).
func NewFastFuncStatusIntStream(sym uintptr) FastFuncStatusIntStream {
	return func(res unsafe.Pointer, arr uintptr, axis int32, stream uintptr) int32 {
		return int32(fastCall4(sym, uintptr(res), arr, uintptr(axis), stream))
	}
}

// FastFuncStatusAxes calls a C function with signature: int foo(void* res, void* arr, int* axes, size_t len, void* stream).
type FastFuncStatusAxes = func(unsafe.Pointer, uintptr, unsafe.Pointer, uint, uintptr) int32

// NewFastFuncStatusAxes returns a zero-allocation callable for a C function
// with signature int foo(void* res, void* arr, int* axes, size_t len, void* stream).
func NewFastFuncStatusAxes(sym uintptr) FastFuncStatusAxes {
	return func(res unsafe.Pointer, arr uintptr, axes unsafe.Pointer, length uint, stream uintptr) int32 {
		return int32(fastCall5(sym, uintptr(res), arr, uintptr(axes), uintptr(length), stream))
	}
}

// FastFuncStatus2IntStream calls a C function with signature: int foo(void* res, void* arr, void* indices, int axis, void* stream).
type FastFuncStatus2IntStream = func(unsafe.Pointer, uintptr, uintptr, int32, uintptr) int32

// NewFastFuncStatus2IntStream returns a zero-allocation callable for a C function
// with signature int foo(void* res, void* arr, void* indices, int axis, void* stream).
func NewFastFuncStatus2IntStream(sym uintptr) FastFuncStatus2IntStream {
	return func(res unsafe.Pointer, arr, indices uintptr, axis int32, stream uintptr) int32 {
		return int32(fastCall5(sym, uintptr(res), arr, indices, uintptr(axis), stream))
	}
}
