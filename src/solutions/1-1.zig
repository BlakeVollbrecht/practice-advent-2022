const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/1.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var local_sum: i32 = 0;
    var max: i32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (std.mem.eql(u8, line, "")) {
            if (local_sum > max) {
                max = local_sum;
            }
            local_sum = 0;
            // std.debug.print("max: {d}\n", .{max});
        } else {
            const num = try std.fmt.parseInt(i32, line, 10);
            local_sum += num;
        }
    }

    try stdout.print("{d}\n", .{max});
    try bw.flush();
}

// test "simple test" {
//     var list = std.ArrayList(i32).init(std.testing.allocator);
//     defer list.deinit(); // try commenting this out and see if zig detects the memory leak!
//     try list.append(42);
//     try std.testing.expectEqual(@as(i32, 42), list.pop());
// }
