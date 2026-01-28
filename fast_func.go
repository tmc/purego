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

// forceFuncPtr assigns the function pointer in src to the reflect.Value dst,
// bypassing type checking. This is safe when the source and destination function
// types have identical ABI layouts (same number of pointer-sized args/returns).
func forceFuncPtr(dst reflect.Value, src any) {
	dstPtr := dst.Addr().UnsafePointer() // *func(...)
	srcPtr := (*unsafe.Pointer)(reflect.ValueOf(src).UnsafePointer())
	*(*unsafe.Pointer)(dstPtr) = *srcPtr
}
