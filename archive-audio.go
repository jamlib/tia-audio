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
  "github.com/brinkt/archive-audio/details"
  "github.com/brinkt/archive-audio/ffmpeg"
  "github.com/brinkt/archive-audio/utils"
)

// go-args: define app args
type args struct {
  Quality string `help:"mp3 quality: 320, V0"`
  Dir string `help:"directory where files will be saved"`
  Url string `arg:"positional,required" help:"archive.org details url"`
}

// go-args: print app description
func (args) Description() string {
  return "\narchive.org lossless audio downloader, transcoder and tagger"
}

// go-args: print app version
func (args) Version() string {
  return "archive-audio 0.0.3\n"
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
func albumArtwork(imgUrl, outPath string) string {
  downloadPath := download(imgUrl, outPath)
  if len(downloadPath) > 0 {
    f := ffmpeg.Create(path.Join(outPath, downloadPath))

    optImg := path.Join(outPath, "folder.jpg")
    err := f.OptimizeAlbumArt(optImg)
    if err != nil {
      log.Fatal(err)
    }
    return optImg
  }
  return ""
}

// process archive.org details
func process(d details.Details, args args) {
  fmt.Println("Aritst:", d.Artist)
  fmt.Println("Album:", d.Album)

  meta := ffmpeg.Metadata{
    Artist: d.Artist,
    Album: d.Album,
    Date: strings.Replace(d.Date, ".", "-", -1),
  }

  // create directory
  outPath := path.Join(args.Dir, d.Artist + "/" + d.Date[:4] + "/" + d.Album)
  os.MkdirAll(outPath, 0775)

  albumArt := albumArtwork(d.Artwork, outPath)
  if len(albumArt) > 0 {
    fmt.Println("Album Art:", albumArt)
    meta.Artwork = albumArt
  }

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

    f := ffmpeg.Create(inFile)
    err := f.ToMp3(args.Quality, meta, outFile)
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
  p := arg.MustParse(&args)

  // check ffmpeg installed
  if _, e0 := ffmpeg.Which(); e0 != nil {
    p.Fail(e0.Error())
  }

  // check url
  if e1 := details.InvalidUrl(args.Url); e1 != nil {
    p.Fail(e1.Error())
  }

  // check output directory
  dirStat, e2 := os.Stat(args.Dir)
  if e2 != nil || !dirStat.IsDir() {
    // default to working directory
    d, e3 := filepath.Abs("./")
    if e3 != nil {
      log.Fatal(e3)
    }
    args.Dir = d
  }

  // set default mp3 quality
  if args.Quality != "V0" && args.Quality != "320" {
    args.Quality = "320"
  }

  fmt.Printf("\nProcessing URL: %s...\n\n", args.Url)
  d, err := details.ProcessUrl(args.Url)
  if err != nil {
    // debug Details{}
    fmt.Printf("%s\n", err.Error())
    fmt.Printf("\n%#v\n\n", d)
    os.Exit(1)
  }

  process(d, args)
}
