// Validates Args floating-point return recovery.

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"testing"

	"github.com/ebitengine/purego"
)

func TestArgs_FloatReturn(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := purego.Dlsym(lib, "ret_double_c")
	if err != nil {
		t.Fatal(err)
	}
	var a purego.Args
	a.Float64(2.0)
	a.Float64(3.0)
	a.Call(fn)
	if got := a.Float64Result(); got != 203.0 {
		t.Errorf("Args.Float64Result = %v, want 203", got)
	}
}

func TestArgs_Float32Return(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := purego.Dlsym(lib, "ret_float_c")
	if err != nil {
		t.Fatal(err)
	}
	var a purego.Args
	a.Float32(1.5)
	a.Float32(2.0)
	a.Call(fn)
	if got := a.Float32Result(); got != 152.0 { // 1.5*100 + 2
		t.Errorf("Args.Float32Result = %v, want 152", got)
	}
}

func TestArgs_FloatReturn_ZeroAllocs(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := purego.Dlsym(lib, "ret_double_c")
	if err != nil {
		t.Fatal(err)
	}
	allocs := testing.AllocsPerRun(1000, func() {
		var a purego.Args
		a.Float64(2.0)
		a.Float64(3.0)
		a.Call(fn)
		_ = a.Float64Result()
	})
	if allocs != 0 {
		t.Errorf("float-return path allocates %v allocs/op, want 0", allocs)
	}
}
