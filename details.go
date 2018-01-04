// archive.org details page
// ie, https://archive.org/details/*

package main

import (
  "fmt"
  "log"
  "regexp"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type Track struct {
  Num string
  Title string `json:"title"`
  Source string `json:"orig"`
}

type DetailsResponse struct {
  Body string
}

type Details struct {
  Url string
  Artwork string
  Artist string
  Album string
  Date string
  Venue string
  Location string
  Tracks []Track
}

// validates an archive.org/details url
func validUrl(url string) (bool, string) {
  valid, _ := regexp.MatchString("^htt(p|ps)://archive.org/details/.+", url)
  if valid {
    return true, ""
  }
  return false, "URL must be a valid archive.org details url\n" + 
    "ie, https://archive.org/details/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16\n"
}

// make url request & parse data from archive.org details response
func processUrl(url string) Details {
  // http request
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  body, e := ioutil.ReadAll(resp.Body)
  if e != nil {
    log.Fatal(e)
  }

  dResp := DetailsResponse{Body: fixWhitespace(string(body))}

  d := Details{
    Url: url,
    Artwork: dResp.parseArtwork(),
    Artist: dResp.parseArtist(),
    Date: dResp.parseDate(),
    Venue: dResp.parseVenue(),
    Location: dResp.parseLocation(),
    Tracks: dResp.parseTracks(),
  }
  d.Album = fmt.Sprintf("%s %s, %s", d.Date, d.Venue, d.Location)

  for i := range d.Tracks {
    d.Tracks[i].Num = fmt.Sprintf("%02d", i + 1)

    if len(d.Tracks[i].Title) == 0 || len(d.Tracks[i].Source) == 0 {
      log.Fatalf("Error: blank track metadata")
    }
  }

  return d
}

// parse 'artist' from HTML body
func (d *DetailsResponse) parseArtist() string {
  return removeHtml(regexpBetween(`<div class="key-val-big"> by `, `</div>`, d.Body))
}

// parse 'artwork url' from HTML body
func (d *DetailsResponse) parseArtwork() string {
  s := regexpBetween(`<div id="theatre-controls">`, `<div id="cher-modal"`, d.Body)
  s = regexpBetween(`<img src="`, `"`, s)

  // replace after .ext
  // ie, image.jpg?other-data => image.jpg
  return regexp.MustCompile(`\?.+$`).ReplaceAllString(s, "")
}

// parse 'date' from HTML body
func (d *DetailsResponse) parseDate() string {
  s := removeHtml(regexpBetween(
    `<div class="key-val-big"> Publication date `, `</a>`, d.Body))

  // use periods instead of dashes, ie 2018.01.01
  return strings.Replace(s, "-", ".", -1)
}

// parse 'venue' & 'location' from HTML body
func (d *DetailsResponse) parseVenueAndLocation(id string) string {
  s := regexpBetween(`href="/search.php?query=`+id, `/a>`, d.Body)
  return regexpBetween(`>`, `<`, s)
}

func (d *DetailsResponse) parseVenue() string {
  return d.parseVenueAndLocation("venue")
}

func (d *DetailsResponse) parseLocation() string {
  return d.parseVenueAndLocation("coverage")
}

// parse 'tracks' from HTML body JSON
func (d *DetailsResponse) parseTracks() []Track {
  s := regexpBetween(`Play('jw6', `, `, {"start"`, d.Body)
  tracks := []Track{}

  if len(s) > 0 {
    json.Unmarshal([]byte(s), &tracks)
  }

  for i := range tracks {
    // remove preceding track number
    tracks[i].Title = regexp.MustCompile(`^\d+\.\s+`).
      ReplaceAllString(tracks[i].Title, "")

    // only allow certain chars
    reg2 := `([A-Z]|[a-z]|[0-9]|[',./!?&> ()])+`
    tracks[i].Title = regexp.MustCompile(reg2).FindString(tracks[i].Title)

    // fix remaining whitespace
    tracks[i].Title = fixWhitespace(tracks[i].Title)
  }

  return tracks
}
