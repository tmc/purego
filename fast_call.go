// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || windows

package purego

import "unsafe"

// Exported fast-call API
//
// FastCall0 through FastCall10 call a C function with the given number of
// integer/pointer arguments without heap allocations. All arguments and
// the return value are passed as uintptr. The caller is responsible for
// converting Go types to uintptr (e.g. unsafe.Pointer, math.Float32bits).
//
// FastCallNF1 and FastCallNF2 variants additionally populate the f1 (and f2)
// float registers. Float arguments must be pre-converted to their bit
// representation (math.Float32bits or math.Float64bits) and passed as uintptr.

// FastCall0 calls a C function with 0 arguments without heap allocations.
func FastCall0(cfn uintptr) uintptr { return fastCall0(cfn) }

// FastCall1 calls a C function with 1 argument without heap allocations.
func FastCall1(cfn, a1 uintptr) uintptr { return fastCall1(cfn, a1) }

// FastCall2 calls a C function with 2 arguments without heap allocations.
func FastCall2(cfn, a1, a2 uintptr) uintptr { return fastCall2(cfn, a1, a2) }

// FastCall3 calls a C function with 3 arguments without heap allocations.
func FastCall3(cfn, a1, a2, a3 uintptr) uintptr { return fastCall3(cfn, a1, a2, a3) }

// FastCall4 calls a C function with 4 arguments without heap allocations.
func FastCall4(cfn, a1, a2, a3, a4 uintptr) uintptr { return fastCall4(cfn, a1, a2, a3, a4) }

// FastCall5 calls a C function with 5 arguments without heap allocations.
func FastCall5(cfn, a1, a2, a3, a4, a5 uintptr) uintptr {
	return fastCall5(cfn, a1, a2, a3, a4, a5)
}

// FastCall6 calls a C function with 6 arguments without heap allocations.
func FastCall6(cfn, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
	return fastCall6(cfn, a1, a2, a3, a4, a5, a6)
}

// FastCall7 calls a C function with 7 arguments without heap allocations.
func FastCall7(cfn, a1, a2, a3, a4, a5, a6, a7 uintptr) uintptr {
	return fastCall7(cfn, a1, a2, a3, a4, a5, a6, a7)
}

// FastCall8 calls a C function with 8 arguments without heap allocations.
func FastCall8(cfn, a1, a2, a3, a4, a5, a6, a7, a8 uintptr) uintptr {
	return fastCall8(cfn, a1, a2, a3, a4, a5, a6, a7, a8)
}

// FastCall9 calls a C function with 9 arguments without heap allocations.
func FastCall9(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr {
	return fastCall9(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9)
}

// FastCall10 calls a C function with 10 arguments without heap allocations.
func FastCall10(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) uintptr {
	return fastCall10(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
}

// FastCall3F1 calls a C function with 3 integer args and 1 float arg.
// The float arg f1 must be pre-converted via math.Float32bits or math.Float64bits.
func FastCall3F1(cfn, a1, a2, a3, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, 0, 0, 0, 0, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall4F1 calls a C function with 4 integer args and 1 float arg.
func FastCall4F1(cfn, a1, a2, a3, a4, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, 0, 0, 0, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall5F1 calls a C function with 5 integer args and 1 float arg.
func FastCall5F1(cfn, a1, a2, a3, a4, a5, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, 0, 0, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall5F2 calls a C function with 5 integer args and 2 float args.
func FastCall5F2(cfn, a1, a2, a3, a4, a5, f1, f2 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, 0, 0, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, f2, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall6F1 calls a C function with 6 integer args and 1 float arg.
func FastCall6F1(cfn, a1, a2, a3, a4, a5, a6, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, 0, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall7F1 calls a C function with 7 integer args and 1 float arg.
func FastCall7F1(cfn, a1, a2, a3, a4, a5, a6, a7, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, 0
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall8F1 calls a C function with 8 integer args and 1 float arg.
func FastCall8F1(cfn, a1, a2, a3, a4, a5, a6, a7, a8, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall8F2 calls a C function with 8 integer args and 2 float args.
func FastCall8F2(cfn, a1, a2, a3, a4, a5, a6, a7, a8, f1, f2 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9, s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, f2, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall9F1 calls a C function with 9 integer args and 1 float arg.
func FastCall9F1(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9 = a9
	s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall10F1 calls a C function with 10 integer args and 1 float arg.
func FastCall10F1(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9, s.a10 = a9, a10
	s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}
