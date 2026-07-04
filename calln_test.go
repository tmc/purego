// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

func openBenchmarkLibrary(t testing.TB) uintptr {
	t.Helper()

	libFileName := filepath.Join(t.TempDir(), "libbenchmark.so")
	if err := buildSharedLib(t, "CC", libFileName, filepath.Join("testdata", "benchmarktest", "benchmark.c")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(libFileName)
	})

	lib, err := load.OpenLibrary(libFileName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := load.CloseLibrary(lib); err != nil {
			t.Errorf("failed to close library: %v", err)
		}
	})
	return lib
}

func TestCallN_abs(t *testing.T) {
	library, err := getSystemLibrary()
	if err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(library)
	if err != nil {
		t.Fatal(err)
	}
	sym, err := purego.Dlsym(lib, "abs")
	if err != nil {
		t.Fatal(err)
	}
	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0] = uintptr(uint32(0xFFFFFFD6)) // -42 as uint32
	r1, _, _, _ := purego.CallN(sym, &ints, &floats, 0)
	if got := int32(r1); got != 42 {
		t.Fatalf("abs(-42) = %d, want 42", got)
	}
}

func TestCallN_atoi(t *testing.T) {
	library, err := getSystemLibrary()
	if err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(library)
	if err != nil {
		t.Fatal(err)
	}
	sym, err := purego.Dlsym(lib, "atoi")
	if err != nil {
		t.Fatal(err)
	}
	str := []byte("12345\x00")
	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0] = uintptr(unsafe.Pointer(&str[0]))
	r1, _, _, _ := purego.CallN(sym, &ints, &floats, 0)
	if got := int(r1); got != 12345 {
		t.Fatalf("atoi(\"12345\") = %d, want 12345", got)
	}
	runtime.KeepAlive(str)
}

func TestCallN_Order5(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("CallN not supported on Windows")
	}
	lib := openBenchmarkLibrary(t)
	sym, err := load.OpenSymbol(lib, "order5_c")
	if err != nil {
		t.Fatal(err)
	}

	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0], ints[1], ints[2], ints[3], ints[4] = 1, 2, 3, 4, 5
	r1, _, _, _ := purego.CallN(sym, &ints, &floats, 0)
	if got := int64(r1); got != 54321 {
		t.Fatalf("order5_c = %d, want 54321", got)
	}
}

func TestCallN_MixedIntFloat(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("CallN not supported on Windows")
	}
	lib := openBenchmarkLibrary(t)
	sym, err := load.OpenSymbol(lib, "mixed_int_float_c")
	if err != nil {
		t.Fatal(err)
	}

	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0], ints[1] = 3, 4
	floats[0] = uintptr(math.Float64bits(1.5))
	floats[1] = uintptr(math.Float64bits(2.5))
	r1, _, _, _ := purego.CallN(sym, &ints, &floats, 0)
	if got := int64(r1); got != 2918 {
		t.Fatalf("mixed_int_float_c = %d, want 2918", got)
	}
}

func TestCallN_ZeroAllocs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("CallN not supported on Windows")
	}
	library, err := getSystemLibrary()
	if err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(library)
	if err != nil {
		t.Fatal(err)
	}
	sym, err := purego.Dlsym(lib, "abs")
	if err != nil {
		t.Fatal(err)
	}
	// Warm up
	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0] = 42
	purego.CallN(sym, &ints, &floats, 0)

	allocs := testing.AllocsPerRun(1000, func() {
		var ints [purego.MaxArgs]uintptr
		var floats [8]uintptr
		ints[0] = 42
		purego.CallN(sym, &ints, &floats, 0)
	})
	if allocs != 0 {
		t.Fatalf("CallN allocs = %v, want 0", allocs)
	}
}

func BenchmarkCallN(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("CallN not supported on Windows")
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
		var ints [purego.MaxArgs]uintptr
		var floats [8]uintptr
		ints[0] = 42
		purego.CallN(sym, &ints, &floats, 0)
	}
}

func BenchmarkSyscallN_1arg(b *testing.B) {
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
		purego.SyscallN(sym, 42)
	}
}
