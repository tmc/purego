// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || windows

package purego

import (
	"reflect"
	"unsafe"
)

// registerFastFunc sets fn to a zero-allocation closure that calls cfn
// without using reflect.MakeFunc. This works because all integer-ABI types
// (int*, uint*, ptr, bool, uintptr) have identical calling convention layout
// on 64-bit platforms, so a func(uintptr, ...) uintptr closure can be
// force-assigned to the user's function pointer.
func registerFastFunc(fn reflect.Value, cfn uintptr, numArgs int, hasReturn bool) {
	switch numArgs {
	case 0:
		if hasReturn {
			impl := func() uintptr { return fastCall0(cfn) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func() { fastCall0(cfn) }
			forceFuncPtr(fn, &impl)
		}
	case 1:
		if hasReturn {
			impl := func(a1 uintptr) uintptr { return fastCall1(cfn, a1) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1 uintptr) { fastCall1(cfn, a1) }
			forceFuncPtr(fn, &impl)
		}
	case 2:
		if hasReturn {
			impl := func(a1, a2 uintptr) uintptr { return fastCall2(cfn, a1, a2) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2 uintptr) { fastCall2(cfn, a1, a2) }
			forceFuncPtr(fn, &impl)
		}
	case 3:
		if hasReturn {
			impl := func(a1, a2, a3 uintptr) uintptr { return fastCall3(cfn, a1, a2, a3) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3 uintptr) { fastCall3(cfn, a1, a2, a3) }
			forceFuncPtr(fn, &impl)
		}
	case 4:
		if hasReturn {
			impl := func(a1, a2, a3, a4 uintptr) uintptr { return fastCall4(cfn, a1, a2, a3, a4) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4 uintptr) { fastCall4(cfn, a1, a2, a3, a4) }
			forceFuncPtr(fn, &impl)
		}
	case 5:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5 uintptr) uintptr { return fastCall5(cfn, a1, a2, a3, a4, a5) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5 uintptr) { fastCall5(cfn, a1, a2, a3, a4, a5) }
			forceFuncPtr(fn, &impl)
		}
	case 6:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5, a6 uintptr) uintptr { return fastCall6(cfn, a1, a2, a3, a4, a5, a6) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5, a6 uintptr) { fastCall6(cfn, a1, a2, a3, a4, a5, a6) }
			forceFuncPtr(fn, &impl)
		}
	case 7:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5, a6, a7 uintptr) uintptr { return fastCall7(cfn, a1, a2, a3, a4, a5, a6, a7) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5, a6, a7 uintptr) { fastCall7(cfn, a1, a2, a3, a4, a5, a6, a7) }
			forceFuncPtr(fn, &impl)
		}
	case 8:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8 uintptr) uintptr { return fastCall8(cfn, a1, a2, a3, a4, a5, a6, a7, a8) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8 uintptr) { fastCall8(cfn, a1, a2, a3, a4, a5, a6, a7, a8) }
			forceFuncPtr(fn, &impl)
		}
	case 9:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr { return FastCall9(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) { FastCall9(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9) }
			forceFuncPtr(fn, &impl)
		}
	case 10:
		if hasReturn {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) uintptr { return FastCall10(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10) }
			forceFuncPtr(fn, &impl)
		} else {
			impl := func(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) { FastCall10(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10) }
			forceFuncPtr(fn, &impl)
		}
	}
}

func fastCall0(cfn uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall1(cfn, a1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, 0, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall2(cfn, a1, a2 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall3(cfn, a1, a2, a3 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall4(cfn, a1, a2, a3, a4 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall5(cfn, a1, a2, a3, a4, a5 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall6(cfn, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall7(cfn, a1, a2, a3, a4, a5, a6, a7 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

func fastCall8(cfn, a1, a2, a3, a4, a5, a6, a7, a8 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall9 calls a C function with 9 arguments without heap allocations.
func FastCall9(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9 = a9
	s.a10, s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall10 calls a C function with 10 arguments without heap allocations.
func FastCall10(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9, s.a10 = a9, a10
	s.a11, s.a12, s.a13, s.a14, s.a15 = 0, 0, 0, 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// FastCall13 calls a C function with 13 arguments without heap allocations.
func FastCall13(cfn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, a4, a5, a6, a7, a8
	s.a9, s.a10, s.a11, s.a12, s.a13 = a9, a10, a11, a12, a13
	s.a14, s.a15 = 0, 0
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = 0, 0, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// Exported aliases for zero-allocation FFI calls.

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
func FastCall5(cfn, a1, a2, a3, a4, a5 uintptr) uintptr { return fastCall5(cfn, a1, a2, a3, a4, a5) }

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

// FastCallNF1 variants call C functions with N integer/pointer args and 1 float32 arg
// passed via the first float register (f1). The float arg is passed as a uintptr
// containing the IEEE 754 bit pattern (use math.Float32bits to convert).

// FastCall3F1 calls a C function with 3 integer args and 1 float arg.
func FastCall3F1(cfn, a1, a2, a3, f1 uintptr) uintptr {
	s := thePool.Get().(*syscall15Args)
	s.fn = cfn
	s.a1, s.a2, s.a3, s.a4, s.a5, s.a6, s.a7, s.a8 = a1, a2, a3, 0, 0, 0, 0, 0
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
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, 0, 0, 0, 0, 0, 0, 0
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
	s.f1, s.f2, s.f3, s.f4, s.f5, s.f6, s.f7, s.f8 = f1, f2, 0, 0, 0, 0, 0, 0
	s.arm64_r8 = 0
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(s))
	r := s.a1
	thePool.Put(s)
	return r
}

// forceFuncPtr assigns the function pointer in src to the reflect.Value dst,
// bypassing type checking. This is safe when the source and destination function
// types have identical ABI layouts (same number of pointer-sized args/returns).
func forceFuncPtr(dst reflect.Value, src any) {
	dstPtr := dst.Addr().UnsafePointer() // *func(...)
	srcPtr := (*unsafe.Pointer)(reflect.ValueOf(src).UnsafePointer())
	*(*unsafe.Pointer)(dstPtr) = *srcPtr
}
