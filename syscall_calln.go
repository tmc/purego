// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego

import (
	"runtime"
	"unsafe"
)

// MaxArgs is the maximum number of integer arguments accepted by [CallN].
const MaxArgs = maxArgs

// CallN calls the C function fn with integer arguments taken from ints and
// floating-point arguments taken from floats.
//
// CallN is an unsafe, low-level primitive. It does NOT keep Go pointers alive for
// the duration of the call: passing the address of a Go object as an argument
// without keeping that object independently reachable is a use-after-free (see
// "Pointer arguments and lifetime" below). Most callers should use [Args], the
// safe builder over CallN, which pins pointer arguments automatically. Use CallN
// directly only for hot paths that pass scalars and already-pinned buffers.
//
// CallN returns the first two integer result registers as r1 and r2, and the raw
// bit patterns of the first two floating-point result registers as rf1 and rf2. A
// function that returns a double leaves its result in the full 64 bits of rf1;
// recover it with [math.Float64frombits]. A function that returns a float leaves
// its result in one half of rf1 — the low 32 bits on the little-endian
// architectures and the high 32 bits on big-endian s390x; use [Args.Float32Result],
// which recovers it on either. A function that returns an integer or pointer leaves
// it in r1, and rf1 and rf2 are unspecified. CallN does not report errno; a caller
// that needs the C error code must use [SyscallN].
//
// Unlike [SyscallN], CallN forms no variadic slice at the call site: the caller
// supplies fixed-size arrays by pointer, so a correctly written call performs no
// heap allocation. CallN also passes integer and floating-point arguments in
// separate register files, so it calls functions with mixed int and float
// parameters correctly, which [SyscallN] does not. It also recovers a
// floating-point return value, which [SyscallN] cannot.
//
// Integer arguments occupy ints[0], ints[1], ... in order; unused trailing slots
// must be zero. Float arguments occupy floats[0..] and must be pre-converted to
// their bit patterns with [math.Float32bits] or [math.Float64bits]. A double
// occupies the full 64 bits of its slot; a float occupies the low 32 bits on the
// little-endian architectures and the high 32 bits on big-endian s390x. There are
// eight floating-point argument registers on amd64, arm64, loong64, and riscv64,
// and four on s390x; a floating-point argument beyond that count cannot be passed,
// because floating-point arguments never spill to the stack. structRet is
// the indirect result-location pointer used by the arm64 AAPCS calling convention
// (passed in x8) for a function that returns a struct too large for registers; it
// is ignored on other architectures. Pass 0 when the function does not return a
// struct by memory.
//
// ints and floats must be non-nil; CallN panics otherwise. CallN is not supported
// on Windows and panics there; use [SyscallN].
//
// # Structs
//
// CallN classifies each argument by which array it is placed in, not by
// inspecting its type, so it does not pass structs by value. A struct that the C
// ABI passes in registers (a small struct, or a homogeneous floating-point
// aggregate) cannot be expressed through the ints and floats arrays, because the
// register file its fields land in depends on the field types. Pass a large
// struct by pointer in an integer slot, receive a large struct return through
// structRet, and use [RegisterFunc] for functions that take or return small
// structs by value; RegisterFunc walks the argument types with reflection to
// place struct fields in the correct registers.
//
// # Pointer arguments and lifetime
//
// CallN does NOT keep Go pointers alive for the duration of the call. This is a
// deliberate consequence of its allocation-free design: because arguments arrive
// through a *[MaxArgs]uintptr rather than as variadic uintptr parameters, CallN
// cannot carry the //go:uintptrescapes pragma that [SyscallN] relies on to pin
// pointer arguments. Passing the address of a Go object as an integer argument
// without keeping that object independently reachable is a use-after-free: the
// garbage collector may free or move it while the C function still holds the
// address.
//
// Callers passing Go pointers must ensure the pointed-to memory stays alive and
// unmoving across the call, for example with [runtime.KeepAlive] on a variable
// that holds the object, or by using memory that is already pinned or C-owned.
// This is stricter than [SyscallN], whose //go:uintptrescapes pragma keeps
// pointer arguments alive automatically; the rules in [unsafe.Pointer],
// especially point 4, apply to CallN without that safety net. For the common case
// of passing a Go pointer safely, prefer [SyscallN] with the unsafe.Pointer→uintptr
// conversion written directly in the call expression, or use [runtime.Pinner] or
// [Args].
func CallN(fn uintptr, ints *[MaxArgs]uintptr, floats *[8]uintptr, structRet uintptr) (r1, r2, rf1, rf2 uintptr) {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	if ints == nil {
		panic("purego: CallN called with nil integer arguments")
	}
	if floats == nil {
		panic("purego: CallN called with nil floating-point arguments")
	}
	if runtime.GOOS == "windows" {
		panic("purego: CallN is not supported on Windows; use SyscallN")
	}
	s := thePool.Get().(*syscallArgs)
	*s = syscallArgs{
		fn: fn,
		a1: ints[0], a2: ints[1], a3: ints[2], a4: ints[3],
		a5: ints[4], a6: ints[5], a7: ints[6], a8: ints[7],
		a9: ints[8], a10: ints[9], a11: ints[10], a12: ints[11],
		a13: ints[12], a14: ints[13], a15: ints[14], a16: ints[15],
		a17: ints[16], a18: ints[17], a19: ints[18], a20: ints[19],
		a21: ints[20], a22: ints[21], a23: ints[22], a24: ints[23],
		a25: ints[24], a26: ints[25], a27: ints[26], a28: ints[27],
		a29: ints[28], a30: ints[29], a31: ints[30], a32: ints[31],
		f1: floats[0], f2: floats[1], f3: floats[2], f4: floats[3],
		f5: floats[4], f6: floats[5], f7: floats[6], f8: floats[7],
		arm64_r8: structRet,
	}
	runtime_cgocall(syscallXABI0, unsafe.Pointer(s))
	// The assembly funnel writes the integer result registers back into a1/a2 and
	// the floating-point result registers into f1/f2 after the call returns. Read
	// them out before returning s to the pool.
	r1, r2, rf1, rf2 = s.a1, s.a2, s.f1, s.f2
	thePool.Put(s)
	return
}
