// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || windows

package purego

import (
	"math"
	"unsafe"
)

// FastFunc0 is a specialized wrapper for zero-argument C functions returning uintptr.
// Common use case: constructors like mlx_array_new().
type FastFunc0 func() uintptr

// NewFastFunc0 creates a fast function wrapper for signature:
//
//	uintptr_t fn(void)
func NewFastFunc0(cfn uintptr) FastFunc0 {
	return func() uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFunc1 is a specialized high-performance wrapper for single-argument C functions.
// It bypasses reflect.MakeFunc entirely for maximum performance.
// Usage:
//   var fastDtype = FastFunc1(sym)
//   result := fastDtype(arrayPtr)
type FastFunc1 func(uintptr) uintptr

// NewFastFunc1 creates a fast function wrapper for a C function with signature:
//   uintptr_t fn(uintptr_t arg)
// This is ~10x faster than RegisterFunc for this common pattern.
func NewFastFunc1(cfn uintptr) FastFunc1 {
	return func(arg uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFunc2 is a specialized wrapper for two-argument C functions.
type FastFunc2 func(uintptr, uintptr) uintptr

// NewFastFunc2 creates a fast function wrapper for signature:
//   uintptr_t fn(uintptr_t arg1, uintptr_t arg2)
func NewFastFunc2(cfn uintptr) FastFunc2 {
	return func(arg1, arg2 uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFunc3 is a specialized wrapper for three-argument C functions.
type FastFunc3 func(uintptr, uintptr, uintptr) uintptr

// NewFastFunc3 creates a fast function wrapper for signature:
//   uintptr_t fn(uintptr_t arg1, uintptr_t arg2, uintptr_t arg3)
func NewFastFunc3(cfn uintptr) FastFunc3 {
	return func(arg1, arg2, arg3 uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFunc4 is a specialized wrapper for four-argument C functions.
type FastFunc4 func(uintptr, uintptr, uintptr, uintptr) uintptr

// NewFastFunc4 creates a fast function wrapper for signature:
//   uintptr_t fn(uintptr_t arg1, uintptr_t arg2, uintptr_t arg3, uintptr_t arg4)
func NewFastFunc4(cfn uintptr) FastFunc4 {
	return func(arg1, arg2, arg3, arg4 uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		syscall.a4 = arg4
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncVoid1 is for void-returning single-argument functions.
type FastFuncVoid1 func(uintptr)

// NewFastFuncVoid1 creates a fast wrapper for: void fn(uintptr_t arg)
func NewFastFuncVoid1(cfn uintptr) FastFuncVoid1 {
	return func(arg uintptr) {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		thePool.Put(syscall)
	}
}

// FastFuncVoid2 is for void-returning two-argument functions.
type FastFuncVoid2 func(uintptr, uintptr)

// NewFastFuncVoid2 creates a fast wrapper for: void fn(uintptr_t arg1, uintptr_t arg2)
func NewFastFuncVoid2(cfn uintptr) FastFuncVoid2 {
	return func(arg1, arg2 uintptr) {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		thePool.Put(syscall)
	}
}

// FastFuncVoid3 is for void-returning three-argument functions.
type FastFuncVoid3 func(uintptr, uintptr, uintptr)

// NewFastFuncVoid3 creates a fast wrapper for: void fn(uintptr_t, uintptr_t, uintptr_t)
func NewFastFuncVoid3(cfn uintptr) FastFuncVoid3 {
	return func(arg1, arg2, arg3 uintptr) {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		thePool.Put(syscall)
	}
}

// FastFuncVoid4 is for void-returning four-argument functions.
type FastFuncVoid4 func(uintptr, uintptr, uintptr, uintptr)

// NewFastFuncVoid4 creates a fast wrapper for: void fn(uintptr_t, uintptr_t, uintptr_t, uintptr_t)
func NewFastFuncVoid4(cfn uintptr) FastFuncVoid4 {
	return func(arg1, arg2, arg3, arg4 uintptr) {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		syscall.a4 = arg4
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		thePool.Put(syscall)
	}
}

// ============================================================================
// Status-returning variants (int32 return type)
// These are the most common MLX function patterns.
// ============================================================================

// FastFuncStatus1 is for functions returning int32 with one uintptr argument.
// Common use case: mlx_array_eval(arr).
type FastFuncStatus1 func(uintptr) int32

// NewFastFuncStatus1 creates a fast wrapper for: int32_t fn(uintptr_t)
func NewFastFuncStatus1(cfn uintptr) FastFuncStatus1 {
	return func(arg uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus2 is for functions returning int32 with out-pointer + one arg.
// Common use case: mlx_array_item_*(out, arr).
type FastFuncStatus2 func(unsafe.Pointer, uintptr) int32

// NewFastFuncStatus2 creates a fast wrapper for: int32_t fn(T* out, uintptr_t arr)
func NewFastFuncStatus2(cfn uintptr) FastFuncStatus2 {
	return func(out unsafe.Pointer, arg uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus3 is for functions returning int32 with out-pointer + two args.
// Common use case: mlx_closure_apply(out, closure, inputs).
type FastFuncStatus3 func(unsafe.Pointer, uintptr, uintptr) int32

// NewFastFuncStatus3 creates a fast wrapper for: int32_t fn(T* out, uintptr_t, uintptr_t)
func NewFastFuncStatus3(cfn uintptr) FastFuncStatus3 {
	return func(out unsafe.Pointer, arg1, arg2 uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg1
		syscall.a3 = arg2
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus4 is for functions returning int32 with out-pointer + three args.
// Common use case: mlx_matmul(out, a, b, stream) - binary operations.
type FastFuncStatus4 func(unsafe.Pointer, uintptr, uintptr, uintptr) int32

// NewFastFuncStatus4 creates a fast wrapper for: int32_t fn(T* out, uintptr_t, uintptr_t, uintptr_t)
func NewFastFuncStatus4(cfn uintptr) FastFuncStatus4 {
	return func(out unsafe.Pointer, arg1, arg2, arg3 uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg1
		syscall.a3 = arg2
		syscall.a4 = arg3
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus5 is for functions returning int32 with out-pointer + four args.
// Common use case: ternary operations with stream parameter.
type FastFuncStatus5 func(unsafe.Pointer, uintptr, uintptr, uintptr, uintptr) int32

// NewFastFuncStatus5 creates a fast wrapper for: int32_t fn(T* out, uintptr_t, uintptr_t, uintptr_t, uintptr_t)
func NewFastFuncStatus5(cfn uintptr) FastFuncStatus5 {
	return func(out unsafe.Pointer, arg1, arg2, arg3, arg4 uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg1
		syscall.a3 = arg2
		syscall.a4 = arg3
		syscall.a5 = arg4
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus6 is for functions returning int32 with out-pointer + five args.
type FastFuncStatus6 func(unsafe.Pointer, uintptr, uintptr, uintptr, uintptr, uintptr) int32

// NewFastFuncStatus6 creates a fast wrapper for: int32_t fn(T* out, 5 x uintptr_t)
func NewFastFuncStatus6(cfn uintptr) FastFuncStatus6 {
	return func(out unsafe.Pointer, arg1, arg2, arg3, arg4, arg5 uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg1
		syscall.a3 = arg2
		syscall.a4 = arg3
		syscall.a5 = arg4
		syscall.a6 = arg5
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// Mixed type variants (int32, uint, bool parameters)
// ============================================================================

// FastFuncStatusInt3 is for: int32_t fn(T* out, uintptr_t arr, int32_t val)
// Common use case: mlx_array_dim(arr, axis).
type FastFuncStatusInt3 func(unsafe.Pointer, uintptr, int32) int32

// NewFastFuncStatusInt3 creates a fast wrapper.
func NewFastFuncStatusInt3(cfn uintptr) FastFuncStatusInt3 {
	return func(out unsafe.Pointer, arg uintptr, val int32) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg
		syscall.a3 = uintptr(val)
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusIntStream is for: int32_t fn(T* out, uintptr_t arr, int32_t val, uintptr_t stream)
// Common use case: mlx_reshape(out, arr, axis, stream).
type FastFuncStatusIntStream func(unsafe.Pointer, uintptr, int32, uintptr) int32

// NewFastFuncStatusIntStream creates a fast wrapper.
func NewFastFuncStatusIntStream(cfn uintptr) FastFuncStatusIntStream {
	return func(out unsafe.Pointer, arg uintptr, val int32, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg
		syscall.a3 = uintptr(val)
		syscall.a4 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusBoolStream is for: int32_t fn(T* out, uintptr_t arr, bool val, uintptr_t stream)
type FastFuncStatusBoolStream func(unsafe.Pointer, uintptr, bool, uintptr) int32

// NewFastFuncStatusBoolStream creates a fast wrapper.
func NewFastFuncStatusBoolStream(cfn uintptr) FastFuncStatusBoolStream {
	return func(out unsafe.Pointer, arg uintptr, val bool, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arg
		var b uintptr
		if val {
			b = 1
		}
		syscall.a3 = b
		syscall.a4 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// Axis operation variants (common MLX pattern with shape/axes arrays)
// ============================================================================

// FastFuncAxisOp is for: int32_t fn(T* out, uintptr_t arr, int32_t* axes, size_t naxes, bool keepdims, uintptr_t stream)
// Common use case: mlx_sum_axes, mlx_mean_axes, etc.
type FastFuncAxisOp func(unsafe.Pointer, uintptr, unsafe.Pointer, uint, bool, uintptr) int32

// NewFastFuncAxisOp creates a fast wrapper for axis reduction operations.
func NewFastFuncAxisOp(cfn uintptr) FastFuncAxisOp {
	return func(out unsafe.Pointer, arr uintptr, axes unsafe.Pointer, naxes uint, keepdims bool, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		syscall.a3 = uintptr(axes)
		syscall.a4 = uintptr(naxes)
		var kd uintptr
		if keepdims {
			kd = 1
		}
		syscall.a5 = kd
		syscall.a6 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// Float parameter variants (using float registers f1-f8)
// ============================================================================

// FastFuncFloat1 is for functions with a single float64 argument returning uintptr.
// Common use case: constructors that take a scalar value.
type FastFuncFloat1 func(float64) uintptr

// NewFastFuncFloat1 creates a fast wrapper for: uintptr_t fn(double)
func NewFastFuncFloat1(cfn uintptr) FastFuncFloat1 {
	return func(f float64) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.f1 = uintptr(math.Float64bits(f))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncFloat32_1 is for functions with a single float32 argument returning uintptr.
type FastFuncFloat32_1 func(float32) uintptr

// NewFastFuncFloat32_1 creates a fast wrapper for: uintptr_t fn(float)
func NewFastFuncFloat32_1(cfn uintptr) FastFuncFloat32_1 {
	return func(f float32) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		// Float32 in lower 32 bits of float register
		syscall.f1 = uintptr(math.Float32bits(f))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncFloat32_2 is for functions with two float32 arguments returning uintptr.
// Common use case: mlx_array_new_complex(real, imag).
type FastFuncFloat32_2 func(float32, float32) uintptr

// NewFastFuncFloat32_2 creates a fast wrapper for: uintptr_t fn(float, float)
func NewFastFuncFloat32_2(cfn uintptr) FastFuncFloat32_2 {
	return func(f1, f2 float32) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.f1 = uintptr(math.Float32bits(f1))
		syscall.f2 = uintptr(math.Float32bits(f2))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusFloat32 is for: int32_t fn(T* out, float val)
type FastFuncStatusFloat32 func(unsafe.Pointer, float32) int32

// NewFastFuncStatusFloat32 creates a fast wrapper.
func NewFastFuncStatusFloat32(cfn uintptr) FastFuncStatusFloat32 {
	return func(out unsafe.Pointer, f float32) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.f1 = uintptr(math.Float32bits(f))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusFloat32_2 is for: int32_t fn(T* out, float f1, float f2)
// Common use case: mlx_array_set_complex(out, real, imag).
type FastFuncStatusFloat32_2 func(unsafe.Pointer, float32, float32) int32

// NewFastFuncStatusFloat32_2 creates a fast wrapper.
func NewFastFuncStatusFloat32_2(cfn uintptr) FastFuncStatusFloat32_2 {
	return func(out unsafe.Pointer, f1, f2 float32) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.f1 = uintptr(math.Float32bits(f1))
		syscall.f2 = uintptr(math.Float32bits(f2))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// uint-returning variants (for size/count functions)
// ============================================================================

// FastFuncUint1 is for functions returning uint with one uintptr argument.
// Common use case: mlx_array_size(arr), mlx_array_ndim(arr).
type FastFuncUint1 func(uintptr) uint

// NewFastFuncUint1 creates a fast wrapper for: size_t fn(uintptr_t)
func NewFastFuncUint1(cfn uintptr) FastFuncUint1 {
	return func(arg uintptr) uint {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := uint(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncInt32_1 is for functions returning int32 with one uintptr argument.
// Common use case: mlx_array_dim(arr, axis) returning dimension size.
type FastFuncInt32_1 func(uintptr) int32

// NewFastFuncInt32_1 creates a fast wrapper for: int32_t fn(uintptr_t)
func NewFastFuncInt32_1(cfn uintptr) FastFuncInt32_1 {
	return func(arg uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncInt32_2 is for functions returning int32 with two uintptr arguments.
type FastFuncInt32_2 func(uintptr, int32) int32

// NewFastFuncInt32_2 creates a fast wrapper for: int32_t fn(uintptr_t, int32_t)
func NewFastFuncInt32_2(cfn uintptr) FastFuncInt32_2 {
	return func(arg uintptr, val int32) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		syscall.a2 = uintptr(val)
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// bool-returning variants
// ============================================================================

// FastFuncBool0 is for functions returning bool with no arguments.
type FastFuncBool0 func() bool

// NewFastFuncBool0 creates a fast wrapper for: bool fn(void)
func NewFastFuncBool0(cfn uintptr) FastFuncBool0 {
	return func() bool {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1 != 0
		thePool.Put(syscall)
		return r
	}
}

// FastFuncBool2 is for functions returning bool with two uintptr arguments.
// Common use case: mlx_device_equal(dev1, dev2).
type FastFuncBool2 func(uintptr, uintptr) bool

// NewFastFuncBool2 creates a fast wrapper for: bool fn(uintptr_t, uintptr_t)
func NewFastFuncBool2(cfn uintptr) FastFuncBool2 {
	return func(arg1, arg2 uintptr) bool {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1 != 0
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// Pointer-returning variants
// ============================================================================

// FastFuncPtr1 is for functions returning a typed pointer with one uintptr argument.
// Common use case: mlx_array_shape(arr) returning *int32.
type FastFuncPtr1 func(uintptr) unsafe.Pointer

// NewFastFuncPtr1 creates a fast wrapper for: T* fn(uintptr_t)
func NewFastFuncPtr1(cfn uintptr) FastFuncPtr1 {
	return func(arg uintptr) unsafe.Pointer {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := unsafe.Pointer(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// Extended argument count variants (5-8 args)
// ============================================================================

// FastFunc5 is a specialized wrapper for five-argument C functions.
type FastFunc5 func(uintptr, uintptr, uintptr, uintptr, uintptr) uintptr

// NewFastFunc5 creates a fast function wrapper for signature:
//
//	uintptr_t fn(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t)
func NewFastFunc5(cfn uintptr) FastFunc5 {
	return func(arg1, arg2, arg3, arg4, arg5 uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		syscall.a4 = arg4
		syscall.a5 = arg5
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFunc6 is a specialized wrapper for six-argument C functions.
type FastFunc6 func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr

// NewFastFunc6 creates a fast function wrapper for signature:
//
//	uintptr_t fn(6 x uintptr_t)
func NewFastFunc6(cfn uintptr) FastFunc6 {
	return func(arg1, arg2, arg3, arg4, arg5, arg6 uintptr) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = arg1
		syscall.a2 = arg2
		syscall.a3 = arg3
		syscall.a4 = arg4
		syscall.a5 = arg5
		syscall.a6 = arg6
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// ============================================================================
// MLX-specific variants (handle common MLX patterns)
// ============================================================================

// FastFuncStatusBool4 is for: int32_t fn(*T out, uintptr arr, bool keepdims, uintptr stream)
// Common use case: mlx_all, mlx_any without axes.
type FastFuncStatusBool4 func(unsafe.Pointer, uintptr, bool, uintptr) int32

// NewFastFuncStatusBool4 creates a fast wrapper.
func NewFastFuncStatusBool4(cfn uintptr) FastFuncStatusBool4 {
	return func(out unsafe.Pointer, arr uintptr, keepdims bool, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		var kd uintptr
		if keepdims {
			kd = 1
		}
		syscall.a3 = kd
		syscall.a4 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusInt5 is for: int32_t fn(*T out, uintptr arr, int32 axis, bool keepdims, uintptr stream)
// Common use case: mlx_all_axis, mlx_any_axis, mlx_argmax_axis.
type FastFuncStatusInt5 func(unsafe.Pointer, uintptr, int32, bool, uintptr) int32

// NewFastFuncStatusInt5 creates a fast wrapper.
func NewFastFuncStatusInt5(cfn uintptr) FastFuncStatusInt5 {
	return func(out unsafe.Pointer, arr uintptr, axis int32, keepdims bool, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		syscall.a3 = uintptr(axis)
		var kd uintptr
		if keepdims {
			kd = 1
		}
		syscall.a4 = kd
		syscall.a5 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusInt4 is for: int32_t fn(*T out, uintptr arr, int32 val, uintptr stream)
// Common use case: mlx_argsort_axis, mlx_argpartition.
type FastFuncStatusInt4 func(unsafe.Pointer, uintptr, int32, uintptr) int32

// NewFastFuncStatusInt4 creates a fast wrapper.
func NewFastFuncStatusInt4(cfn uintptr) FastFuncStatusInt4 {
	return func(out unsafe.Pointer, arr uintptr, val int32, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		syscall.a3 = uintptr(val)
		syscall.a4 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusInt2_5 is for: int32_t fn(*T out, uintptr arr, int32 a, int32 b, uintptr stream)
// Common use case: mlx_argpartition_axis.
type FastFuncStatusInt2_5 func(unsafe.Pointer, uintptr, int32, int32, uintptr) int32

// NewFastFuncStatusInt2_5 creates a fast wrapper.
func NewFastFuncStatusInt2_5(cfn uintptr) FastFuncStatusInt2_5 {
	return func(out unsafe.Pointer, arr uintptr, a, b int32, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		syscall.a3 = uintptr(a)
		syscall.a4 = uintptr(b)
		syscall.a5 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusBool5 is for: int32_t fn(*T out, uintptr a, uintptr b, bool val, uintptr stream)
// Common use case: mlx_array_equal.
type FastFuncStatusBool5 func(unsafe.Pointer, uintptr, uintptr, bool, uintptr) int32

// NewFastFuncStatusBool5 creates a fast wrapper.
func NewFastFuncStatusBool5(cfn uintptr) FastFuncStatusBool5 {
	return func(out unsafe.Pointer, a, b uintptr, val bool, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = a
		syscall.a3 = b
		var v uintptr
		if val {
			v = 1
		}
		syscall.a4 = v
		syscall.a5 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatus2IntStream is for: int32_t fn(*T out, uintptr arr1, uintptr arr2, int32 val, uintptr stream)
// Common use case: mlx_take_along_axis.
type FastFuncStatus2IntStream func(unsafe.Pointer, uintptr, uintptr, int32, uintptr) int32

// NewFastFuncStatus2IntStream creates a fast wrapper.
func NewFastFuncStatus2IntStream(cfn uintptr) FastFuncStatus2IntStream {
	return func(out unsafe.Pointer, arr1, arr2 uintptr, val int32, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr1
		syscall.a3 = arr2
		syscall.a4 = uintptr(val)
		syscall.a5 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncInt1 is for functions returning int32 with one int32 argument.
// Common use case: mlx_array_new_int(int32).
type FastFuncInt1 func(int32) uintptr

// NewFastFuncInt1 creates a fast wrapper for: uintptr_t fn(int32_t)
func NewFastFuncInt1(cfn uintptr) FastFuncInt1 {
	return func(val int32) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(val)
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncBool1 is for functions with one bool argument returning uintptr.
// Common use case: mlx_array_new_bool(bool).
type FastFuncBool1 func(bool) uintptr

// NewFastFuncBool1 creates a fast wrapper for: uintptr_t fn(bool)
func NewFastFuncBool1(cfn uintptr) FastFuncBool1 {
	return func(val bool) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		var v uintptr
		if val {
			v = 1
		}
		syscall.a1 = v
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncFloat64_1 is for functions with a single float64 argument returning uintptr.
// Common use case: mlx_array_new_float64.
type FastFuncFloat64_1 func(float64) uintptr

// NewFastFuncFloat64_1 creates a fast wrapper for: uintptr_t fn(double)
func NewFastFuncFloat64_1(cfn uintptr) FastFuncFloat64_1 {
	return func(f float64) uintptr {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.f1 = uintptr(math.Float64bits(f))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := syscall.a1
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusFloat64 is for: int32_t fn(*T out, float64 val)
// Common use case: mlx_array_set_float64.
type FastFuncStatusFloat64 func(unsafe.Pointer, float64) int32

// NewFastFuncStatusFloat64 creates a fast wrapper.
func NewFastFuncStatusFloat64(cfn uintptr) FastFuncStatusFloat64 {
	return func(out unsafe.Pointer, f float64) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.f1 = uintptr(math.Float64bits(f))
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusBool2 is for: int32_t fn(*T out, bool val)
// Common use case: mlx_array_set_bool.
type FastFuncStatusBool2 func(unsafe.Pointer, bool) int32

// NewFastFuncStatusBool2 creates a fast wrapper.
func NewFastFuncStatusBool2(cfn uintptr) FastFuncStatusBool2 {
	return func(out unsafe.Pointer, val bool) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		var v uintptr
		if val {
			v = 1
		}
		syscall.a2 = v
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusInt2 is for: int32_t fn(*T out, int32 val)
// Common use case: mlx_array_set_int.
type FastFuncStatusInt2 func(unsafe.Pointer, int32) int32

// NewFastFuncStatusInt2 creates a fast wrapper.
func NewFastFuncStatusInt2(cfn uintptr) FastFuncStatusInt2 {
	return func(out unsafe.Pointer, val int32) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = uintptr(val)
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}

// FastFuncStatusAxes is for: int32_t fn(*T out, uintptr arr, *int32 axes, uint len, uintptr stream)
// Common use case: mlx_transpose_axes, mlx_expand_dims_axes.
type FastFuncStatusAxes func(unsafe.Pointer, uintptr, unsafe.Pointer, uint, uintptr) int32

// NewFastFuncStatusAxes creates a fast wrapper.
func NewFastFuncStatusAxes(cfn uintptr) FastFuncStatusAxes {
	return func(out unsafe.Pointer, arr uintptr, axes unsafe.Pointer, axesLen uint, stream uintptr) int32 {
		syscall := thePool.Get().(*syscall15Args)
		syscall.fn = cfn
		syscall.a1 = uintptr(out)
		syscall.a2 = arr
		syscall.a3 = uintptr(axes)
		syscall.a4 = uintptr(axesLen)
		syscall.a5 = stream
		runtime_cgocall(syscall20XABI0, unsafe.Pointer(syscall))
		r := int32(syscall.a1)
		thePool.Put(syscall)
		return r
	}
}
