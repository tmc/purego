// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build darwin && arm64

package purego_test

import (
	"path/filepath"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

func openFloat16TestLibrary(t *testing.T) uintptr {
	t.Helper()
	path := filepath.Join(t.TempDir(), "float16-test.dylib")
	if err := buildSharedLib(t, "CC", path, filepath.Join("testdata", "abitest", "float16_test.c")); err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := load.CloseLibrary(lib); err != nil {
			t.Error(err)
		}
	})
	return lib
}

func TestFloat16ArgumentBits(t *testing.T) {
	lib := openFloat16TestLibrary(t)
	var bits func(purego.Float16) uint16
	purego.RegisterLibFunc(&bits, lib, "float16_bits")
	for i := 0; i <= 0xffff; i++ {
		if got := bits(purego.Float16(i)); got != uint16(i) {
			t.Fatalf("bits(%#04x) = %#04x", i, got)
		}
	}
}

func TestFloat16MixedRegisterBanks(t *testing.T) {
	lib := openFloat16TestLibrary(t)
	var mixed func(uintptr, uintptr, uintptr, uintptr, uintptr, purego.Float16) uint64
	purego.RegisterLibFunc(&mixed, lib, "float16_mixed_bank")
	x := [...]uintptr{0x11, 0x23, 0x47, 0x89, 0x107}
	const h = purego.Float16(0x7e55)
	want := uint64(x[0] ^ x[1]<<1 ^ x[2]<<2 ^ x[3]<<3 ^ x[4]<<4 ^ uintptr(h))
	if got := mixed(x[0], x[1], x[2], x[3], x[4], h); got != want {
		t.Fatalf("mixed banks = %#x, want %#x", got, want)
	}
}

func TestFloat16ClassificationUsesExactTypeIdentity(t *testing.T) {
	type namedUint16 uint16
	lib := openFloat16TestLibrary(t)
	var ordinary func(uintptr, uintptr, uintptr, uintptr, uintptr, uint16) uint64
	var named func(uintptr, uintptr, uintptr, uintptr, uintptr, namedUint16) uint64
	purego.RegisterLibFunc(&ordinary, lib, "uint16_mixed_bank")
	purego.RegisterLibFunc(&named, lib, "uint16_mixed_bank")
	x := [...]uintptr{0x11, 0x23, 0x47, 0x89, 0x107}
	const value = uint16(0x7e55)
	want := uint64(x[0] ^ x[1]<<1 ^ x[2]<<2 ^ x[3]<<3 ^ x[4]<<4 ^ uintptr(value))
	if got := ordinary(x[0], x[1], x[2], x[3], x[4], value); got != want {
		t.Fatalf("ordinary uint16 = %#x, want %#x", got, want)
	}
	if got := named(x[0], x[1], x[2], x[3], x[4], namedUint16(value)); got != want {
		t.Fatalf("named uint16 = %#x, want %#x", got, want)
	}
}

func TestFloat16IndependentRegisterOrdinals(t *testing.T) {
	lib := openFloat16TestLibrary(t)
	var interleaved func(uintptr, purego.Float16, uintptr, purego.Float16) uint64
	purego.RegisterLibFunc(&interleaved, lib, "float16_interleaved")
	const x0, x1 = uintptr(0x123), uintptr(0x456)
	const h0, h1 = purego.Float16(0x3c00), purego.Float16(0xc000)
	want := uint64(x0 ^ x1<<1 ^ uintptr(h0)<<32 ^ uintptr(h1)<<48)
	if got := interleaved(x0, h0, x1, h1); got != want {
		t.Fatalf("interleaved = %#x, want %#x", got, want)
	}
}

func TestFloat16NinthArgumentUsesPackedStack(t *testing.T) {
	lib := openFloat16TestLibrary(t)
	var nine func(purego.Float16, purego.Float16, purego.Float16, purego.Float16, purego.Float16, purego.Float16, purego.Float16, purego.Float16, purego.Float16) uint64
	purego.RegisterLibFunc(&nine, lib, "float16_nine")
	values := [...]purego.Float16{0x0000, 0x8000, 0x3c00, 0xc000, 0x0001, 0x7c00, 0xfc00, 0x7e55, 0x3555}
	var want uint64
	for _, value := range values {
		want = want*65599 + uint64(value)
	}
	if got := nine(values[0], values[1], values[2], values[3], values[4], values[5], values[6], values[7], values[8]); got != want {
		t.Fatalf("nine = %#x, want %#x", got, want)
	}
}

func TestFloat16UnsupportedRoles(t *testing.T) {
	tests := []struct {
		name string
		fn   any
	}{
		{"result", new(func() purego.Float16)},
		{"struct argument", new(func(struct{ Value purego.Float16 }))},
		{"callback argument", new(func(func(purego.Float16)))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("RegisterFunc did not reject unsupported Float16 role")
				}
			}()
			purego.RegisterFunc(tt.fn, 1)
		})
	}
}

func TestFloat16CallbackAndVariadicAreRejected(t *testing.T) {
	for _, callback := range []any{
		func(purego.Float16) {},
		func() purego.Float16 { return 0 },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("NewCallback did not reject Float16")
				}
			}()
			purego.NewCallback(callback)
		}()
	}
	var variadic func(...any)
	purego.RegisterFunc(&variadic, 1)
	defer func() {
		if recover() == nil {
			t.Fatal("variadic call did not reject Float16")
		}
	}()
	variadic(purego.Float16(0x3c00))
}

func TestFloat16VariadicAggregatesAreRejected(t *testing.T) {
	var variadic func(...any)
	purego.RegisterFunc(&variadic, 1)
	tests := []struct {
		name  string
		value any
	}{
		{"struct", struct{ Value purego.Float16 }{}},
		{"nested struct", struct {
			Inner struct{ Value purego.Float16 }
		}{}},
		{"array", [1]purego.Float16{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if got := recover(); got != "purego: Float16 in aggregate variadic arguments is not supported" {
					t.Fatalf("panic = %v", got)
				}
			}()
			variadic(test.value)
		})
	}
}
