// archive.org details page
// ie, https://archive.org/details/*

package main

import (
  "log"
  "regexp"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type Track struct {
  Num int
  Title string `json:"title"`
  Source string `json:"orig"`
}

type Details struct {
  Url string
  Body string
  Artist string
  Artwork string
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
    log.Fatalf("http archive.org error")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  d := Details{
    Url: url,
    Body: fixWhitespace(string(body)),
  }

  // parse all from details.Body
  d.parseArtist()
  d.parseArtwork()
  d.parseDate()
  d.parseVenue()
  d.parseLocation()
  d.parseTracks()

  return d
}

// parse 'artist' from HTML body
func (d *Details) parseArtist() {
  d.Artist = removeHtml(regexpBetween(`<div class="key-val-big"> by `, `</div>`, d.Body))
}

// parse 'artwork url' from HTML body
func (d *Details) parseArtwork() {
  s := regexpBetween(`<div id="theatre-controls">`, `<div id="cher-modal"`, d.Body)
  s = regexpBetween(`<img src="`, `"`, s)

  // replace after .ext? including ?
  // ie, image.jpg?other-data => image.jpg
  d.Artwork = regexp.MustCompile(`\?.+$`).ReplaceAllString(s, "")
}

// parse 'date' from HTML body
func (d *Details) parseDate() {
  s := removeHtml(regexpBetween(
    `<div class="key-val-big"> Publication date `, `</a>`, d.Body))

  // use periods instead of dashes, ie 2018.01.01
  d.Date = strings.Replace(s, "-", ".", -1)
}

// parse 'venue' & 'location' from HTML body
func (d *Details) parseVenueAndLocation(id string) string {
  s := regexpBetween(`href="/search.php?query=`+id, `/a>`, d.Body)
  return regexpBetween(`>`, `<`, s)
}

func (d *Details) parseVenue() {
  d.Venue = d.parseVenueAndLocation("venue")
}

func (d *Details) parseLocation() {
  d.Location = d.parseVenueAndLocation("coverage")
}

// parse 'tracks' from HTML body JSON
func (d *Details) parseTracks() {
  s := regexpBetween(`Play('jw6', `, `, {"start"`, d.Body)
  if len(s) == 0 {
    return
  }

  tracks := []Track{}
  json.Unmarshal([]byte(s), &tracks)

  for i := range tracks {
    // remove preceding track number
    tracks[i].Title = regexp.MustCompile(`^\d+\.\s+`).
      ReplaceAllString(tracks[i].Title, "")

    // only allow certain chars
    reg2 := `([A-Z]|[a-z]|[0-9]|[',./!?&> ()])+`
    tracks[i].Title = regexp.MustCompile(reg2).FindString(tracks[i].Title)

    // trim remaining whitespace
    tracks[i].Title = fixWhitespace(tracks[i].Title)
  }

  d.Tracks = tracks
}
