// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build !s390x && !386 && !arm && !ppc64le && (darwin || freebsd || linux || netbsd)

package purego

// float32ArgBits places a single-precision argument's bits within a
// floating-point argument register, and float32ResultBits extracts a
// single-precision result from a floating-point result register. On the
// little-endian architectures the value occupies the low 32 bits of the register.
func float32ArgBits(bits uint32) uintptr   { return uintptr(bits) }
func float32ResultBits(reg uintptr) uint32 { return uint32(reg) }
