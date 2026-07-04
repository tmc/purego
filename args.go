// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego

import (
	"math"
	"runtime"
	"unsafe"
)

// Args builds an argument list for a C call and invokes it without allocating.
//
// Args is the safe, allocation-free counterpart to [CallN]. A caller declares it
// as a local variable, appends arguments with the typed Int, Uintptr, Float32,
// Float64, and Ptr methods, and invokes the function with Call:
//
//	var a purego.Args
//	a.Int(3)
//	a.Float64(1.5)
//	a.Ptr(unsafe.Pointer(buf)) // pinned for the duration of the call
//	r1, _ := a.Call(fn)
//
// Unlike [CallN], Args keeps every pointer added with Ptr alive and unmoving for
// the duration of the call using a [runtime.Pinner], so passing a Go pointer is
// safe with no further care from the caller. Because Args holds its argument
// registers in fixed-size arrays and its pinner inline, a stack-allocated Args
// performs no heap allocation, including when pointer arguments are pinned.
//
// The zero value is ready to use. An Args is not safe for concurrent use. After
// Call returns, the Args is reset and may be reused for another call.
//
// Args is not supported on Windows; Call panics there. Use [SyscallN].
type Args struct {
	ints   [MaxArgs]uintptr
	floats [8]uintptr
	nint   int
	nfloat int
	r8     uintptr
	pinner runtime.Pinner
	pinned bool

	// retf1 holds the bit pattern of the first floating-point result register from
	// the most recent Call, so Float64Result and Float32Result can recover a float
	// return. It survives reset so it can be read after Call returns.
	retf1 uintptr
}

// Int appends a signed integer argument.
func (a *Args) Int(v int) { a.Uintptr(uintptr(v)) }

// Uintptr appends a machine-word integer argument. Use Ptr, not Uintptr, to pass
// a Go pointer: a pointer passed as a bare uintptr is not kept alive.
func (a *Args) Uintptr(v uintptr) {
	if a.nint >= MaxArgs {
		panic("purego: too many integer arguments to Args")
	}
	a.ints[a.nint] = v
	a.nint++
}

// Float32 appends a 32-bit floating-point argument.
func (a *Args) Float32(v float32) { a.float(float32ArgBits(math.Float32bits(v))) }

// Float64 appends a 64-bit floating-point argument.
func (a *Args) Float64(v float64) { a.float(uintptr(math.Float64bits(v))) }

func (a *Args) float(bits uintptr) {
	if a.nfloat >= len(a.floats) {
		panic("purego: too many floating-point arguments to Args")
	}
	a.floats[a.nfloat] = bits
	a.nfloat++
}

// Ptr appends a pointer argument and pins the pointed-to object so it stays alive
// and unmoving for the duration of the next Call. The pointer occupies an integer
// argument slot. Passing a nil pointer is allowed and pins nothing.
func (a *Args) Ptr(p unsafe.Pointer) {
	if p != nil {
		a.pinner.Pin(p)
		a.pinned = true
	}
	a.Uintptr(uintptr(p))
}

// StructReturn sets the indirect struct-return pointer passed in x8 on arm64. It
// is ignored on other architectures. Pass a pointer to caller-owned memory large
// enough to hold the returned struct.
func (a *Args) StructReturn(p unsafe.Pointer) {
	if p != nil {
		a.pinner.Pin(p)
		a.pinned = true
	}
	a.r8 = uintptr(p)
}

// Call invokes fn with the accumulated arguments and returns the first two
// integer result registers. To recover a floating-point return value, call
// Float64Result or Float32Result after Call. It unpins any pinned pointers and
// resets the Args for reuse before returning.
func (a *Args) Call(fn uintptr) (r1, r2 uintptr) {
	var rf1 uintptr
	r1, r2, rf1, _ = CallN(fn, &a.ints, &a.floats, a.r8)
	if a.pinned {
		a.pinner.Unpin()
	}
	a.reset()
	a.retf1 = rf1
	return
}

// Float64Result returns the float64 result of the most recent Call. It is only
// meaningful when the called function returns a double.
func (a *Args) Float64Result() float64 { return math.Float64frombits(uint64(a.retf1)) }

// Float32Result returns the float32 result of the most recent Call. It is only
// meaningful when the called function returns a float.
func (a *Args) Float32Result() float32 { return math.Float32frombits(float32ResultBits(a.retf1)) }

func (a *Args) reset() {
	for i := 0; i < a.nint; i++ {
		a.ints[i] = 0
	}
	for i := 0; i < a.nfloat; i++ {
		a.floats[i] = 0
	}
	a.nint = 0
	a.nfloat = 0
	a.r8 = 0
	a.pinned = false
}
