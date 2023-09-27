const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    // combined each 3 lines of the original input into one space-delimited line using the
    // paste command: paste -d " " - - - < 3.txt | cat >> 3-triplets.txt
    var file = try std.fs.cwd().openFile("inputs/3-triplets.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var priority_sum: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var packs = std.mem.split(u8, line, " ");
        var group_item_flags: u64 = std.math.maxInt(u64);

        while (packs.next()) |pack| {
            const pack_flags = getPackFlags(pack);
            group_item_flags &= pack_flags;
        }

        const group_priority = priorityFromItemFlag(group_item_flags);
        priority_sum += group_priority;
        // std.debug.print("{d}, {b}\n", .{ group_priority, group_item_flags });
    }

    try stdout.print("{d}\n", .{priority_sum});
    try bw.flush();
}

fn getPackFlags(pack: []const u8) u64 {
    var pack_flags: u64 = 0;

    for (pack) |item| {
        const item_flag = getItemFlag(item);
        pack_flags = pack_flags | item_flag;
    }

    return pack_flags;
}

fn getItemFlag(item: u8) u64 {
    var item_flag: u64 = 1;
    const flag_index = getPriority(item); // priority converts alpha chars to numbers 1-52
    item_flag <<= @as(u6, @intCast(flag_index));
    return item_flag;
}

fn priorityFromItemFlag(item_flag: u64) u8 {
    var flag_shifter = item_flag;
    const end_condition: u64 = 1;
    var priority: u8 = 0;

    while (flag_shifter & end_condition == 0) {
        flag_shifter >>= 1;
        priority += 1;
    }

    return priority;
}

fn getPriority(item: u8) u8 {
    const priority = switch (item) {
        'A'...'Z' => item - 38,
        'a'...'z' => item - 96,
        else => 0,
    };
    return priority;
}
