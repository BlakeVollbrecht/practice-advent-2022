const std = @import("std");

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/2.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var my_score: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        const opponent_move = line[0];
        const outcome = line[2];

        var my_move = opponent_move; // default to outcome = 'Y'

        if (outcome == 'X') { // must lose
            my_move = switch (opponent_move) {
                'A' => 'C',
                'B' => 'A',
                'C' => 'B',
                else => 'E',
            };
        } else if (outcome == 'Z') { // must win
            my_move = switch (opponent_move) {
                'A' => 'B',
                'B' => 'C',
                'C' => 'A',
                else => 'E',
            };
        }

        var points = my_move - 64; // convert utf-8 A,B,C to int 1,2,3

        if (my_move == opponent_move) {
            // add 3 points for a draw
            points += 3;
        } else if ((my_move == 'A' and opponent_move == 'C') or
            (my_move == 'B' and opponent_move == 'A') or
            (my_move == 'C' and opponent_move == 'B'))
        {
            // add 6 points for a win
            points += 6;
        }

        my_score += points;

        // std.debug.print("{c} {c} {d}\n", .{ my_move, opponent_move, points });
    }

    try stdout.print("{d}\n", .{my_score});
    try bw.flush();
}
