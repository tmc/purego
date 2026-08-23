// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

#include <stdint.h>
#include <string.h>

uint16_t float16_bits(_Float16 value) {
    uint16_t bits;
    memcpy(&bits, &value, sizeof(bits));
    return bits;
}

uint64_t float16_mixed_bank(uintptr_t x0, uintptr_t x1, uintptr_t x2,
                            uintptr_t x3, uintptr_t x4, _Float16 value) {
    uint16_t bits;
    memcpy(&bits, &value, sizeof(bits));
    return x0 ^ (x1 << 1) ^ (x2 << 2) ^ (x3 << 3) ^ (x4 << 4) ^ bits;
}

uint64_t uint16_mixed_bank(uintptr_t x0, uintptr_t x1, uintptr_t x2,
                           uintptr_t x3, uintptr_t x4, uint16_t value) {
    return x0 ^ (x1 << 1) ^ (x2 << 2) ^ (x3 << 3) ^ (x4 << 4) ^ value;
}

uint64_t float16_interleaved(uintptr_t x0, _Float16 h0, uintptr_t x1,
                             _Float16 h1) {
    uint16_t b0, b1;
    memcpy(&b0, &h0, sizeof(b0));
    memcpy(&b1, &h1, sizeof(b1));
    return x0 ^ (x1 << 1) ^ ((uint64_t)b0 << 32) ^ ((uint64_t)b1 << 48);
}

uint64_t float16_nine(_Float16 h0, _Float16 h1, _Float16 h2, _Float16 h3,
                      _Float16 h4, _Float16 h5, _Float16 h6, _Float16 h7,
                      _Float16 h8) {
    _Float16 values[] = {h0, h1, h2, h3, h4, h5, h6, h7, h8};
    uint64_t hash = 0;
    for (unsigned i = 0; i < 9; i++) {
        uint16_t bits;
        memcpy(&bits, &values[i], sizeof(bits));
        hash = hash * 65599 + bits;
    }
    return hash;
}
