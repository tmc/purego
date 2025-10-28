// Copyright 2025 The Ebitengine Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate go run .

package main

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/txtar"
)

//go:embed *.txtar
var templatesFS embed.FS

func main() {
	tmpl := template.Must(template.New("").Funcs(templateFuncs).ParseFS(templatesFS, "*.txtar"))

	for _, t := range tmpl.Templates() {
		if t.Name() == "" {
			continue
		}
		var buf strings.Builder
		t.Execute(&buf, testSpecs)
		for _, file := range txtar.Parse([]byte(buf.String())).Files {
			os.WriteFile(filepath.Join("..", "..", file.Name), file.Data, 0644)
		}
	}
}

// TestSpec defines a test function specification
type TestSpec struct {
	Name   string
	Params []string // e.g., ["int32:11", "float:3", "string:2"]
}

// Test Specifications - Template context
// Top-level keys are arbitrary identifiers used in templates
// Second-level keys are Go test names that will be generated
// Values are arrays of TestSpecs that become sub-tests
var testSpecs = map[string]any{
	"abitest": map[string][]TestSpec{
		"TestABIGenerated": {
			// Basic integer tests
			{"test_11_int8", []string{"int8:11"}},
			{"test_11_int16", []string{"int16:11"}},
			{"test_11_int32", []string{"int32:11"}},
			{"test_11_int64", []string{"int64:11"}},
			{"test_9_int32", []string{"int32:9"}},
			{"test_15_int32", []string{"int32:15"}},

			// Bool tests
			{"test_11_bool", []string{"bool:11"}},

			// Pointer tests
			// {"test_9_ptrs", []string{"ptr:9"}},

			// Float tests
			// TODO: Float32 packing has issues
			// {"test_11_float32", []string{"float:11"}},
			{"test_11_float64", []string{"double:11"}},

			// Mixed tests - safe combinations that work with current code
			{"test_mixed_8r_2u8_1u32", []string{"uint32:8", "uint8:2", "uint32:1"}},
			{"test_mixed_8i32_3i16", []string{"int32:8", "int16:3"}},
			{"test_mixed_varied", []string{"int32:8", "int8:1", "int16:1", "int32:1"}},
			{"test_mixed_8i32_3bool", []string{"int32:8", "bool:3"}},
			{"test_mixed_8i32_3f64", []string{"int32:8", "double:3"}}, // float64 works
			// TODO: Alternating types has issues
			{"test_alternating_i32_bool", []string{"int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1"}},

			// TODO: These tests expose bugs in func.go:295 where floats go to stack when int registers are full
			// instead of using independent float registers D0-D7. Uncomment after fixing the OR logic bug.
			{"test_mixed_8i32_3f32", []string{"int32:8", "float:3"}},
			{"test_mixed_8i32_8f32", []string{"int32:8", "float:8"}},
			{"test_mixed_8i32_1f32_1i32", []string{"int32:8", "float:1", "int32:1"}},

			// Edge cases
			{"test_8i32_u8_u32", []string{"int32:8", "uint8:1", "uint32:1"}},
			{"test_8i32_u8_u16_u32", []string{"int32:8", "uint8:1", "uint16:1", "uint32:1"}},
			{"test_8i32_u8_u64", []string{"int32:8", "uint8:1", "uint64:1"}},
			{"test_8i32_4u8", []string{"int32:8", "uint8:4"}},
			{"test_8i32_3u16", []string{"int32:8", "uint16:3"}},
			{"test_8i32_u8_u64_u8_u64", []string{"int32:8", "uint8:1", "uint64:1", "uint8:1", "uint64:1"}},

			// String tests - verify each string gets correct value (not duplicates)
			{"test_3strings", []string{"string:3"}},                                              // Simple test: "foo", "bar", "baz"
			{"test_2i32_3strings", []string{"int32:2", "string:3"}},                              // Clearer: gives "foo", "bar", "baz"
			{"test_8i32_3strings", []string{"int32:8", "string:3"}},                              // Strings on stack after ints
			{"test_mixed_stack_strings", []string{"int32:8", "string:1", "int32:1", "string:1"}}, // Interleaved
			{"test_9_strings", []string{"string:9"}},                                             // Multiple strings - tests val.String() vs v.String() bug
			{"test_string_first_stack", []string{"int32:8", "string:1"}},
			{"test_5strings", []string{"string:5"}}, // Simpler multi-string test
		},
	},
}
