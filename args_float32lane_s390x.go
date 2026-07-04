// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build s390x && linux

package purego

// float32ArgBits places a single-precision argument's bits within a
// floating-point argument register, and float32ResultBits extracts a
// single-precision result from a floating-point result register. On big-endian
// s390x a single-precision value occupies the high 32 bits of the 64-bit register.
func float32ArgBits(bits uint32) uintptr   { return uintptr(bits) << 32 }
func float32ResultBits(reg uintptr) uint32 { return uint32(reg >> 32) }
