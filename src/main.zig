const std = @import("std");

const solution = @import("solutions/12-1.zig");

pub fn main() !void {
    const start = std.time.nanoTimestamp();

    try solution.solve();

    const end = std.time.nanoTimestamp();
    const time_ms = (@as(f64, @floatFromInt(end)) - @as(f64, @floatFromInt(start))) / 1000000;
    std.debug.print("Time: {d}ms", .{time_ms});
}
