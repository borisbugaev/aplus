const std = @import("std");

pub fn quiz(line: []u8) void {}

pub fn main() void {
    var gpa = try std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    std.debug.print("ahoy\n", .{});
    var zout = std.io.getStdOut().writer();
    try zout.print("Read\n", .{});
    const args = try std.process.args();
    const input_file = try std.fs.cwd().openFile(args[1], .{});
    defer input_file.close();
    var bffr_reader = try std.io.bufferedReader(input_file.reader());
    var input_stream = try bffr_reader.reader();
    var bffr: [1024]u8 = undefined;
    while (try input_stream.readUntilDelimiterOrEof(&bffr, '\n')) |line| {
        try quiz(line);
    }
}
