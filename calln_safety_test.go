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

// TestCallN_PointerArg exercises the documented safe pattern for passing a Go
// pointer through CallN: pin it with a runtime.Pinner for the duration of the
// call. The C helper reads a uint64 back through the pointer, so a correct call
// returns the sentinel value.
//
// This is deterministic: it does not depend on winning a race against the
// garbage collector. It documents and regression-tests the pin-the-pointer
// contract described on CallN.
func TestCallN_PointerArg(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	deref, err := load.OpenSymbol(lib, "deref_u64")
	if err != nil {
		t.Fatalf("OpenSymbol(deref_u64): %v", err)
	}

	const magic = 0xDEADBEEFCAFEF00D
	v := uint64(magic)

	var pinner runtime.Pinner
	pinner.Pin(&v)
	defer pinner.Unpin()

	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	ints[0] = uintptr(unsafe.Pointer(&v))

	got, _, _, _ := purego.CallN(deref, &ints, &floats, 0)
	if uint64(got) != magic {
		t.Fatalf("CallN(deref_u64, &v) = %#x, want %#x", uint64(got), uint64(magic))
	}
}
