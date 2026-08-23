// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package purego

import (
	"reflect"
	"runtime"
)

// Float16 carries the raw bits of an IEEE 754 binary16 value.
//
// RegisterFunc supports Float16 scalar arguments on Apple ARM64 platforms.
// Float16 stores bits and does not implement binary16 arithmetic. Convert its
// bits explicitly when numeric operations are needed.
type Float16 uint16

var float16Type = reflect.TypeFor[Float16]()

func isFloat16Type(ty reflect.Type) bool {
	return ty == float16Type
}

func isFloatABIType(ty reflect.Type) bool {
	return isFloat16Type(ty) || ty.Kind() == reflect.Float32 || ty.Kind() == reflect.Float64
}

func supportsFloat16Arguments() bool {
	return (runtime.GOOS == "darwin" || runtime.GOOS == "ios") && runtime.GOARCH == "arm64"
}

func containsFloat16Value(ty reflect.Type) bool {
	if isFloat16Type(ty) {
		return true
	}
	switch ty.Kind() {
	case reflect.Array:
		return containsFloat16Value(ty.Elem())
	case reflect.Struct:
		for i := 0; i < ty.NumField(); i++ {
			if containsFloat16Value(ty.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

func rejectFloat16CallbackType(ty reflect.Type) {
	if ty == nil || ty.Kind() != reflect.Func {
		return
	}
	for i := 0; i < ty.NumIn(); i++ {
		if containsFloat16Value(ty.In(i)) {
			panic("purego: Float16 callbacks are not supported")
		}
	}
	for i := 0; i < ty.NumOut(); i++ {
		if containsFloat16Value(ty.Out(i)) {
			panic("purego: Float16 callbacks are not supported")
		}
	}
}
