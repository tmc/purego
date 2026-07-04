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

int64_t sum14_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5,
                int64_t a6, int64_t a7, int64_t a8, int64_t a9, int64_t a10,
                int64_t a11, int64_t a12, int64_t a13, int64_t a14) {
  return a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8 + a9 + a10 + a11 + a12 + a13 +
         a14;
}

int64_t sum15_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5,
                int64_t a6, int64_t a7, int64_t a8, int64_t a9, int64_t a10,
                int64_t a11, int64_t a12, int64_t a13, int64_t a14,
                int64_t a15) {
  return a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8 + a9 + a10 + a11 + a12 + a13 +
         a14 + a15;
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

typedef int64_t (*callback14_t)(int64_t, int64_t, int64_t, int64_t, int64_t,
                                int64_t, int64_t, int64_t, int64_t, int64_t,
                                int64_t, int64_t, int64_t, int64_t);
int64_t call_callback14(callback14_t cb, int64_t a1, int64_t a2, int64_t a3,
                        int64_t a4, int64_t a5, int64_t a6, int64_t a7,
                        int64_t a8, int64_t a9, int64_t a10, int64_t a11,
                        int64_t a12, int64_t a13, int64_t a14) {
  return cb(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14);
}

typedef int64_t (*callback15_t)(int64_t, int64_t, int64_t, int64_t, int64_t,
                                int64_t, int64_t, int64_t, int64_t, int64_t,
                                int64_t, int64_t, int64_t, int64_t, int64_t);
int64_t call_callback15(callback15_t cb, int64_t a1, int64_t a2, int64_t a3,
                        int64_t a4, int64_t a5, int64_t a6, int64_t a7,
                        int64_t a8, int64_t a9, int64_t a10, int64_t a11,
                        int64_t a12, int64_t a13, int64_t a14, int64_t a15) {
  return cb(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15);
}

// order5_c weights its 5 integer args by position so argument ordering is checked.
int64_t order5_c(int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5) {
	return a1 + a2 * 10 + a3 * 100 + a4 * 1000 + a5 * 10000;
}

// mixed_int_float_c interleaves int and float params; SyscallN cannot call this
// correctly but CallN can, since it passes ints and floats in separate registers.
int64_t mixed_int_float_c(int64_t a1, double f1, int64_t a2, double f2) {
	return a1 + (int64_t)(f1 * 10) + a2 * 100 + (int64_t)(f2 * 1000);
}

// deref_u64 reads and returns the uint64 at address p. Used by the CallN
// pointer-argument test.
uint64_t deref_u64(uintptr_t p) {
	return *(uint64_t *)p;
}

// ret_double_c and ret_float_c return a floating-point value, used to test that
// CallN recovers a floating-point return value.
double ret_double_c(double a, double b) { return a * 100.0 + b; }

float ret_float_c(float a, float b) { return a * 100.0f + b; }
