// bridge.zig - Compiled via build.zig into the runtime layer
const std = @import("std");

// External link to the raw Assembly routine
extern fn aura_transform_buffer_abi(input: [*]const u8, output: [*]u8, len: u64) callconv(.C) void;

// Exporting a guaranteed C-ABI function name to be consumed by Go's runtime
export fn aura_process_payload(in_ptr: [*]const u8, out_ptr: [*]u8, length: u64) callconv(.C) i32 {
    if (length == 0) return -1;
    
    // Invoke the optimized assembly pathway securely 
    aura_transform_buffer_abi(in_ptr, out_ptr, length);
    
    return 0; // Success code
}
