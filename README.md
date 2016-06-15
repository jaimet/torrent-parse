# torrent-parse
A simple program for dumping metadata from a BitTorrent file.

This tool is intended for quickly inspecting BitTorrent files without
opening them in a client program.

To avoid terminal corruption, `torrent-parse` currently only prints
fields that I've explicitly added support for. (Of course, I could
simply escape non-printable characters in unknown fields, but I also
like being able to list fields under easily readable names.) For all
unsupported fields, a warning is printed to standard error.

`bencode.go` contains a general-purpose bencode parser which should be
easily usable in other projects. If you try this and have any issues or
bug reports, please share your experience.
