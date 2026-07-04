// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

// TestArgs_Order5 checks integer argument ordering through the Args builder.
func TestArgs_Order5(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := load.OpenSymbol(lib, "order5_c")
	if err != nil {
		t.Fatal(err)
	}
	var a purego.Args
	a.Int(1)
	a.Int(2)
	a.Int(3)
	a.Int(4)
	a.Int(5)
	r1, _ := a.Call(fn)
	if int64(r1) != 54321 {
		t.Fatalf("order5_c via Args = %d, want 54321", int64(r1))
	}
}

// TestArgs_MixedIntFloat checks that interleaved int/float args land in the right
// register files — the case SyscallN cannot express.
func TestArgs_MixedIntFloat(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := load.OpenSymbol(lib, "mixed_int_float_c")
	if err != nil {
		t.Fatal(err)
	}
	var a purego.Args
	a.Int(3)       // a1
	a.Float64(1.5) // f1
	a.Int(4)       // a2
	a.Float64(2.5) // f2
	r1, _ := a.Call(fn)
	if int64(r1) != 2918 {
		t.Fatalf("mixed_int_float_c via Args = %d, want 2918", int64(r1))
	}
}

// TestArgs_Reuse verifies an Args is reset and reusable after Call.
func TestArgs_Reuse(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := load.OpenSymbol(lib, "order5_c")
	if err != nil {
		t.Fatal(err)
	}
	var a purego.Args
	for i := 0; i < 3; i++ {
		a.Int(1)
		a.Int(2)
		a.Int(3)
		a.Int(4)
		a.Int(5)
		r1, _ := a.Call(fn)
		if int64(r1) != 54321 {
			t.Fatalf("reuse iter %d: got %d, want 54321", i, int64(r1))
		}
	}
}

// TestArgs_Ptr passes a Go pointer via Args.Ptr, which pins it for the duration
// of the call. The C helper reads a uint64 back through the pointer, so a correct
// call returns the sentinel. This documents and regression-tests that Args.Ptr
// makes pointer arguments safe without any caller-side pinning.
func TestArgs_Ptr(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	deref, err := load.OpenSymbol(lib, "deref_u64")
	if err != nil {
		t.Fatal(err)
	}
	const magic = 0xDEADBEEFCAFEF00D
	v := uint64(magic)

	var a purego.Args
	a.Ptr(unsafe.Pointer(&v))
	got, _ := a.Call(deref)
	if uint64(got) != magic {
		t.Fatalf("Args.Ptr(&v) -> deref_u64 = %#x, want %#x", uint64(got), uint64(magic))
	}
}

// TestArgs_ZeroAllocs is the headline property: the safe builder allocates
// nothing, including when a pointer argument is pinned.
func TestArgs_ZeroAllocs(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	deref, err := load.OpenSymbol(lib, "deref_u64")
	if err != nil {
		t.Fatal(err)
	}
	v := uint64(0xDEADBEEFCAFEF00D)

	// Warm up.
	{
		var a purego.Args
		a.Ptr(unsafe.Pointer(&v))
		a.Call(deref)
	}

	n := testing.AllocsPerRun(1000, func() {
		var a purego.Args
		a.Ptr(unsafe.Pointer(&v)) // pins v
		a.Call(deref)
	})
	if n != 0 {
		t.Fatalf("Args.Call with a pinned pointer = %v allocs/op, want 0", n)
	}
}

func BenchmarkArgs_1arg(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Args not supported on Windows")
	}
	library, err := getSystemLibrary()
	if err != nil {
		b.Fatal(err)
	}
	lib, err := load.OpenLibrary(library)
	if err != nil {
		b.Fatal(err)
	}
	sym, err := purego.Dlsym(lib, "abs")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		var a purego.Args
		a.Int(42)
		a.Call(sym)
	}
}

func BenchmarkArgs_ptr(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Args not supported on Windows")
	}
	library, err := getSystemLibrary()
	if err != nil {
		b.Fatal(err)
	}
	lib, err := load.OpenLibrary(library)
	if err != nil {
		b.Fatal(err)
	}
	sym, err := purego.Dlsym(lib, "abs")
	if err != nil {
		b.Fatal(err)
	}
	// A stable heap object whose address is passed (and pinned) each call.
	buf := new([16]byte)
	b.ResetTimer()
	for b.Loop() {
		var a purego.Args
		a.Int(42)
		a.Ptr(unsafe.Pointer(buf)) // exercises the pin+unpin path
		a.Call(sym)
	}
	runtime.KeepAlive(buf)
}
