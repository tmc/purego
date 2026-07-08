# Design: zero-allocation call primitives for purego

Status: proposal / request for comment
Relates to: #399, #445, and the mixed int/float limitation noted in `syscall.go`

## Summary

purego callers on hot paths want to call C functions with **no per-call heap
allocation** and **correct handling of floating-point arguments**. `SyscallN`
gives neither: its variadic slice escapes to the heap across a module boundary
(one allocation per call), and it passes integer and float arguments through the
same register file, so it "does not properly call functions that have both
integer and float parameters" (its own doc comment).

This proposal adds two symbols that together close both gaps, and sketches a third
that extends them to struct-by-value:

1. **`CallN`** — a low-level, allocation-free call primitive. Fixed-size argument
   arrays passed by pointer (no variadic slice), integers and floats in separate
   register files, and recovery of integer and floating-point results. Unsafe with
   respect to Go pointer lifetimes.
2. **`Args`** — a safe, allocation-free builder over `CallN`. Typed append
   methods, automatic pinning of pointer arguments via an inline `runtime.Pinner`,
   and accessors for floating-point results.
3. **`StructArg[T]`** (appendix, prototype) — by-value struct arguments with the
   platform C ABI, zero-allocation after first use of each type, behind a single
   generic function. Not proposed for merge yet; included to show the design is
   reachable without a libffi-style prepared-interface API.

All three reuse the existing `syscall_SyscallN` funnel and require **no new
assembly**.

## Background: why `SyscallN` allocates and mis-handles floats

Two independent problems:

- **Allocation.** `SyscallN(fn uintptr, args ...uintptr)` builds a variadic
  slice. For a call site in a different module, escape analysis cannot prove the
  slice does not outlive the call, so it heap-allocates — one allocation per call,
  matching the argument-slice size. The alloc-free variadic alternative
  (`//go:nosplit` + `//go:uintptrkeepalive`, as `syscall.SyscallN` uses on
  Windows) is **not reachable from a third-party module**: the compiler gates the
  bare `//go:uintptrkeepalive` pragma to the standard library
  (`noder.go`: `flag == ir.UintptrKeepAlive && !base.Flag.Std`), verified by
  compiling on go1.26 and tip. Dropping `//go:uintptrescapes` without keepalive
  reintroduces a use-after-free. So a fixed-arity path is the only allocation-free
  option third parties actually get. (This is the conclusion of the #445
  discussion.)

- **Floats.** The C ABI passes integer and floating-point arguments in separate
  register files (amd64: general registers vs XMM0–7; arm64: x0–7 vs d0–d7), and
  a float argument does not consume an integer slot or vice versa. `SyscallN`
  takes one flat `...uintptr` and copies the same values into both files, so it
  cannot express "this value goes in a float register, that one in an integer
  register." That is the mixed int/float limitation its doc records.

## Proposal 1: `CallN`

```go
func CallN(fn uintptr, ints *[MaxArgs]uintptr, floats *[8]uintptr, structRet uintptr) (r1, r2, rf1, rf2 uintptr)
```

- **No allocation.** Arguments arrive through fixed-size arrays passed by pointer;
  no variadic slice is formed at the call site, so a correctly written call
  performs no heap allocation.
- **Correct floats.** Integer arguments occupy `ints[0..]`; floating-point
  arguments occupy `floats[0..7]` as bit patterns (`math.Float64bits`), and the
  funnel loads them into the float register file. This calls mixed int/float
  functions correctly.
- **Integer and float results.** Returns the two integer result registers (`r1`,
  `r2`) and the raw bit patterns of the two floating-point result registers (`rf1`,
  `rf2`). Both are already captured by the funnel; recovering them costs nothing.
  A float32 return occupies the low 32 bits of `rf1` on the little-endian
  architectures and the high 32 bits on big-endian s390x; `Args.Float32Result`
  recovers it on either.
- **No errno.** errno is intentionally omitted from v1 (see "errno" below).
- **Bounded by construction.** `*[MaxArgs]uintptr` makes over-length arguments a
  compile error rather than `SyscallN`'s runtime panic. On the architectures
  `CallN` supports (amd64, arm64, loong64, riscv64) there are exactly eight float
  registers; a ninth float argument cannot be expressed, matching the ABI (floats
  never spill to the stack).

### Alignment with `SyscallN`

`CallN` is the low-level peer of `SyscallN` and deliberately matches its
conventions where they are load-bearing, and improves them where `SyscallN`'s are
weaker:

| aspect | `SyscallN` | `CallN` | |
|---|---|---|---|
| first two returns | `r1, r2` | `r1, r2` | same |
| `fn == 0` | panics | panics (same message) | same |
| over-length args | runtime panic | compile error | improved |
| float return | not available | `rf1, rf2` | new capability |
| errno | `err` (last return) | omitted from v1 | see below |
| pointer lifetime | kept alive by `//go:uintptrescapes` | **caller responsibility** | stricter (see below) |

### errno

`SyscallN` returns errno as its third value, but the funnel only writes errno into
`a3` under some OS and cgo configurations — under `GOOS_darwin` in the assembly
funnels, and via the cgo funnel on Linux — while the pure-Go funnels on non-Darwin
minor architectures leave `a3` as the third input argument. Exporting a new API
that documents a cross-platform errno on top of that non-uniform behavior would be
dishonest and was caught misbehaving on the Linux minor-architecture CI job.
Therefore v1 of `CallN`/`Args` omits errno; a caller that needs the C error code
uses `SyscallN`. errno can be added later, once the funnels share a uniform
contract with tests across Darwin, Linux, and the minor architectures.

### Pointer lifetime (the one unsafe aspect)

`CallN` cannot carry `//go:uintptrescapes` — that pragma applies only to variadic
`...uintptr` parameters, not to `*[MaxArgs]uintptr`. So `CallN` does **not** keep
Go pointers alive across the call. Passing the address of a Go object as an
integer argument without keeping it independently reachable is a use-after-free.
The rules in `unsafe.Pointer`, especially point 4, apply without `SyscallN`'s
safety net. This is documented as a caller responsibility; callers pass scalars
and already-pinned buffers, or use `Args` (below), which pins for them.

## Proposal 2: `Args`

`Args` is the safe, allocation-free builder over `CallN`.

```go
var a purego.Args
a.Int(3)
a.Float64(1.5)
a.Ptr(unsafe.Pointer(buf)) // pinned for the duration of the call
r1, _ := a.Call(fn)
// optional result recovery:
d := a.Float64Result()     // double return
```

- **Safe pointers, still zero-alloc.** `Ptr` pins its argument with an inline
  `runtime.Pinner`. Pinning an already-heap object with a stack `Pinner` does not
  allocate, so `Args` is safe **and** allocation-free simultaneously — it closes
  `CallN`'s use-after-free window at no cost to the allocation-free property.
- **Ergonomic results.** `Call` returns the two integer result registers, keeping
  the common signature small; `Float64Result` and `Float32Result` recover a
  floating-point return for the calls that need it. `Float32` and `Float32Result`
  place and recover a single-precision value in the register half the target ABI
  uses (low on little-endian, high on big-endian s390x), so both work on every
  supported architecture.
- **Zero value ready, reusable.** Not safe for concurrent use.

## Allocation and correctness (measured)

Allocations are deterministic and are the durable story; exact ns/op is
machine-dependent.

| call | allocs/op |
|---|---|
| `SyscallN` (1 arg) | 1 (8 B) |
| `CallN` | 0 |
| `Args` (scalar) | 0 |
| `Args` (pinned pointer) | 0 |

Correctness is covered by tests for argument ordering, mixed int/float calls (the
case `SyscallN` cannot express), int stack-spill combined with float arguments,
and float and double returns — verified on every supported architecture, including
big-endian s390x (which places single-precision values in the high half of the
floating-point register), through the CI minor-architecture job.

## Scope: what these do *not* do

`CallN`/`Args` classify each argument by which array or method it is given, not by
inspecting its type. They therefore do **not** pass structs by value: a struct's
fields land in a type-dependent distribution across the integer and float register
files that can only be computed by walking the type. Large structs are passed by
pointer (`Args.Ptr`) and returned via `CallN`'s `structRet` parameter (or
`Args.StructReturn`); small structs and homogeneous floating-point aggregates by
value stay with `RegisterFunc`, which already walks argument types with reflection.

This is a deliberate boundary, not an oversight: `CallN`'s speed comes from *not*
doing a type walk. It positions these primitives beside reflection-based
`RegisterFunc` and beside libffi-style FFIs (e.g. goffi), rather than replacing
them — a zero-reflection fast lane for scalar, float, and pointer calls.

## Staging

1. **`CallN`** targeting `main`.
2. **`Args`** stacked on `CallN` (it cannot compile without it).
3. **`StructArg`** as a follow-up, once the primitive and builder are accepted.

## Appendix: `StructArg[T]` — zero-alloc by-value structs without a big API

By-value struct support looks like it needs a libffi-style prepared-interface API
(`Prepare` a call-interface descriptor, then `Call` it). It does not. The
placement plan for a struct is a pure function of its type, so it can be computed
once per type and cached, and the struct can be taken by typed pointer via a
generic function — no `any`-boxing, no per-call reflection:

```go
func StructArg[T any](a *Args, v T)
```

```go
var a purego.Args
a.Int(n)
purego.StructArg(&a, vec2{X: 1, Y: 2}) // HFA -> two float registers
a.Float64(s)
r1, _ := a.Call(fn)
```

A prototype (simplified classifier covering HFA, small structs, and
large-struct-by-pointer) measures **0 allocs/op** after the first use of each
type, at roughly the cost of a bare `CallN`. The public surface added is exactly
one generic function; the per-type plan and its cache are unexported.

Two alternatives were measured and rejected:

- A non-generic `Args.Struct(v any)` boxes the struct into `any`, which escapes:
  **1 alloc/op**. Convenient but defeats the purpose.
- An exposed `Prepare`/`CallInterface` API (the goffi shape) is also zero-alloc
  but adds several exported types the caller must manage. `StructArg[T]` gets the
  same performance behind one function.

To be merge-ready, `StructArg` needs purego's existing ~600-line `placeRegisters`
classifier wired into the plan computation (which runs once per type, so its cost
is irrelevant). That is a larger change than `CallN`/`Args` and is why it is a
follow-up rather than part of this proposal.
