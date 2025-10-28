// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

#include <stdint.h>
#include <stdbool.h>
#include <assert.h>
#include <stdio.h>
#include <string.h>

uint32_t stack_uint8_t(uint32_t a, uint32_t b, uint32_t c, uint32_t d, uint32_t e, uint32_t f, uint32_t g, uint32_t h, uint8_t i, uint8_t j, uint32_t k ) {
    assert(i == 1);
    assert(j == 2);
    assert(k == 1024);
    return a | b | c | d | e | f | g | h | (uint32_t) i | (uint32_t) j | k;
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

const char* test_8i32_3strings(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                                int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                                const char* s1, const char* s2, const char* s3) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%s:%s:%s",
             a1, a2, a3, a4, a5, a6, a7, a8, s1, s2, s3);
    return result;
}

const char* test_8i32_3f32_independent_regs(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                                             int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                                             float f1, float f2, float f3) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%.1f:%.1f:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, f1, f2, f3);
    return result;
}

const char* test_11_float32_packing(float f1, float f2, float f3, float f4, float f5,
                                    float f6, float f7, float f8, float f9, float f10, float f11) {
    static char result[256];
    snprintf(result, sizeof(result), "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11);
    return result;
}

const char* test_alternating_i32_bool(int32_t a1, bool b1, int32_t a2, bool b2,
                                      int32_t a3, bool b3, int32_t a4, bool b4,
                                      int32_t a5, bool b5, int32_t a6) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, b1, a2, b2, a3, b3, a4, b4, a5, b5, a6);
    return result;
}
