const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/10.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var cycle: i32 = 0;
    var x: i32 = 1;
    var signal_strength_sum: i32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.split(u8, line, " ");
        var parts: [2][]const u8 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            parts[i] = value;
            i += 1;
        }
        // std.debug.print("LINE: {s}\n", .{line});

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
            cycle += 1;

            if (@rem(cycle, 20) == 0 and !(@rem(cycle, 40) == 0)) {
                signal_strength_sum += x * cycle;
                // std.debug.print("%20: {d}  {d}  {d}\n", .{ cycle, x, signal_strength_sum });
            }
        }

        if (std.mem.eql(u8, instruction, "addx")) {
            const value = try std.fmt.parseInt(i32, parts[1], 10);
            x += value;
        }
    }

    try stdout.print("{d}\n", .{signal_strength_sum});
    try bw.flush();
}
