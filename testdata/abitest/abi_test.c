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

void test_8i32_3strings(char* result, size_t size, int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                        int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                        const char* s1, const char* s2, const char* s3) {
    snprintf(result, size, "%d:%d:%d:%d:%d:%d:%d:%d:%s:%s:%s",
             a1, a2, a3, a4, a5, a6, a7, a8, s1, s2, s3);
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

const char* test_9_int32(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                         int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                         int32_t a9) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9);
    return result;
}

void test_10_int32(char* buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                   int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                   int32_t a9, int32_t a10) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9, a10);
}

void test_11_int32(char* buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                   int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                   int32_t a9, int32_t a10, int32_t a11) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11);
}

const char* test_9_int8(int8_t a1, int8_t a2, int8_t a3, int8_t a4,
                        int8_t a5, int8_t a6, int8_t a7, int8_t a8,
                        int8_t a9) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9);
    return result;
}

const char* test_9_bool(bool b1, bool b2, bool b3, bool b4,
                        bool b5, bool b6, bool b7, bool b8,
                        bool b9) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d",
             b1, b2, b3, b4, b5, b6, b7, b8, b9);
    return result;
}

const char* test_9_float32(float f1, float f2, float f3, float f4,
                           float f5, float f6, float f7, float f8,
                           float f9) {
    static char result[256];
    snprintf(result, sizeof(result), "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9);
    return result;
}

void test_10_float32(char* buf, size_t bufsize, float f1, float f2, float f3, float f4,
                     float f5, float f6, float f7, float f8,
                     float f9, float f10) {
    snprintf(buf, bufsize, "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9, f10);
}

const char* test_11_float32(float f1, float f2, float f3, float f4,
                            float f5, float f6, float f7, float f8,
                            float f9, float f10, float f11) {
    static char result[256];
    snprintf(result, sizeof(result), "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11);
    return result;
}

const char* test_12_int32(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                          int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                          int32_t a9, int32_t a10, int32_t a11, int32_t a12) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12);
    return result;
}

const char* test_13_int32(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                          int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                          int32_t a9, int32_t a10, int32_t a11, int32_t a12,
                          int32_t a13) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13);
    return result;
}

const char* test_12_float32(float f1, float f2, float f3, float f4,
                            float f5, float f6, float f7, float f8,
                            float f9, float f10, float f11, float f12) {
    static char result[256];
    snprintf(result, sizeof(result), "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12);
    return result;
}

const char* test_13_float32(float f1, float f2, float f3, float f4,
                            float f5, float f6, float f7, float f8,
                            float f9, float f10, float f11, float f12,
                            float f13) {
    static char result[256];
    snprintf(result, sizeof(result), "%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f",
             f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13);
    return result;
}

const char* test_10_intermixed(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                                int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                                float f1, float f2, float f3, float f4,
                                float f5, float f6, float f7, float f8,
                                int32_t a9, float f9) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%.1f:%d:%.1f",
             a1, a2, a3, a4, a5, a6, a7, a8, f1, f2, f3, f4, f5, f6, f7, f8, a9, f9);
    return result;
}

const char* test_mixed_stack(int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                              int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                              bool b1, const char* s1, int32_t a9) {
    static char result[256];
    snprintf(result, sizeof(result), "%d:%d:%d:%d:%d:%d:%d:%d:%d:%s:%d",
             a1, a2, a3, a4, a5, a6, a7, a8, b1, s1, a9);
    return result;
}

void test_mixed_stack_4args(char* buf, size_t bufsize, int32_t a1, int32_t a2, int32_t a3, int32_t a4,
                             int32_t a5, int32_t a6, int32_t a7, int32_t a8,
                             const char* s1, bool b1, int32_t a9, const char* s2) {
    snprintf(buf, bufsize, "%d:%d:%d:%d:%d:%d:%d:%d:%s:%d:%d:%s",
             a1, a2, a3, a4, a5, a6, a7, a8, s1, b1, a9, s2);
}
