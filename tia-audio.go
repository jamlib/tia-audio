package main

import (
  "os"
  "fmt"
  "io"
  "log"
  "strings"
  "net/http"
  "path"
  "path/filepath"

  "github.com/alexflint/go-arg"
  "github.com/JamTools/goff/ffmpeg"
  "github.com/JamTools/tia-audio/details"
  "github.com/JamTools/tia-audio/utils"
)

// go-args: define app args
type args struct {
  Quality string `help:"mp3 quality: 320, V0"`
  Dir string `help:"directory where files will be saved"`
  Url string `arg:"positional,required" help:"archive.org details url"`
}

// go-args: print app description
func (args) Description() string {
  return "\nThe Internet Archive (archive.org) lossless audio downloader, " +
    "transcoder and tagger"
}

// go-args: print app version
func (args) Version() string {
  return "tia-audio 0.0.4\n"
}

// download file
func download(fileUrl, outPath string) string {
  if len(fileUrl) == 0 {
    return ""
  }

  a := strings.Split(fileUrl, "/")
  filePath := a[len(a)-1]

  fmt.Printf("\nDownloading: %s...\n", filePath)

  res, e := http.Get(fileUrl)
  if e != nil {
    log.Fatal(e)
  }
  defer res.Body.Close()

  file, err := os.Create(path.Join(outPath, filePath))
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  _, err = io.Copy(file, res.Body)
  if err != nil {
    log.Fatal(err)
  }

  return filePath
}

// download & optimize album artwork
func albumArtwork(imgUrl, outPath string, meta *ffmpeg.Metadata) {
  fileName := download(imgUrl, outPath)
  if len(fileName) > 0 {
    outFile := path.Join(outPath, "folder.jpg")
    _, err := ffmpeg.OptimizeAlbumArt(path.Join(outPath, fileName), outFile)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println("Album Art:", outFile)
    meta.Artwork = outFile
  }
}

// process archive.org details
func process(d *details.Details, args *args) {
  meta := ffmpeg.Metadata{
    Artist: d.Artist,
    Date: d.Date,
    // use periods instead of dashes, ie 2018.01.01
    Album: fmt.Sprintf("%s %s, %s",
      strings.Replace(d.Date, "-", ".", -1), d.Venue, d.Location),
  }

  fmt.Println("Aritst:", meta.Artist)
  fmt.Println("Album:", meta.Album)

  // create directory
  outPath := path.Join(args.Dir, d.Artist + "/" + d.Date[:4] + "/" + meta.Album)
  os.MkdirAll(outPath, 0775)

  albumArtwork(d.Artwork, outPath, &meta)

  for i := range d.Tracks {
    meta.Track = d.Tracks[i].Num
    meta.Title = d.Tracks[i].Title

    fmt.Println("\nTrack:", meta.Track)
    fmt.Println("Title:", meta.Title)

    downloadUrl := strings.Replace(d.Url, "/details/", "/download/", -1)
    download(downloadUrl + "/" + d.Tracks[i].Source, outPath)

    fmt.Printf("Converting '%s' to '%s' mp3...\n", d.Tracks[i].Source, args.Quality)

    inFile := path.Join(outPath, d.Tracks[i].Source)
    outFile:= path.Join(outPath, utils.SafeFilename(
      meta.Track + " - " + meta.Title + ".mp3"))

    _, err := ffmpeg.ToMp3(inFile, args.Quality, meta, outFile)
    if err != nil {
      log.Fatal(err)
    }

    os.Remove(inFile)
  }

  fmt.Println("\nProcess Completed!\n")
}

// package main entry
func main() {
  var args args
  arg.MustParse(&args)

  // check url
  if err := details.InvalidUrl(args.Url); err != nil {
    log.Fatal(err)
  }

  // check ffmpeg installed
  if _, err := ffmpeg.Which(); err != nil {
    log.Fatal(err)
  }

  // check output directory
  dirStat, err := os.Stat(args.Dir)
  if err != nil || !dirStat.IsDir() {
    // default to working directory
    dir, err := filepath.Abs("./")
    if err != nil {
      log.Fatal(err)
    }
    args.Dir = dir
  }

  // set default mp3 quality
  if args.Quality != "V0" && args.Quality != "320" {
    args.Quality = "320"
  }

  fmt.Printf("\nProcessing URL: %s...\n\n", args.Url)
  d, err := details.ProcessUrl(args.Url)
  if err != nil {
    // debug Details{}
    log.Fatalf("%s\n\n%#v\n\n", err.Error(), d)
  }

  process(d, &args)
}
