package ffmpeg

import (
  "fmt"
  "log"
  "bytes"
  "errors"
  "os/exec"
)

type ffmpeg struct {
  *exec.Cmd
}

type Metadata struct {
  Artist string
  Album string
  Title string
  Track string
  Date string
  Artwork string
}

// check that ffmpeg is installed on system
func Which() string {
  path, err := exec.LookPath("ffmpeg")
  if err != nil {
    log.Fatalf("Error: ffmpeg not found on system")
  }

  return path
}

// new ffmpeg wrapper where args can be added
func Create(input string) *ffmpeg {
  return &ffmpeg{exec.Command(Which(), "-i", input)}
}

// add additional arguments
func (f *ffmpeg) setArgs(args ...string) {
  f.Args = append(f.Args, args...)
}

// append output as final arg & run ffmpeg
func (f *ffmpeg) run(output string) error {
  // include ffmpeg debug message
  var out bytes.Buffer
  var stderr bytes.Buffer
  f.Stdout = &out
  f.Stderr = &stderr

  f.setArgs(output)

  err := f.Run()
  if err != nil {
    return errors.New(fmt.Sprint(err) + ": " + stderr.String())
  }
  return err
}

// optimize image as embedded album art
func (f *ffmpeg) OptimizeAlbumArt(output string) error {
  f.setArgs("-y", "-qscale:v", "2", "-vf", "scale=500:-1")
  return f.run(output)
}

// convert lossless to mp3
func (f *ffmpeg) ToMp3(quality string, meta Metadata, output string) error {
  if len(meta.Artwork) > 0 {
    f.setArgs("-i", meta.Artwork)
  }

  f.setArgs("-y", "-map", "0:a", "-codec:a", "libmp3lame")

  if quality == "320" {
    f.setArgs("-b:a", "320k")
  } else {
    f.setArgs("-qscale:a", "0")
  }

  f.setArgs("-metadata", "artist=" + meta.Artist)
  f.setArgs("-metadata", "album=" + meta.Album)
  f.setArgs("-metadata", "title=" + meta.Title)
  f.setArgs("-metadata", "track=" + meta.Track)
  f.setArgs("-metadata", "date=" + meta.Date)

  if len(meta.Artwork) > 0 {
    f.setArgs("-map", "1:v", "-c:v", "copy", "-metadata:s:v", "title=Album cover",
      "-metadata:s:v", "comment=Cover (Front)")
  }

  f.setArgs("-id3v2_version", "4")
  return f.run(output)
}
