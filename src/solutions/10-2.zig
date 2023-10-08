const std = @import("std");

pub fn solve() !void {
    var file = try std.fs.cwd().openFile("inputs/10.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var cycle: i32 = 0;
    var x: i32 = 1;
    var pixels: [240]bool = undefined;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.split(u8, line, " ");
        var parts: [2][]const u8 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            parts[i] = value;
            i += 1;
        }

        const instruction = parts[0];
        var instruction_cycles: u32 = 0;

        if (std.mem.eql(u8, instruction, "noop")) {
            instruction_cycles = 1;
        } else if (std.mem.eql(u8, instruction, "addx")) {
            instruction_cycles = 2;
        } else {
            std.debug.print("unknown instruction: {s}\n", .{line});
        }

        for (0..instruction_cycles) |_| {
            const current_pixel = @rem(cycle, 240);
            const current_column = @rem(cycle, 40);

            if (x - current_column >= -1 and x - current_column <= 1) {
                pixels[@as(u32, @intCast(current_pixel))] = true;
            } else {
                pixels[@as(u32, @intCast(current_pixel))] = false;
            }

            // std.debug.print("%20: {d}  {d}  {d}\n", .{ cycle, x, current_column });
            cycle += 1;
        }

        if (std.mem.eql(u8, instruction, "addx")) {
            const value = try std.fmt.parseInt(i32, parts[1], 10);
            x += value;
        }
    }

    try printPixels(pixels);
}

fn printPixels(pixels: [240]bool) !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    for (pixels, 0..) |pixel, i| {
        if (i % 40 == 0) {
            try stdout.print("\n", .{});
        }

        if (pixel) {
            try stdout.print("#", .{});
        } else {
            try stdout.print(".", .{});
        }

        try bw.flush();
    }
}
