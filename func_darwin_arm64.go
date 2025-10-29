// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

//go:build darwin && arm64

package purego

import "reflect"

// shouldBundleStackArgs determines if we need to start C-style packing for
// Darwin ARM64 stack arguments. This happens when registers are exhausted.
func shouldBundleStackArgs(v reflect.Value, numInts, numFloats int) bool {
	// Check primitives first
	isFloat := v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64
	isInt := !isFloat && v.Kind() != reflect.Struct
	primitiveOnStack := (isInt && numInts >= numOfIntegerRegisters()) ||
		(isFloat && numFloats >= numOfFloatRegisters)

	// Check if struct would go on stack
	structOnStack := false
	if v.Kind() == reflect.Struct {
		hfa := isHFA(v.Type())
		hva := isHVA(v.Type())
		size := v.Type().Size()

		if hfa || hva || size <= 16 {
			if hfa && numFloats+v.NumField() > numOfFloatRegisters {
				structOnStack = true
			} else if hva && numInts+v.NumField() > numOfIntegerRegisters() {
				structOnStack = true
			} else if size <= 16 {
				slotsNeeded := int((size + 7) / 8)
				if numInts+slotsNeeded > numOfIntegerRegisters() {
					structOnStack = true
				}
			}
		}
	}

	return primitiveOnStack || structOnStack
}

// structFitsInRegisters determines if a struct can still fit in remaining
// registers, used during stack argument bundling to decide if a struct
// should go through normal register allocation or be bundled with stack args.
func structFitsInRegisters(val reflect.Value, tempNumInts, tempNumFloats int) (bool, int, int) {
	hfa := isHFA(val.Type())
	hva := isHVA(val.Type())
	size := val.Type().Size()

	if hfa {
		// HFA: check if elements fit in float registers
		if tempNumFloats+val.NumField() <= numOfFloatRegisters {
			return true, tempNumInts, tempNumFloats + val.NumField()
		}
	} else if hva {
		// HVA: check if elements fit in int registers
		if tempNumInts+val.NumField() <= numOfIntegerRegisters() {
			return true, tempNumInts + val.NumField(), tempNumFloats
		}
	} else if size <= 16 {
		// Non-HFA/HVA small structs use int registers for byte-packing
		slotsNeeded := int((size + 7) / 8)
		if tempNumInts+slotsNeeded <= numOfIntegerRegisters() {
			return true, tempNumInts + slotsNeeded, tempNumFloats
		}
	}

	return false, tempNumInts, tempNumFloats
}
