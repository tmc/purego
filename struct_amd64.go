// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package purego

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/ebitengine/purego/internal/strings"
)

func getStruct(outType reflect.Type, syscall syscall15Args) (v reflect.Value) {
	if hasStringFields(outType) {
		return getStructWithStrings(outType, syscall)
	}
	outSize := outType.Size()
	switch {
	case outSize == 0:
		return reflect.New(outType).Elem()
	case outSize <= 8:
		if isAllFloats(outType) {
			// 2 float32s or 1 float64s are return in the float register
			return reflect.NewAt(outType, unsafe.Pointer(&struct{ a uintptr }{syscall.f1})).Elem()
		}
		// up to 8 bytes is returned in RAX
		return reflect.NewAt(outType, unsafe.Pointer(&struct{ a uintptr }{syscall.a1})).Elem()
	case outSize <= 16:
		r1, r2 := syscall.a1, syscall.a2
		if isAllFloats(outType) {
			r1 = syscall.f1
			r2 = syscall.f2
		} else {
			// check first 8 bytes if it's floats
			hasFirstFloat := false
			f1 := outType.Field(0).Type
			if f1.Kind() == reflect.Float64 || f1.Kind() == reflect.Float32 && outType.Field(1).Type.Kind() == reflect.Float32 {
				r1 = syscall.f1
				hasFirstFloat = true
			}

			// find index of the field that starts the second 8 bytes
			var i int
			for i = 0; i < outType.NumField(); i++ {
				if outType.Field(i).Offset == 8 {
					break
				}
			}

			// check last 8 bytes if they are floats
			f1 = outType.Field(i).Type
			if f1.Kind() == reflect.Float64 || f1.Kind() == reflect.Float32 && i+1 == outType.NumField() {
				r2 = syscall.f1
			} else if hasFirstFloat {
				// if the first field was a float then that means the second integer field
				// comes from the first integer register
				r2 = syscall.a1
			}
		}
		return reflect.NewAt(outType, unsafe.Pointer(&struct{ a, b uintptr }{r1, r2})).Elem()
	default:
		// create struct from the Go pointer created above
		// weird pointer dereference to circumvent go vet
		return reflect.NewAt(outType, *(*unsafe.Pointer)(unsafe.Pointer(&syscall.a1))).Elem()
	}
}

// getStructWithStrings handles struct return values that contain string fields.
// C returns char* pointers which must be converted to Go strings.
func getStructWithStrings(outType reflect.Type, syscall syscall15Args) reflect.Value {
	v := reflect.New(outType)
	goPtr := v.UnsafePointer()
	cSize := calculateCStructSize(outType)

	var regs [2]uintptr
	var ptr unsafe.Pointer
	if cSize <= 16 {
		regs = [2]uintptr{syscall.a1, syscall.a2}
		ptr = unsafe.Pointer(&regs[0])
	} else {
		// large struct returned via hidden pointer
		ptr = *(*unsafe.Pointer)(unsafe.Pointer(&syscall.a1))
	}

	// Walk struct fields and reconstruct from C layout
	cOffset := uintptr(0)
	for i := 0; i < outType.NumField(); i++ {
		field := outType.Field(i)
		var fieldSize, fieldAlign uintptr
		if field.Type.Kind() == reflect.String {
			fieldSize = 8
			fieldAlign = 8
		} else {
			fieldSize = field.Type.Size()
			fieldAlign = uintptr(field.Type.Align())
		}
		cOffset = (cOffset + fieldAlign - 1) &^ (fieldAlign - 1)

		fieldPtr := unsafe.Add(ptr, cOffset)
		if field.Type.Kind() == reflect.String {
			charPtr := *(*uintptr)(fieldPtr)
			goStr := strings.GoString(charPtr)
			*(*string)(unsafe.Add(goPtr, field.Offset)) = goStr
		} else {
			// Copy raw bytes from C layout to Go struct at the correct Go offset
			dst := unsafe.Add(goPtr, field.Offset)
			for j := uintptr(0); j < fieldSize; j++ {
				*(*byte)(unsafe.Add(dst, j)) = *(*byte)(unsafe.Add(fieldPtr, j))
			}
		}
		cOffset += fieldSize
	}
	return v.Elem()
}

func isAllFloats(ty reflect.Type) bool {
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		switch f.Type.Kind() {
		case reflect.Float64, reflect.Float32:
		default:
			return false
		}
	}
	return true
}

// https://refspecs.linuxbase.org/elf/x86_64-abi-0.99.pdf
// https://gitlab.com/x86-psABIs/x86-64-ABI
// Class determines where the 8 byte value goes.
// Higher value classes win over lower value classes
const (
	_NO_CLASS = 0b0000
	_SSE      = 0b0001
	_X87      = 0b0011 // long double not used in Go
	_INTEGER  = 0b0111
	_MEMORY   = 0b1111
)

func addStruct(v reflect.Value, numInts, numFloats, numStack *int, addInt, addFloat, addStack func(uintptr), keepAlive []any) []any {
	if v.Type().Size() == 0 {
		return keepAlive
	}

	// if greater than 64 bytes place on stack
	if v.Type().Size() > 8*8 {
		return placeStack(v, addStack, keepAlive)
	}
	var (
		savedNumFloats = *numFloats
		savedNumInts   = *numInts
		savedNumStack  = *numStack
	)
	var ok bool
	ok, keepAlive = tryPlaceRegister(v, addFloat, addInt, keepAlive)
	placeOnStack := postMerger(v.Type()) || !ok
	if placeOnStack {
		// reset any values placed in registers
		*numFloats = savedNumFloats
		*numInts = savedNumInts
		*numStack = savedNumStack
		keepAlive = placeStack(v, addStack, keepAlive)
	}
	return keepAlive
}

func postMerger(t reflect.Type) (passInMemory bool) {
	// (c) If the size of the aggregate exceeds two eightbytes and the first eight- byte isn’t SSE or any other
	// eightbyte isn’t SSEUP, the whole argument is passed in memory.
	if t.Kind() != reflect.Struct {
		return false
	}
	if t.Size() <= 2*8 {
		return false
	}
	return true // Go does not have an SSE/SSEUP type so this is always true
}

func tryPlaceRegister(v reflect.Value, addFloat func(uintptr), addInt func(uintptr), keepAlive []any) (ok bool, newKeepAlive []any) {
	ok = true
	newKeepAlive = keepAlive
	var val uint64
	var shift byte // # of bits to shift
	var flushed bool
	class := _NO_CLASS
	flushIfNeeded := func() {
		if flushed {
			return
		}
		flushed = true
		if class == _SSE {
			addFloat(uintptr(val))
		} else {
			addInt(uintptr(val))
		}
		val = 0
		shift = 0
		class = _NO_CLASS
	}
	var place func(v reflect.Value)
	place = func(v reflect.Value) {
		var numFields int
		if v.Kind() == reflect.Struct {
			numFields = v.Type().NumField()
		} else {
			numFields = v.Type().Len()
		}

		for i := 0; i < numFields; i++ {
			flushed = false
			var f reflect.Value
			if v.Kind() == reflect.Struct {
				f = v.Field(i)
			} else {
				f = v.Index(i)
			}
			switch f.Kind() {
			case reflect.Struct:
				place(f)
			case reflect.Bool:
				if f.Bool() {
					val |= 1 << shift
				}
				shift += 8
				class |= _INTEGER
			case reflect.Pointer, reflect.UnsafePointer:
				val = uint64(f.Pointer())
				shift = 64
				class = _INTEGER
			case reflect.String:
				ptr := strings.CString(f.String())
				newKeepAlive = append(newKeepAlive, ptr)
				val = uint64(uintptr(unsafe.Pointer(ptr)))
				shift = 64
				class = _INTEGER
			case reflect.Int8:
				val |= uint64(f.Int()&0xFF) << shift
				shift += 8
				class |= _INTEGER
			case reflect.Int16:
				val |= uint64(f.Int()&0xFFFF) << shift
				shift += 16
				class |= _INTEGER
			case reflect.Int32:
				val |= uint64(f.Int()&0xFFFF_FFFF) << shift
				shift += 32
				class |= _INTEGER
			case reflect.Int64, reflect.Int:
				val = uint64(f.Int())
				shift = 64
				class = _INTEGER
			case reflect.Uint8:
				val |= f.Uint() << shift
				shift += 8
				class |= _INTEGER
			case reflect.Uint16:
				val |= f.Uint() << shift
				shift += 16
				class |= _INTEGER
			case reflect.Uint32:
				val |= f.Uint() << shift
				shift += 32
				class |= _INTEGER
			case reflect.Uint64, reflect.Uint, reflect.Uintptr:
				val = f.Uint()
				shift = 64
				class = _INTEGER
			case reflect.Float32:
				val |= uint64(math.Float32bits(float32(f.Float()))) << shift
				shift += 32
				class |= _SSE
			case reflect.Float64:
				if v.Type().Size() > 16 {
					ok = false
					return
				}
				val = uint64(math.Float64bits(f.Float()))
				shift = 64
				class = _SSE
			case reflect.Array:
				place(f)
			default:
				panic("purego: unsupported kind " + f.Kind().String())
			}

			if shift == 64 {
				flushIfNeeded()
			} else if shift > 64 {
				// Should never happen, but may if we forget to reset shift after flush (or forget to flush),
				// better fall apart here, than corrupt arguments.
				panic("purego: tryPlaceRegisters shift > 64")
			}
		}
	}

	place(v)
	flushIfNeeded()
	return ok, newKeepAlive
}

func placeStack(v reflect.Value, addStack func(uintptr), keepAlive []any) []any {
	if hasStringFields(v.Type()) {
		return placeStackWithStrings(v, addStack, keepAlive)
	}
	// Copy the struct as a contiguous block of memory in eightbyte (8-byte)
	// chunks. The x86-64 ABI requires structs passed on the stack to be
	// laid out exactly as in memory, including padding and field packing
	// within eightbytes. Decomposing field-by-field would place each field
	// as a separate stack slot, breaking structs with mixed-type fields
	// that share an eightbyte (e.g. int32 + float32).
	if !v.CanAddr() {
		tmp := reflect.New(v.Type()).Elem()
		tmp.Set(v)
		v = tmp
	}
	ptr := v.Addr().UnsafePointer()
	size := v.Type().Size()
	for off := uintptr(0); off < size; off += 8 {
		chunk := *(*uintptr)(unsafe.Add(ptr, off))
		addStack(chunk)
	}
	return keepAlive
}

// placeStackWithStrings builds a C-layout buffer for structs with string fields,
// converting Go strings to C char* pointers, then copies eightbyte chunks to the stack.
func placeStackWithStrings(v reflect.Value, addStack func(uintptr), keepAlive []any) []any {
	// Make sure the struct is addressable so we can read fields via raw pointer
	if !v.CanAddr() {
		tmp := reflect.New(v.Type()).Elem()
		tmp.Set(v)
		v = tmp
	}
	structPtr := v.Addr().UnsafePointer()

	cSize := calculateCStructSize(v.Type())
	buf := make([]byte, cSize)
	cOffset := uintptr(0)

	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i)
		var fieldSize, fieldAlign uintptr
		if field.Type.Kind() == reflect.String {
			fieldSize = 8
			fieldAlign = 8
		} else {
			fieldSize = field.Type.Size()
			fieldAlign = uintptr(field.Type.Align())
		}
		cOffset = (cOffset + fieldAlign - 1) &^ (fieldAlign - 1)

		dst := unsafe.Pointer(&buf[cOffset])
		if field.Type.Kind() == reflect.String {
			goStr := *(*string)(unsafe.Add(structPtr, field.Offset))
			ptr := strings.CString(goStr)
			keepAlive = append(keepAlive, ptr)
			*(*uintptr)(dst) = uintptr(unsafe.Pointer(ptr))
		} else {
			// Copy field bytes from Go struct memory using field offset
			src := unsafe.Add(structPtr, field.Offset)
			for j := uintptr(0); j < fieldSize; j++ {
				*(*byte)(unsafe.Add(dst, j)) = *(*byte)(unsafe.Add(src, j))
			}
		}
		cOffset += fieldSize
	}

	// Copy buffer as eightbyte chunks
	for off := uintptr(0); off < cSize; off += 8 {
		chunk := *(*uintptr)(unsafe.Pointer(&buf[off]))
		addStack(chunk)
	}
	return keepAlive
}

func placeRegisters(v reflect.Value, addFloat func(uintptr), addInt func(uintptr), keepAlive []any) []any {
	panic("purego: placeRegisters not implemented on amd64")
}

// shouldBundleStackArgs always returns false on non-Darwin platforms
// since C-style stack argument bundling is only needed on Darwin ARM64.
func shouldBundleStackArgs(v reflect.Value, numInts, numFloats int) bool {
	return false
}

// structFitsInRegisters is not used on amd64.
func structFitsInRegisters(val reflect.Value, tempNumInts, tempNumFloats int) (bool, int, int) {
	panic("purego: structFitsInRegisters should not be called on amd64")
}

// collectStackArgs is not used on amd64.
func collectStackArgs(args []reflect.Value, startIdx int, numInts, numFloats int,
	keepAlive []any, addInt, addFloat, addStack func(uintptr),
	pNumInts, pNumFloats, pNumStack *int) ([]reflect.Value, []any) {
	panic("purego: collectStackArgs should not be called on amd64")
}

// bundleStackArgs is not used on amd64.
func bundleStackArgs(stackArgs []reflect.Value, addStack func(uintptr)) {
	panic("purego: bundleStackArgs should not be called on amd64")
}
