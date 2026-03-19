// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

#include <stdint.h>

int64_t sum1_c(int64_t a1) { return a1; }

int64_t sum2_c(int64_t a1, int64_t a2) { return a1 + a2; }

int64_t sum3_c(int64_t a1, int64_t a2, int64_t a3) { return a1 + a2 + a3; }

int64_t sum5_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5) {
  return a1 + a2 + a3 + a4 + a5;
}

int64_t sum10_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5,
                int64_t a6, int64_t a7, int64_t a8, int64_t a9, int64_t a10) {
  return a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8 + a9 + a10;
}

// weighted_sum functions: N integer args + 1 trailing float
int64_t weighted_sum3f_c(int64_t a1, int64_t a2, int64_t a3, float w) {
  return (int64_t)((double)(a1 + a2 + a3) * (double)w);
}

int64_t weighted_sum5f_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4,
                         int64_t a5, float w) {
  return (int64_t)((double)(a1 + a2 + a3 + a4 + a5) * (double)w);
}

int64_t weighted_sum3d_c(int64_t a1, int64_t a2, int64_t a3, double w) {
  return (int64_t)((double)(a1 + a2 + a3) * w);
}

// interleaved float: int, float, int, int -> weighted sum
int64_t interleaved_if_c(int64_t a, float w, int64_t b, int64_t c) {
  return (int64_t)((double)(a + b + c) * (double)w);
}

// interleaved: int, float, int, float, int -> a*w1 + b*w2 + c
int64_t interleaved_2f_c(int64_t a, float w1, int64_t b, float w2, int64_t c) {
  return (int64_t)((double)a * (double)w1 + (double)b * (double)w2 + (double)c);
}

// mlx-go hot-path shapes: interleaved float at position 5 (SDPA/Rope shape)
// 9 args: 4 ints, 1 float, 4 ints -> int32
int32_t sdpa_shape_c(int64_t q, int64_t k, int64_t v, int64_t mask,
                     float scale, int64_t mode, int64_t arr, int64_t sinks,
                     int64_t stream) {
  return (int32_t)(q + k + v + mask + (int64_t)scale + mode + arr + sinks +
                   stream);
}

// 10 args: all ints -> int32 (QuantizedMatmul shape)
int32_t qmm_shape_c(int64_t res, int64_t w, int64_t scales, int64_t biases,
                     int64_t transpose, int64_t group_size, int64_t bits,
                     int64_t type_, int64_t x, int64_t stream) {
  return (int32_t)(res + w + scales + biases + transpose + group_size + bits +
                   type_ + x + stream);
}

// 5 args: 3 ints, 1 float, 1 int -> int32 (RMSNorm shape)
int32_t rmsnorm_shape_c(int64_t res, int64_t x, int64_t weight, float eps,
                        int64_t stream) {
  return (int32_t)(res + x + weight + (int64_t)eps + stream);
}

// 3 args: float at pos 0, 2 ints -> weighted sum (edge case: float first)
int64_t interleaved_3_f0_c(float w, int64_t a, int64_t b) {
  return (int64_t)((double)(a + b) * (double)w);
}

typedef int64_t (*callback1_t)(int64_t);
int64_t call_callback1(callback1_t cb, int64_t a1) { return cb(a1); }

typedef int64_t (*callback2_t)(int64_t, int64_t);
int64_t call_callback2(callback2_t cb, int64_t a1, int64_t a2) {
  return cb(a1, a2);
}

typedef int64_t (*callback3_t)(int64_t, int64_t, int64_t);
int64_t call_callback3(callback3_t cb, int64_t a1, int64_t a2, int64_t a3) {
  return cb(a1, a2, a3);
}

typedef int64_t (*callback5_t)(int64_t, int64_t, int64_t, int64_t, int64_t);
int64_t call_callback5(callback5_t cb, int64_t a1, int64_t a2, int64_t a3,
                       int64_t a4, int64_t a5) {
  return cb(a1, a2, a3, a4, a5);
}

typedef int64_t (*callback10_t)(int64_t, int64_t, int64_t, int64_t, int64_t,
                                int64_t, int64_t, int64_t, int64_t, int64_t);
int64_t call_callback10(callback10_t cb, int64_t a1, int64_t a2, int64_t a3,
                        int64_t a4, int64_t a5, int64_t a6, int64_t a7,
                        int64_t a8, int64_t a9, int64_t a10) {
  return cb(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10);
}

