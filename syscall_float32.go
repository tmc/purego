// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd || windows)

package purego

import (
	"math"
	"runtime"
)

// The Syscall<N>Float1_32 family is the single-precision counterpart of the
// [Syscall1Float1] family. Each variant calls a C function that takes N integer
// (or pointer) arguments and exactly one trailing C float (float32) argument.
//
// The integer arguments a1..aN are passed in the integer argument registers in
// order; f1 is passed as a single-precision float in the first floating-point
// argument register. On platforms whose C ABI assigns integer and
// floating-point arguments to independent register files (the AArch64 AAPCS and
// the AMD64 System V ABI), a float argument occupies the low 32 bits of the
// first floating-point register, so the float bits are written there directly;
// this is correct regardless of where the float appears in the C declaration
// relative to the integer arguments.
//
// Use this family, not the float64 [Syscall1Float1] family, when the C function
// declares its scalar parameter as float rather than double. The two are not
// interchangeable: a double bit pattern read by a float parameter (or the
// reverse) yields a different value.
//
// The same restrictions as the float64 family apply: only the (N integers,
// exactly one trailing float) shape is provided; multiple floats, or arguments
// that spill onto the stack, are not expressible and are intentionally omitted.
// On Windows the underlying helper cannot place a value in a floating-point
// register, so these variants fall back to passing the raw bits of f1 as an
// additional integer argument; they are correct only on platforms that use the
// purego funnel (see [SyscallN]).

// Syscall1Float1_32 calls fn with one integer argument and one trailing float32.
// See the Syscall<N>Float1_32 family documentation for details and caveats.
//
//go:uintptrescapes
func Syscall1Float1_32(fn, a1 uintptr, f1 float32) (r1, r2, err uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if runtime.GOOS == "windows" {
		return syscall_syscallN(fn, a1, uintptr(math.Float32bits(f1)))
	}
	var tmp [maxArgs]uintptr
	tmp[0] = a1
	var floats [maxArgs]uintptr
	floats[0] = uintptr(math.Float32bits(f1))
	s := syscall_SyscallN(fn, tmp[:], floats[:], 0)
	defer thePool.Put(s)
	return s.a1, s.a2, s.a3
}

// Syscall2Float1_32 calls fn with two integer arguments and one trailing float32.
// See the Syscall<N>Float1_32 family documentation for details and caveats.
//
//go:uintptrescapes
func Syscall2Float1_32(fn, a1, a2 uintptr, f1 float32) (r1, r2, err uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if runtime.GOOS == "windows" {
		return syscall_syscallN(fn, a1, a2, uintptr(math.Float32bits(f1)))
	}
	var tmp [maxArgs]uintptr
	tmp[0], tmp[1] = a1, a2
	var floats [maxArgs]uintptr
	floats[0] = uintptr(math.Float32bits(f1))
	s := syscall_SyscallN(fn, tmp[:], floats[:], 0)
	defer thePool.Put(s)
	return s.a1, s.a2, s.a3
}

// Syscall3Float1_32 calls fn with three integer arguments and one trailing float32.
// See the Syscall<N>Float1_32 family documentation for details and caveats.
//
//go:uintptrescapes
func Syscall3Float1_32(fn, a1, a2, a3 uintptr, f1 float32) (r1, r2, err uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if runtime.GOOS == "windows" {
		return syscall_syscallN(fn, a1, a2, a3, uintptr(math.Float32bits(f1)))
	}
	var tmp [maxArgs]uintptr
	tmp[0], tmp[1], tmp[2] = a1, a2, a3
	var floats [maxArgs]uintptr
	floats[0] = uintptr(math.Float32bits(f1))
	s := syscall_SyscallN(fn, tmp[:], floats[:], 0)
	defer thePool.Put(s)
	return s.a1, s.a2, s.a3
}

// Syscall4Float1_32 calls fn with four integer arguments and one trailing float32.
// See the Syscall<N>Float1_32 family documentation for details and caveats.
//
//go:uintptrescapes
func Syscall4Float1_32(fn, a1, a2, a3, a4 uintptr, f1 float32) (r1, r2, err uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if runtime.GOOS == "windows" {
		return syscall_syscallN(fn, a1, a2, a3, a4, uintptr(math.Float32bits(f1)))
	}
	var tmp [maxArgs]uintptr
	tmp[0], tmp[1], tmp[2], tmp[3] = a1, a2, a3, a4
	var floats [maxArgs]uintptr
	floats[0] = uintptr(math.Float32bits(f1))
	s := syscall_SyscallN(fn, tmp[:], floats[:], 0)
	defer thePool.Put(s)
	return s.a1, s.a2, s.a3
}

// Syscall8Float1_32 calls fn with eight integer arguments and one trailing float32.
// See the Syscall<N>Float1_32 family documentation for details and caveats.
//
//go:uintptrescapes
func Syscall8Float1_32(fn, a1, a2, a3, a4, a5, a6, a7, a8 uintptr, f1 float32) (r1, r2, err uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if runtime.GOOS == "windows" {
		return syscall_syscallN(fn, a1, a2, a3, a4, a5, a6, a7, a8, uintptr(math.Float32bits(f1)))
	}
	var tmp [maxArgs]uintptr
	tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5], tmp[6], tmp[7] = a1, a2, a3, a4, a5, a6, a7, a8
	var floats [maxArgs]uintptr
	floats[0] = uintptr(math.Float32bits(f1))
	s := syscall_SyscallN(fn, tmp[:], floats[:], 0)
	defer thePool.Put(s)
	return s.a1, s.a2, s.a3
}
