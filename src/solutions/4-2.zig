const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/4.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var sum: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.splitAny(u8, line, ",-");
        var numbers: [4]u32 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            numbers[i] = try std.fmt.parseInt(u32, value, 10);
            i += 1;
        }

        // std.debug.print("{s}\n", .{line});

        if (numbers[1] < numbers[2] or numbers[3] < numbers[0]) {
            // std.debug.print("no overlap\n", .{});
        } else {
            sum += 1;
        }
    }

    try stdout.print("{d}\n", .{sum});
    try bw.flush();
}
