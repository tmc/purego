// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd

package purego_test

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

// TestConcurrentRegisterFunc exercises the caching and pooling logic
// under concurrent access. Multiple goroutines register and call
// functions with identical signatures simultaneously.
func TestConcurrentRegisterFunc(t *testing.T) {
	libName, err := getSystemLibrary()
	if err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(libName)
	if err != nil {
		t.Fatal(err)
	}

	var abs func(int32) int32
	purego.RegisterLibFunc(&abs, lib, "abs")

	const goroutines = 16
	const calls = 1000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			// Each goroutine registers its own function to exercise
			// concurrent sync.Map + sync.Pool access.
			var localAbs func(int32) int32
			purego.RegisterLibFunc(&localAbs, lib, "abs")
			for i := 0; i < calls; i++ {
				got := localAbs(-42)
				if got != 42 {
					t.Errorf("abs(-42) = %d, want 42", got)
					return
				}
			}
		}()
	}
	wg.Wait()
}

// TestCallbackPathUnaffected verifies that callback-based round trips
// (Go -> C -> Go) still work correctly with the optimization changes.
func TestCallbackPathUnaffected(t *testing.T) {
	libName, err := getSystemLibrary()
	if err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(libName)
	if err != nil {
		t.Fatal(err)
	}

	data := []int{88, 56, 100, 2, 25}
	sorted := []int{2, 25, 56, 88, 100}
	compare := func(_ purego.CDecl, a, b *int) int {
		return *a - *b
	}
	var qsort func(data []int, nitms uintptr, size uintptr, compar func(_ purego.CDecl, a, b *int) int)
	purego.RegisterLibFunc(&qsort, lib, "qsort")
	qsort(data, uintptr(len(data)), unsafe.Sizeof(int(0)), compare)
	for i := range data {
		if data[i] != sorted[i] {
			t.Errorf("qsort: got %d wanted %d at index %d", data[i], sorted[i], i)
		}
	}
}
