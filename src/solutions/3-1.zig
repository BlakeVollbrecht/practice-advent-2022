const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/3.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var priority_sum: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        const halfway_index = line.len / 2;
        const compartment1 = line[0..halfway_index];
        const compartment2 = line[halfway_index..];

        outer: for (compartment1) |item1| {
            for (compartment2) |item2| {
                if (item1 == item2) {
                    priority_sum += getPriority(item1);
                    break :outer;
                }
            }
        }
    }

    try stdout.print("{d}\n", .{priority_sum});
    try bw.flush();
}

fn getPriority(item: u8) u8 {
    const priority = switch (item) {
        'A'...'Z' => item - 38,
        'a'...'z' => item - 96,
        else => 0,
    };
    return priority;
}
