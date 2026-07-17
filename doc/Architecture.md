## Compiling with the build.zig Toolchain
The presence of build.zig means the project skips traditional, slower make tools. Instead, it relies on the Zig build system to orchestrate its multi-language compilation path.
Zig natively parses your Assembly files (.yasm), creates standard static library objects, and bundles them into an output format that Go's linker can immediately read.
```zig
// build.zig - Concept compilation script for Aura Moby
const std = @import("std");

pub fn build(b: *std.Build) void {
    // 1. Establish cross-compilation target and release optimization levels
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // 2. Instantiate a static C-ABI artifact library container
    const lib = b.addStaticLibrary(.{
        .name = "aurabridge",
        .target = target,
        .optimize = optimize,
    });

    // 3. Inject native C/Zig bridging interfaces
    lib.addCSourceFile(.{
        .file = b.path("native/bridge.zig"),
        .flags = &[_][]const u8{ "-Wall", "-Wextra" },
    });

    // 4. Compile and link the 21.7% pure assembly code layout
    lib.addAssemblyFile(b.path("native/windows1252.yasm"));

    // 5. Enforce standard system C library linking paths
    lib.linkLibC();

    // 6. Direct the compiler to output the build object to ../build/native
    const artifact = b.addInstallArtifact(lib, .{});
    b.getInstallStep().dependOn(&artifact.step);
}
```
------------------------------
*  Verifying Build Directives in CLAUDE.md
In modern AI-assisted engineering repositories, CLAUDE.md acts as the primary instruction manual for coding agents and developer environments. It establishes the commands required to safely run cross-language compilation passes without breaking the core runtime symbols.
A typical CLAUDE.md file for this exact repository stack dictates the sequential order of operations needed to build the container engine:

> Aura Moby Development Guidelines## Build Commands- Compile native C-ABI/Assembly layers: `zig build -Doptimize=ReleaseFast`

- Run local Go daemon tests: `go test ./engine/...`

- Full project compilation: `zig build && go build -o aura-moby ./cmd/dev`
Code Conventions- Memory Pinning: Always use `runtime.Pinner` or `unsafe.Pointer` checks when passing Go byte slices down to Assembly layers to stop Garbage Collector relocations.- Register Conventions: All assembly files must follow System V AMD64 ABI (Linux) or Microsoft x64 ABI (Windows) register configurations.


------------------------------
> # Navigating the doc/ Architecture Guide

To understand how these pieces fit together, you can inspect the localized markdown files within the cloned workspace. The structural organization separates high-level coordination logic from down-to-the-metal hardware instruction models:

> # doc/ai.md (The Intelligent Orchestrator)

This manual maps out how the ai/scheduler/ uses telemetry metrics to predict container thread blockages. It documents variables like tensor_load_ratio and details the telemetry loops that decide when to trigger hardware-level shifts.

> # doc/native.md (The Register Matrix)

This file defines how data transforms across language boundaries. It documents exactly which CPU registers pass string vectors down to the Assembly layout (windows1252.yasm), providing explicit safety warnings for developers adding new native instructions.

------------------------------
