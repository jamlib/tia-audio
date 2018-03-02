// archive.org details page
// ie, https://archive.org/details/*

package details

import (
  "fmt"
  "errors"
  "regexp"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "github.com/JamTools/tia-audio/utils"
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
  Date string
  Venue string
  Location string
  Tracks []Track
}

// validates an archive.org/details url
func InvalidUrl(url string) error {
  valid, _ := regexp.MatchString("^htt(p|ps)://archive.org/details/.+", url)
  if !valid {
    return errors.New("URL must be a valid archive.org details url\n" +
      "ie, https://archive.org/details/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16\n")
  }
  return nil
}

// make url request & parse data from archive.org details response
func ProcessUrl(url string) (*Details, error) {
  // http request
  resp, err := http.Get(url)
  if err != nil {
    return &Details{}, err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return &Details{}, err
  }

  dResp := &DetailsResponse{Body: utils.FixWhitespace(string(body))}

  d := &Details{
    Url: url,
    Artwork: dResp.parseArtwork(),
    Artist: dResp.parseArtist(),
    Date: dResp.parseDate(),
    Venue: dResp.parseVenue(),
    Location: dResp.parseLocation(),
    Tracks: dResp.parseTracks(),
  }

  if invalid := d.validate(); invalid != nil {
    return d, invalid
  }
  return d, nil
}

// validate Details{}
func (d *Details) validate() error {
  var strErr string

  // include error if metadata incomplete
  if len(d.Artist) == 0 || len(d.Date) == 0 ||
    len(d.Venue) == 0 || len(d.Location) == 0 {
    strErr = "Error in parse of metadata {Artist, Date, Venue, Location}. "
  }

  // include error if no tracks found
  if len(d.Tracks) == 0 {
    strErr += "Error in parse of tracks json data. "
  }

  // include error if tracks json data incomplete
  for i := range d.Tracks {
    if len(d.Tracks[i].Title) == 0 || len(d.Tracks[i].Source) == 0 {
      strErr += "Error in tracks json data {Title, Source}. "
    }
  }

  if len(strErr) > 0 {
    return errors.New(strErr)
  }
  return nil
}

// parse 'artist' from HTML body
func (d *DetailsResponse) parseArtist() string {
  return utils.RemoveHtml(utils.RegexpBetween(
    `<div class="key-val-big"> by `, `</div>`, d.Body))
}

// parse 'artwork url' from HTML body
func (d *DetailsResponse) parseArtwork() string {
  s := utils.RegexpBetween(`<div id="theatre-controls">`, `<div id="cher-modal"`, d.Body)
  s = utils.RegexpBetween(`<img src="`, `"`, s)

  // replace after .ext
  // ie, image.jpg?other-data => image.jpg
  return regexp.MustCompile(`\?.+$`).ReplaceAllString(s, "")
}

// parse 'date' from HTML body
func (d *DetailsResponse) parseDate() string {
  return utils.RemoveHtml(utils.RegexpBetween(
    `<div class="key-val-big"> Publication date `, `</a>`, d.Body))
}

// parse 'venue' & 'location' from HTML body
func (d *DetailsResponse) parseVenueAndLocation(id string) string {
  s := utils.RegexpBetween(`href="/search.php?query=`+id, `/a>`, d.Body)
  return utils.RegexpBetween(`>`, `<`, s)
}

func (d *DetailsResponse) parseVenue() string {
  return d.parseVenueAndLocation("venue")
}

func (d *DetailsResponse) parseLocation() string {
  return d.parseVenueAndLocation("coverage")
}

// parse 'tracks' from HTML body JSON
func (d *DetailsResponse) parseTracks() []Track {
  s := utils.RegexpBetween(`Play('jw6', `, `, {"start"`, d.Body)
  tracks := []Track{}

  if len(s) > 0 {
    json.Unmarshal([]byte(s), &tracks)
  }

  for i := range tracks {
    // set track num
    tracks[i].Num = fmt.Sprintf("%02d", i + 1)

    // remove preceding track number
    tracks[i].Title = regexp.MustCompile(`^\d+\.\s+`).
      ReplaceAllString(tracks[i].Title, "")

    // only allow certain chars
    reg2 := `([A-Z]|[a-z]|[0-9]|[',./!?&> ()])+`
    tracks[i].Title = regexp.MustCompile(reg2).FindString(tracks[i].Title)

    // fix remaining whitespace
    tracks[i].Title = utils.FixWhitespace(tracks[i].Title)
  }

  return tracks
}
