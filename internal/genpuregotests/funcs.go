// Copyright 2025 The Ebitengine Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"
	"text/template"
)

// Template functions available in txtar templates
var templateFuncs = template.FuncMap{
	"expandSpec": expandSpec,
}

// expandSpec converts a test specification into a full test function with all parameters expanded
func expandSpec(spec TestSpec) TestFunc {
	tf := TestFunc{Name: spec.Name}
	paramIdx := 0
	for _, pspec := range spec.Params {
		parts := strings.Split(pspec, ":")
		if len(parts) != 2 {
			panic("invalid param spec: " + pspec)
		}
		var count int
		fmt.Sscanf(parts[1], "%d", &count)
		info := typeInfo[parts[0]]
		for i := 0; i < count; i++ {
			paramIdx++
			name := "a"
			switch parts[0] {
			case "ptr":
				name = "p"
			case "string":
				name = "s"
			case "bool":
				name = "b"
			case "float", "double":
				name = "f"
			}
			tf.Params = append(tf.Params, Param{
				Type: info.CType, Name: fmt.Sprintf("%s%d", name, paramIdx),
				Format: info.Format, FormatCast: info.FormatCast,
			})
		}
	}
	return tf
}

// Type information for C and Go interop - maps type names to their representations and test value generators
var typeInfo = map[string]struct {
	CType, GoType, Format, FormatCast string
	TestValue                         func(int) string
}{
	"int8":   {"int8_t", "int8", "d", "int", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"int16":  {"int16_t", "int16", "d", "int", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"int32":  {"int32_t", "int32", "d", "int", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"int64":  {"int64_t", "int64", "lld", "long long", func(i int) string { return fmt.Sprintf("int64(%d)", i+1) }},
	"uint8":  {"uint8_t", "uint8", "u", "unsigned", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"uint16": {"uint16_t", "uint16", "u", "unsigned", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"uint32": {"uint32_t", "uint32", "u", "unsigned", func(i int) string { return fmt.Sprintf("%d", i+1) }},
	"uint64": {"uint64_t", "uint64", "llu", "unsigned long long", func(i int) string { return fmt.Sprintf("uint64(%d)", i+1) }},
	"float":  {"float", "float32", "f", "double", func(i int) string { return fmt.Sprintf("float32(%d)", i+1) }},
	"double": {"double", "float64", "lf", "double", func(i int) string { return fmt.Sprintf("float64(%d)", i+1) }},
	"bool": {"bool", "bool", "d", "int", func(i int) string {
		if i%2 == 1 {
			return "true"
		}
		return "false"
	}},
	"ptr": {"void*", "unsafe.Pointer", "p", "void*", func(i int) string {
		if i%3 == 0 {
			return "nil"
		}
		return fmt.Sprintf("unsafe.Pointer(uintptr(%d))", i)
	}},
	"string": {"const char*", "string", "s", "const char*", func(i int) string {
		strs := []string{`"hello"`, `"world"`, `"foo"`, `"bar"`, `"baz"`, `"qux"`, `"quux"`, `"corge"`, `"grault"`, `"garply"`, `"waldo"`, `"fred"`}
		if i < len(strs) {
			return strs[i]
		}
		return fmt.Sprintf(`"s%d"`, i)
	}},
}

// Param represents a single function parameter with C/Go type mappings
type Param struct {
	Type, Name, Format, FormatCast string
}

// TestFunc represents a complete test function with all parameters
type TestFunc struct {
	Name   string
	Params []Param
}

// CallbackReturnType returns the C return type for a callback function
func (tf TestFunc) CallbackReturnType() string {
	// Most callbacks return int for testing purposes
	return "int"
}

// CallbackGoReturnType returns the Go return type for a callback function
func (tf TestFunc) CallbackGoReturnType() string {
	// Most callbacks return int for testing purposes
	return "int"
}

// ExpectedReturnValue returns the expected return value for callback tests
func (tf TestFunc) ExpectedReturnValue() string {
	// Use a distinctive return value to verify callback execution
	return "42"
}

// GoType returns the Go type string for this parameter (used in templates)
func (p Param) GoType() string {
	for _, info := range typeInfo {
		if info.CType == p.Type {
			return info.GoType
		}
	}
	return p.Type
}

// TestValue returns the Go test value for this parameter at the given index (used in templates)
func (p Param) TestValue(i int) string {
	for _, info := range typeInfo {
		if info.CType == p.Type {
			return info.TestValue(i)
		}
	}
	return fmt.Sprintf("%d", i+1)
}

// TestValueC returns the C literal value for this parameter at the given index (used in callback templates)
func (p Param) TestValueC(i int) string {
	switch p.Type {
	case "int8_t", "int16_t", "int32_t", "uint8_t", "uint16_t", "uint32_t":
		return fmt.Sprintf("%d", i+1)
	case "int64_t":
		return fmt.Sprintf("%dLL", i+1)
	case "uint64_t":
		return fmt.Sprintf("%dULL", i+1)
	case "float":
		return fmt.Sprintf("%d.0f", i+1)
	case "double":
		return fmt.Sprintf("%d.0", i+1)
	case "bool":
		if i%2 == 1 {
			return "1"
		}
		return "0"
	case "void*":
		if i%3 == 0 {
			return "NULL"
		}
		return fmt.Sprintf("(void*)%d", i)
	case "const char*":
		strs := []string{`"hello"`, `"world"`, `"foo"`, `"bar"`, `"baz"`, `"qux"`, `"quux"`, `"corge"`, `"grault"`, `"garply"`, `"waldo"`, `"fred"`}
		if i < len(strs) {
			return strs[i]
		}
		return fmt.Sprintf(`"s%d"`, i)
	}
	return fmt.Sprintf("%d", i+1)
}

// ExpectedOutput returns the expected C output string for this parameter at the given index (used in templates)
func (p Param) ExpectedOutput(i int) string {
	switch p.Type {
	case "int8_t", "int16_t", "int32_t", "uint8_t", "uint16_t", "uint32_t", "int64_t", "uint64_t":
		return fmt.Sprintf("%d", i+1)
	case "float", "double":
		return fmt.Sprintf("%f", float64(i+1))
	case "bool":
		if i%2 == 1 {
			return "1"
		}
		return "0"
	case "void*":
		if i%3 == 0 {
			return "(nil)"
		}
		return fmt.Sprintf("0x%x", i)
	case "const char*":
		strs := []string{"hello", "world", "foo", "bar", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred"}
		if i < len(strs) {
			return strs[i]
		}
		return fmt.Sprintf("s%d", i)
	}
	return fmt.Sprintf("%d", i+1)
}
