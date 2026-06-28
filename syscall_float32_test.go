// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

// TestSyscallFloat1_32 verifies that the Syscall<N>Float32x1 family delivers
// the trailing argument as a single-precision C float in the first
// floating-point register. Each C helper returns sum(ints) + (long)(f1*1000);
// the float contributes 1500 only when the single-precision bit pattern landed
// in the FP register.
func TestSyscallFloat1_32(t *testing.T) {
	lib := openABITestLib(t)

	sym := func(name string) uintptr {
		s, err := load.OpenSymbol(lib, name)
		if err != nil {
			t.Fatalf("OpenSymbol(%s) failed: %v", name, err)
		}
		return s
	}

	const f1 float32 = 1.5
	const floatPart = 1500 // (long)(1.5*1000)

	fn1 := sym("calln_f32_1")
	fn2 := sym("calln_f32_2")
	fn3 := sym("calln_f32_3")
	fn4 := sym("calln_f32_4")
	fn8 := sym("calln_f32_8")

	{
		got, _, _ := purego.Syscall1Float32x1(fn1, 10, f1)
		if want := uintptr(10 + floatPart); got != want {
			t.Errorf("Syscall1Float32x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall2Float32x1(fn2, 10, 20, f1)
		if want := uintptr(30 + floatPart); got != want {
			t.Errorf("Syscall2Float32x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall3Float32x1(fn3, 10, 20, 30, f1)
		if want := uintptr(60 + floatPart); got != want {
			t.Errorf("Syscall3Float32x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall4Float32x1(fn4, 10, 20, 30, 40, f1)
		if want := uintptr(100 + floatPart); got != want {
			t.Errorf("Syscall4Float32x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall8Float32x1(fn8, 1, 2, 3, 4, 5, 6, 7, 8, f1)
		if want := uintptr(36 + floatPart); got != want {
			t.Errorf("Syscall8Float32x1: got %d, want %d", got, want)
		}
	}
}

// TestSyscallFloat1_32_NotInterchangeableWithFloat64 documents why a separate
// single-precision family is required: passing a double bit pattern (the
// float64 family) to a C function whose parameter is float yields a wrong
// value, and vice versa. This is the exact failure mode that corrupts MLX's
// fast ops (mlx_fast_rms_norm/rope/scaled_dot_product_attention all take a C
// float), which the float64 Syscall<N>Float64x1 family silently mis-delivers.
func TestSyscallFloat1_32_NotInterchangeableWithFloat64(t *testing.T) {
	lib := openABITestLib(t)

	sym := func(name string) uintptr {
		s, err := load.OpenSymbol(lib, name)
		if err != nil {
			t.Fatalf("OpenSymbol(%s) failed: %v", name, err)
		}
		return s
	}

	// calln_f32_4 declares its trailing param as float. The float32 family
	// delivers it correctly; the float64 family delivers a double bit pattern,
	// which the float parameter reads as a different (wrong) value.
	fnFloat := sym("calln_f32_4")

	const eps float32 = 0.01
	wantCorrect := uintptr(100 + (int)(eps*1000)) // 100 + 10 = 110

	if got, _, _ := purego.Syscall4Float32x1(fnFloat, 10, 20, 30, 40, eps); got != wantCorrect {
		t.Errorf("float32 family on a float param: got %d, want %d", got, wantCorrect)
	}

	// The float64 family on the same float param must NOT produce the correct
	// value — that mismatch is the whole reason this family exists. If a future
	// ABI change ever made them coincide, this guard would flag it.
	if got, _, _ := purego.Syscall4Float64x1(fnFloat, 10, 20, 30, 40, float64(eps)); got == wantCorrect {
		t.Errorf("float64 family unexpectedly matched on a float param (got %d); "+
			"the single/double encodings should differ", got)
	}
}

func BenchmarkSyscall4Float32x1(b *testing.B) {
	lib := openABITestLib(b)
	fn, err := load.OpenSymbol(lib, "calln_f32_4")
	if err != nil {
		b.Fatalf("OpenSymbol(calln_f32_4) failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r uintptr
	for i := 0; i < b.N; i++ {
		r, _, _ = purego.Syscall4Float32x1(fn, 10, 20, 30, 40, 1.5)
	}
	b.StopTimer()
	if want := uintptr(100 + 1500); r != want {
		b.Fatalf("Syscall4Float32x1: got %d, want %d", r, want)
	}
}

func BenchmarkSyscall8Float32x1(b *testing.B) {
	lib := openABITestLib(b)
	fn, err := load.OpenSymbol(lib, "calln_f32_8")
	if err != nil {
		b.Fatalf("OpenSymbol(calln_f32_8) failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r uintptr
	for i := 0; i < b.N; i++ {
		r, _, _ = purego.Syscall8Float32x1(fn, 1, 2, 3, 4, 5, 6, 7, 8, 1.5)
	}
	b.StopTimer()
	if want := uintptr(36 + 1500); r != want {
		b.Fatalf("Syscall8Float32x1: got %d, want %d", r, want)
	}
}
