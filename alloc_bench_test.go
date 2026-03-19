// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package purego_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

// These benchmarks measure allocation overhead for the exact function
// signatures used by mlx-go's hot-path FFI calls. The key question:
// does RegisterFunc allocate per-call for interleaved float signatures?

func loadBenchLib(b *testing.B) uintptr {
	b.Helper()
	libFileName := filepath.Join(b.TempDir(), "libbenchmark.so")
	if err := buildSharedLib("CC", libFileName, filepath.Join("testdata", "benchmarktest", "benchmark.c")); err != nil {
		b.Fatalf("build: %v", err)
	}
	b.Cleanup(func() { os.Remove(libFileName) })
	h, err := load.OpenLibrary(libFileName)
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	b.Cleanup(func() { load.CloseLibrary(h) })
	return h
}

// --- Integer-only signatures (should hit registerFastFunc, 0 allocs) ---

// BenchmarkAlloc_IntOnly_5args matches sum5_c: 5 int64 args.
func BenchmarkAlloc_IntOnly_5args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, int64, int64, int64, int64) int64
	purego.RegisterLibFunc(&fn, h, "sum5_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 4, 5)
	}
}

// BenchmarkAlloc_IntOnly_10args matches QuantizedMatmul shape (all ints).
func BenchmarkAlloc_IntOnly_10args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, int64, int64, int64, int64, int64, int64, int64, int64, int64) int64
	purego.RegisterLibFunc(&fn, h, "sum10_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	}
}

// BenchmarkAlloc_UintptrOnly_10args matches QuantizedMatmul using uintptr.
func BenchmarkAlloc_UintptrOnly_10args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) int32
	purego.RegisterLibFunc(&fn, h, "qmm_shape_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	}
}

// --- Trailing float signatures (should hit registerFastFuncFloat, 0 allocs) ---

// BenchmarkAlloc_TrailingFloat_5args matches weighted_sum5f_c: 5 ints + 1 trailing float.
func BenchmarkAlloc_TrailingFloat_5args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, int64, int64, int64, int64, float32) int64
	purego.RegisterLibFunc(&fn, h, "weighted_sum5f_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 4, 5, 1.0)
	}
}

// BenchmarkAlloc_TrailingFloat_3args matches weighted_sum3f_c.
func BenchmarkAlloc_TrailingFloat_3args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, int64, int64, float32) int64
	purego.RegisterLibFunc(&fn, h, "weighted_sum3f_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 1.0)
	}
}

// --- Interleaved float signatures (hits registerFastFuncInterleaved, 0 allocs) ---

// BenchmarkAlloc_InterleavedFloat_4args matches interleaved_if_c: int, float, int, int.
func BenchmarkAlloc_InterleavedFloat_4args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, float32, int64, int64) int64
	purego.RegisterLibFunc(&fn, h, "interleaved_if_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2.0, 3, 4)
	}
}

// BenchmarkAlloc_InterleavedFloat_5args matches interleaved_2f_c: int, float, int, float, int.
func BenchmarkAlloc_InterleavedFloat_5args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(int64, float32, int64, float32, int64) int64
	purego.RegisterLibFunc(&fn, h, "interleaved_2f_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2.0, 3, 4.0, 5)
	}
}

// BenchmarkAlloc_SDPA_9args matches FastSDPAFloat: 4 uintptr, 1 float32, 4 uintptr → int32.
func BenchmarkAlloc_SDPA_9args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(uintptr, uintptr, uintptr, uintptr, float32, uintptr, uintptr, uintptr, uintptr) int32
	purego.RegisterLibFunc(&fn, h, "sdpa_shape_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 4, 1.0, 5, 6, 7, 8)
	}
}

// BenchmarkAlloc_RMSNorm_5args matches FastRMSNormFloat: 3 uintptr, 1 float32, 1 uintptr → int32.
func BenchmarkAlloc_RMSNorm_5args(b *testing.B) {
	h := loadBenchLib(b)
	var fn func(uintptr, uintptr, uintptr, float32, uintptr) int32
	purego.RegisterLibFunc(&fn, h, "rmsnorm_shape_c")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fn(1, 2, 3, 1e-5, 4)
	}
}
