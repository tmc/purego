// Validates CallN floating-point return recovery.

//go:build !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego_test

import (
	"math"
	"testing"

	"github.com/ebitengine/purego"
)

func TestCallN_FloatReturn(t *testing.T) {
	lib := openBenchmarkLibrary(t)
	fn, err := purego.Dlsym(lib, "ret_double_c")
	if err != nil {
		t.Fatal(err)
	}
	var ints [purego.MaxArgs]uintptr
	var floats [8]uintptr
	floats[0] = uintptr(math.Float64bits(2.0))
	floats[1] = uintptr(math.Float64bits(3.0))
	_, _, rf1, _ := purego.CallN(fn, &ints, &floats, 0)
	got := math.Float64frombits(uint64(rf1))
	if got != 203.0 { // 2*100 + 3
		t.Errorf("CallN float return = %v, want 203", got)
	}
}
