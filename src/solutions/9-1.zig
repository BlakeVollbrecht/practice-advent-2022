const std = @import("std");
const ArrayList = std.ArrayList;
const Vector = std.meta.Vector;
const allocator = std.heap.page_allocator;

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/9.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;

    var tail_visited: ArrayList([]const u8) = ArrayList([]const u8).init(allocator);
    defer tail_visited.deinit();
    try tail_visited.append("0,0");
    var head: @Vector(2, i32) = .{ 0, 0 };
    // var tail: @Vector(2, i32) = .{0, 0};

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.split(u8, line, " ");
        var parts: [2][]const u8 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            parts[i] = value;
            i += 1;
        }
        // std.debug.print("LINE: {s}\n", .{line});

        const direction = parts[0];
        const distance = try std.fmt.parseInt(u32, parts[1], 10);

        for (0..distance) |_| {
            const movement: @Vector(2, i32) = switch (direction[0]) {
                'L' => .{ -1, 0 },
                'R' => .{ 1, 0 },
                'U' => .{ 0, 1 },
                'D' => .{ 0, -1 },
                else => unreachable,
            };
            head = head + movement;
        }
    }

    try stdout.print("{any}\n", .{head});
    try bw.flush();
}
