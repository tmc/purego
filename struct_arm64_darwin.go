// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

//go:build darwin && arm64

package purego

import (
	"reflect"
	"unsafe"
)

// placeRegisters implements Darwin ARM64 calling convention for struct arguments.
//
// For HFA/HVA structs, each element must go in a separate register (or stack slot for elements
// that don't fit in registers). We use placeRegistersAArch64 for this.
//
// For non-HFA/HVA structs, Darwin uses byte-level packing. We copy the struct memory in
// 8-byte chunks, which works correctly for both register and stack placement.
func placeRegisters(v reflect.Value, addFloat func(uintptr), addInt func(uintptr)) {
	// Check if this is an HFA/HVA
	hfa := isHFA(v.Type())
	hva := isHVA(v.Type())

	// For HFA/HVA structs, use the standard ARM64 logic which places each element separately
	if hfa || hva {
		placeRegistersAArch64(v, addFloat, addInt)
		return
	}

	// For non-HFA/HVA structs, use byte-level copying
	// If the value is not addressable, create an addressable copy
	if !v.CanAddr() {
		addressable := reflect.New(v.Type()).Elem()
		addressable.Set(v)
		v = addressable
	}
	ptr := unsafe.Pointer(v.Addr().Pointer())
	size := v.Type().Size()

	// Copy the struct memory in 8-byte chunks
	for offset := uintptr(0); offset < size; offset += 8 {
		// Read 8 bytes (or whatever remains) from the struct
		var chunk uintptr
		remaining := size - offset
		if remaining >= 8 {
			chunk = *(*uintptr)(unsafe.Add(ptr, offset))
		} else {
			// For the last partial chunk, read only the remaining bytes
			bytes := (*[8]byte)(unsafe.Add(ptr, offset))
			for i := uintptr(0); i < remaining; i++ {
				chunk |= uintptr(bytes[i]) << (i * 8)
			}
		}
		addInt(chunk)
	}
}
