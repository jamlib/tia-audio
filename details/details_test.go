package details

import (
  "testing"
)

func TestInvalidUrl(t *testing.T) {
  urls := [2]string{
    "https://archive.org/details/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16",
    "http://archive.org/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16",
  }

  err := InvalidUrl(urls[0])
  if err != nil {
    t.Error("Expected: valid url")
  }

  err = InvalidUrl(urls[1])
  if err == nil {
    t.Error("Expected: invalid url")
  }
}

func TestUrlRequest(t *testing.T) {
  var err error

  _, err = urlRequest("invalid-url")
  if err == nil {
    t.Error("Expected: invalid url")
  }

  _, err = urlRequest("http://archive.org")
  if err != nil {
    t.Error("Expected: valid url")
  }
}

func TestDetailsValidate(t *testing.T) {
  var err error
  var d *Details

  // test invalid metadata
  d = &Details{}
  err = d.validate()
  if err == nil {
    t.Error("Expected: Error in parse of metadata")
  }

  // test missing tracks
  d = &Details{ Artist: "Artist", Date: "2018-01-01", Venue: "Venue", Location: "Location" }
  err = d.validate()
  if err == nil {
    t.Error("Expected: Error in parse of tracks json data")
  }

  // test missing track data
  d.Tracks = []Track{ { Num: "1" } }
  err = d.validate()
  if err == nil {
    t.Error("Expected: Error in tracks json data")
  }
}

func TestProcessUrl(t *testing.T) {
  var err error

  urls := [2]string{
    "https://archive.org/details/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16",
    "http://archive.org/jrad2017-01-12.cmc621.cmc64.sbd.matrix.flac16",
  }

  _, err = ProcessUrl(urls[0])
  if err != nil {
    t.Error("Expected: valid details")
  }

  _, err = ProcessUrl(urls[1])
  if err == nil {
    t.Error("Expected: invalid details url")
  }
}
