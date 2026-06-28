// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"path/filepath"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

// openABITestLib builds and opens the abitest shared library used by the
// Syscall<N>Float64x1 tests and benchmarks.
func openABITestLib(tb testing.TB) uintptr {
	tb.Helper()
	libFileName := filepath.Join(tb.TempDir(), "abitest.so")
	if err := buildSharedLib("CC", libFileName, filepath.Join("testdata", "abitest", "abi_test.c")); err != nil {
		tb.Fatal(err)
	}
	lib, err := load.OpenLibrary(libFileName)
	if err != nil {
		tb.Fatalf("Dlopen(%q) failed: %v", libFileName, err)
	}
	tb.Cleanup(func() {
		if err := load.CloseLibrary(lib); err != nil {
			tb.Errorf("failed to close library: %s", err)
		}
	})
	return lib
}

// TestSyscallFloat1 verifies that the Syscall<N>Float64x1 family places the
// trailing float64 in the floating-point argument register rather than an
// integer register. Each C helper returns sum(ints) + (long)(f1*1000); with
// f1 == 1.5 the float contributes 1500 only when it landed in the FP register.
func TestSyscallFloat1(t *testing.T) {
	lib := openABITestLib(t)

	sym := func(name string) uintptr {
		s, err := load.OpenSymbol(lib, name)
		if err != nil {
			t.Fatalf("OpenSymbol(%s) failed: %v", name, err)
		}
		return s
	}

	const f1 = 1.5
	const floatPart = 1500 // (long)(1.5*1000)

	fn1 := sym("calln_float1")
	fn2 := sym("calln_float2")
	fn3 := sym("calln_float3")
	fn4 := sym("calln_float4")
	fn8 := sym("calln_float8")

	{
		got, _, _ := purego.Syscall1Float64x1(fn1, 10, f1)
		if want := uintptr(10 + floatPart); got != want {
			t.Errorf("Syscall1Float64x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall2Float64x1(fn2, 10, 20, f1)
		if want := uintptr(30 + floatPart); got != want {
			t.Errorf("Syscall2Float64x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall3Float64x1(fn3, 10, 20, 30, f1)
		if want := uintptr(60 + floatPart); got != want {
			t.Errorf("Syscall3Float64x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall4Float64x1(fn4, 10, 20, 30, 40, f1)
		if want := uintptr(100 + floatPart); got != want {
			t.Errorf("Syscall4Float64x1: got %d, want %d", got, want)
		}
	}
	{
		got, _, _ := purego.Syscall8Float64x1(fn8, 1, 2, 3, 4, 5, 6, 7, 8, f1)
		if want := uintptr(36 + floatPart); got != want {
			t.Errorf("Syscall8Float64x1: got %d, want %d", got, want)
		}
	}
}

// TestSyscallFloat1_DistinctValues exercises the variants with a non-trivial
// float and asymmetric integer arguments to catch register misordering that a
// symmetric input (all equal) would hide.
func TestSyscallFloat1_DistinctValues(t *testing.T) {
	lib := openABITestLib(t)

	fn8, err := load.OpenSymbol(lib, "calln_float8")
	if err != nil {
		t.Fatalf("OpenSymbol(calln_float8) failed: %v", err)
	}

	// f1 == 3.141 -> (long)(3.141*1000) == 3141.
	got, _, _ := purego.Syscall8Float64x1(fn8, 100, 200, 300, 400, 500, 600, 700, 800, 3.141)
	if want := uintptr(3600 + 3141); got != want {
		t.Errorf("Syscall8Float64x1: got %d, want %d", got, want)
	}
}

func BenchmarkSyscall4Float64x1(b *testing.B) {
	lib := openABITestLib(b)
	fn, err := load.OpenSymbol(lib, "calln_float4")
	if err != nil {
		b.Fatalf("OpenSymbol(calln_float4) failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r uintptr
	for i := 0; i < b.N; i++ {
		r, _, _ = purego.Syscall4Float64x1(fn, 10, 20, 30, 40, 1.5)
	}
	b.StopTimer()
	if want := uintptr(100 + 1500); r != want {
		b.Fatalf("Syscall4Float64x1: got %d, want %d", r, want)
	}
}

func BenchmarkSyscall8Float64x1(b *testing.B) {
	lib := openABITestLib(b)
	fn, err := load.OpenSymbol(lib, "calln_float8")
	if err != nil {
		b.Fatalf("OpenSymbol(calln_float8) failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r uintptr
	for i := 0; i < b.N; i++ {
		r, _, _ = purego.Syscall8Float64x1(fn, 1, 2, 3, 4, 5, 6, 7, 8, 1.5)
	}
	b.StopTimer()
	if want := uintptr(36 + 1500); r != want {
		b.Fatalf("Syscall8Float64x1: got %d, want %d", r, want)
	}
}
