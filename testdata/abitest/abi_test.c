// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

#include <assert.h>
#include <inttypes.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

uint32_t stack_uint8_t(uint32_t a, uint32_t b, uint32_t c, uint32_t d, uint32_t e, uint32_t f, uint32_t g, uint32_t h, uint8_t i, uint8_t j, uint32_t k ) {
    assert(i == 1);
    assert(j == 2);
    assert(k == 1024);
    return a | b | c | d | e | f | g | h | (uint32_t)i | (uint32_t)j | k;
}

uint32_t reg_uint8_t(uint8_t a, uint8_t b, uint32_t c) {
    assert(a == 1);
    assert(b == 2);
    assert(c == 1024);
    return a | b | c;
}

uint32_t stack_string(uint32_t a, uint32_t b, uint32_t c, uint32_t d, uint32_t e, uint32_t f, uint32_t g, uint32_t h, const char * i) {
    assert(i != 0);
    assert(strcmp(i, "test") == 0);
    return a | b | c | d | e | f | g | h;
}

void stack_8i32_3strings(char* result, size_t size, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, const char* s1, const char* s2, const char* s3) {
    snprintf(result, size, "%d:%d:%d:%d:%d:%d:%d:%d:%s:%s:%s", a1, a2, a3, a4, a5, a6, a7, a8, s1, s2, s3);
}

// HFA (Homogeneous Float Aggregate) struct with 2 floats
typedef struct {
    float x;
    float y;
} Float2;

// HFA struct with 4 floats
typedef struct {
    float x;
    float y;
    float z;
    float w;
} Float4;

// Non-HFA struct (mixed types)
typedef struct {
    int32_t a;
    float b;
} MixedStruct;

// Small struct that fits in one register
typedef struct {
    int32_t x;
    int32_t y;
} IntPair;

// Test: 8 int registers exhausted, then HFA struct on stack
void test_8int_hfa2_stack(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, Float2 f) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%.1f:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, f.x, f.y);
}

// Test: 8 int registers exhausted, then multiple structs on stack
void test_8int_2structs_stack(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, IntPair p1, IntPair p2) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, p1.x, p1.y, p2.x, p2.y);
}

// Test: 8 float registers exhausted, then HFA on stack
void test_8float_hfa2_stack(char *buf, size_t bufsize, float f1, float f2, float f3, float f4, float f5, float f6, float f7, float f8, Float2 f) {
    snprintf(buf, bufsize, "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f.x, f.y);
}

// Test: mixed - int regs exhausted, float struct can still use float regs
void test_8int_hfa2_floatregs(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, Float2 f) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%.1f:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, f.x, f.y);
}

// Test: primitives and struct interleaved on stack
void test_8int_int_struct_int(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, int32_t a9, IntPair p, int32_t a10) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9, p.x, p.y, a10);
}

// Test: HFA4 struct on stack (4 floats)
void test_8int_hfa4_stack(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, Float4 f) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%.1f:%.1f:%.1f:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, f.x, f.y, f.z, f.w);
}

// Test: mixed type struct on stack
void test_8int_mixed_struct(char *buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4, int32_t a5, int32_t a6, int32_t a7, int32_t a8, MixedStruct m) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%d:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, m.a, m.b);
}
