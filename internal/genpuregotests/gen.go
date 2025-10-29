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

// StructSpec defines a struct test specification (imported from funcs.go)

// Test Specifications - Template context
// Top-level keys are arbitrary identifiers used in templates
// Second-level keys are Go test names that will be generated
// Values are arrays of TestSpecs that become sub-tests
var testSpecs = map[string]any{
	"abitest": map[string][]TestSpec{
		"TestABIGenerated": {
			// MINIMAL FAILING TESTS - One test per unique failure pattern

			// Pattern 1: Stack int - last param becomes 0
			{"test_11_int32", []string{"int32:11"}}, // Last param: got 0, want 11

			// Pattern 2: Stack int - position 12 gets wrong value (14 instead of 12)
			{"test_12_int32", []string{"int32:12"}}, // Position 12: got 14, want 12

			// Pattern 3: Mixed int+float on stack - both corrupted
			{"test_mixed_8i32_1f32_1i32", []string{"int32:8", "float:1", "int32:1"}}, // Float=0, int=garbage

			// Pattern 4: String duplication - ints then strings
			{"test_8i32_3strings", []string{"int32:8", "string:3"}}, // All strings become "grault"

			// Pattern 5: String duplication - interleaved (int between strings)
			{"test_mixed_stack_strings", []string{"int32:8", "string:1", "int32:1", "string:1"}}, // Both strings same

			// Pattern 6: String duplication - only strings
			{"test_string_only_10", []string{"string:10"}}, // Last 2 strings duplicated

			// Pattern 7: String duplication - few ints then many strings
			{"test_1i32_10strings", []string{"int32:1", "string:10"}}, // Last 3 strings duplicated

			// Pattern 8: String duplication - bools then strings
			{"test_mixed_bool_string", []string{"bool:8", "string:3"}}, // All strings become "grault"

			// Pattern 9: Float then int corruption
			{"test_float_then_int", []string{"float:8", "int32:8"}}, // Ints after floats corrupted

			// Pattern 10: Mixed float types (float32 + float64)
			{"test_mixed_f32_f64_mix", []string{"float:4", "double:4", "float:2"}}, // Last float=0

			// MORE EXAMPLES of same patterns (uncomment to test variations):
			// // Same as pattern 1 (last param = 0):
			// {"test_11_uint32", []string{"uint32:11"}},
			// {"test_alternating_i32_bool", []string{"int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1", "bool:1", "int32:1"}},
			//
			// // Same as pattern 2 (position 12 corrupted):
			// {"test_15_int32", []string{"int32:15"}},
			//
			// // Same as pattern 4 (ints then strings):
			// {"test_8i32_2strings", []string{"int32:8", "string:2"}},
			//
			// // Same as pattern 9 (float then int):
			// {"test_double_then_int", []string{"double:8", "int32:8"}},

			// INTERESTING PASSING TESTS (regression checks - uncomment to verify):
			// Stack boundary cases that work:
			// {"test_9_int32", []string{"int32:9"}},          // Exactly one int on stack - WORKS
			// {"test_10_int32", []string{"int32:10"}},        // Two ints on stack - WORKS
			// {"test_11_int64", []string{"int64:11"}},        // int64 overflow - WORKS
			// {"test_11_uint64", []string{"uint64:11"}},      // uint64 overflow - WORKS
			//
			// Float/int mixing that works:
			// {"test_9_float32", []string{"float:9"}},                                  // One float32 on stack - WORKS
			// {"test_9_float64", []string{"double:9"}},                                 // One float64 on stack - WORKS
			// {"test_mixed_5i32_5f32_5i32", []string{"int32:5", "float:5", "int32:5"}}, // Multiple type transitions - WORKS
			// {"test_mixed_8i32_3f32", []string{"int32:8", "float:3"}},                // 8 ints + 3 floats - WORKS
			// {"test_mixed_8i32_3f64", []string{"int32:8", "double:3"}},               // 8 ints + 3 float64s - WORKS
			//
			// String cases that work:
			// {"test_3strings", []string{"string:3"}},                     // 3 strings in registers - WORKS
			// {"test_9_strings", []string{"string:9"}},                    // 9 strings (some on stack) - WORKS
			// {"test_mixed_8i32_1string_1i32", []string{"int32:8", "string:1", "int32:1"}}, // String between ints - WORKS
			//
			// All types mixed:
			// {"test_mixed_all_types", []string{"int32:3", "float:2", "bool:2", "string:2", "int32:2"}}, // WORKS
		},
	},
}
