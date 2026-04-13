// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package purego_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

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
	r1, _ := purego.CallN(sym, &ints, &floats, 0)
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
	r1, _ := purego.CallN(sym, &ints, &floats, 0)
	if got := int(r1); got != 12345 {
		t.Fatalf("atoi(\"12345\") = %d, want 12345", got)
	}
	runtime.KeepAlive(str)
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
