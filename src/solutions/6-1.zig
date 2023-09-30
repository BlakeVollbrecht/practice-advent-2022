const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/6.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var marker_end: u32 = 0;
    var last_four: [4]u8 = undefined;
    var last_four_index: u8 = 0;

    while (true) {
        const char = in_stream.readByte() catch |err| switch (err) {
            error.EndOfStream => break,
            else => |e| return e,
        };

        last_four[last_four_index] = char;

        if (last_four_index == 3) {
            last_four_index = 0;
        } else {
            last_four_index += 1;
        }

        marker_end += 1;

        if (marker_end > 3 and hasUniqueChars(&last_four)) {
            break;
        }
    }

    try stdout.print("{d}\n", .{marker_end});
    try bw.flush();
}

fn hasUniqueChars(sequence: []const u8) bool {
    for (sequence[0 .. sequence.len - 1], 0..(sequence.len - 1)) |c1, i| {
        for (sequence[i + 1 ..]) |c2| {
            if (c1 == c2) {
                return false;
            }
        }
    }

    return true;
}
