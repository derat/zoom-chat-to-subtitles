# zoom-chat-subtitles

This program converts a [Zoom] chat log to a [SubRip] `.srt` subtitle file:

> The SubRip file format is described on the Matroska multimedia container
> format website as "perhaps the most basic of all subtitle formats."

The `.srt` file can be [added as a closed-captioning track in Google
Drive](https://support.google.com/drive/answer/1372218).

[Zoom]: https://zoom.us/
[SubRip]: https://en.wikipedia.org/wiki/SubRip

## Usage

The Zoom chat log should look like the following:

```
10:58:52	 From Someone : Hi there
10:59:03	 From Someone : How are you doing?
...
```

Figure out the time when the video started and pass it and the chat log to the
program:

```
$ go run main.go -start-time 10:57:06 <chat.txt >captions.srt
```

Subtitle data is written to stdout with timestamps relative to the supplied
start time:

```
1
00:01:46,000 --> 00:01:54,000
Someone: Hi there

2
00:01:57,000 --> 00:02:05,000
Someone: How are you doing?

...
```

Several flags are available to configure behavior:

```
Usage: main [flag]...
Reads a Zoom chat log from stdin and writes SubRip subtitles to stdout.
Flags:
  -first-names
        Only include first names in subtitles
  -show-sec int
        Seconds to show subtitles (default 8)
  -start-time string
        Video start time as HH:MM:SS
```
