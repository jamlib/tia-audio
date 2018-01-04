package main

import (
  "os"
  "fmt"
  "io"
  "log"
  "regexp"
  "strings"
  "net/http"
  "path"
  "path/filepath"

  "github.com/alexflint/go-arg"
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
  return "archive-audio 0.0.1\n"
}

// download file
func Download(fileUrl, outPath string) string {
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
func AlbumArtwork(imgUrl, outPath string) string {
  downloadPath := Download(imgUrl, outPath)
  if len(downloadPath) > 0 {
    f := newFfmpeg(path.Join(outPath, downloadPath))

    optImg := path.Join(outPath, "folder.jpg")
    err := f.optimizeAlbumArt(optImg)
    if err != nil {
      log.Fatal(err)
    }
    return optImg
  }
  return ""
}

// main process
func Handle(d Details, args args) {
  album := fmt.Sprintf("%s %s, %s", d.Date, d.Venue, d.Location)
  fmt.Println("Aritst:", d.Artist)
  fmt.Println("Album:", album)

  meta := Metadata{
    Artist: d.Artist,
    Album: album,
    Date: strings.Replace(d.Date, ".", "-", -1),
  }

  // create directory
  savePath := d.Artist + "/" + d.Date[:4] + "/" + album
  args.Dir = path.Join(args.Dir, savePath)
  os.MkdirAll(args.Dir, 0775)

  albumArt := AlbumArtwork(d.Artwork, args.Dir)
  if len(albumArt) > 0 {
    fmt.Println("Album Art:", albumArt)
    meta.Artwork = albumArt
  }

  for i := range d.Tracks {
    meta.Track = fmt.Sprintf("%02d", i + 1)
    meta.Title = d.Tracks[i].Title

    // title: replaces / or \ with _
    titlePath := regexp.MustCompile(`[\/\\]`).ReplaceAllString(meta.Title, "_")

    if len(meta.Title) == 0 || len(d.Tracks[i].Source) == 0 {
      log.Fatalf("Error: blank track metadata")
    }

    fmt.Println("\nTrack:", meta.Track)
    fmt.Println("Title:", meta.Title)

    downloadUrl := strings.Replace(d.Url, "/details/", "/download/", -1)
    Download(downloadUrl + "/" + d.Tracks[i].Source, args.Dir)

    fmt.Printf("Converting '%s' to '%s' mp3...\n", d.Tracks[i].Source, args.Quality)

    inPath := path.Join(args.Dir, d.Tracks[i].Source)
    outPath := path.Join(args.Dir, meta.Track + " - " + titlePath + ".mp3")

    f := newFfmpeg(inPath)
    err := f.toMp3(args.Quality, meta, outPath)
    if err != nil {
      log.Fatal(err)
    }

    os.Remove(inPath)
  }

  fmt.Println("\nProcess Completed!\n")
}

// package main entry
func main() {
  var args args
  p := arg.MustParse(&args)

  // check ffmpeg installed
  whichFfmpeg()

  // check url
  valid, err := validUrl(args.Url)
  if !valid {
    p.Fail(err)
  }

  // check output directory
  if isDirectory(args.Dir) == false {
    // default to working directory
    d, err := filepath.Abs("./")
    if err != nil {
      log.Fatal(err)
    }
    args.Dir = d
  }

  // set default mp3 quality
  if args.Quality != "V0" && args.Quality != "320" {
    args.Quality = "320"
  }

  fmt.Printf("\nProcessing URL: %s...\n\n", args.Url)
  d := processUrl(args.Url)

  // ensure metadata extraction was successful
  if len(d.Artist) == 0 || len(d.Date) == 0 ||
    len(d.Venue) == 0 || len(d.Location) == 0 {
    log.Fatalf("Error: unable to parse metadata from HTML")
  }

  Handle(d, args)
}
