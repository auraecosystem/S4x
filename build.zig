const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Example step to compile a Go binary utilizing Zig as the high-performance cross-compiler link
    const go_build = b.addSystemCommand(&.{ "go", "build" });
    
    // Define environment arguments to inject Zig's compilation architecture into Go 
    go_build.addArgs(&.{ "-ldflags", "-w -s", "-o", "server" });

    // Expose this runner task via target invocation line: `zig build run-go`
    const run_step = b.step("run-go", "Compile application utilizing standard execution frameworks");
    run_step.dependency(&go_build.step);
}
