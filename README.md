# archive-audio

Download, transcode to mp3, embed album artwork, & tag lossless audio from `archive.org`.

Transcodes to `320 kbps` or `V0` mp3 using `ffmpeg` & `libmp3lame`.

## Usage

```
Usage: archive-audio [--quality QUALITY] [--dir DIR] URL

Positional arguments:
  URL                    archive.org details url

Options:
  --quality QUALITY      mp3 quality: 320, V0
  --dir DIR              directory where files will be saved
  --help, -h             display this help and exit
  --version              display version and exit
```

## License

This code is available open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).