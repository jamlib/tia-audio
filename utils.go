// general utility functions

package main

import (
  "os"
  "regexp"
  "strings"
)

// check if path is valid directory
func isDirectory(p string) bool {
  fileInfo, err := os.Stat(p)
  if err != nil {
    return false
  }
  return fileInfo.IsDir()
}

// replaces \ & / from proposed file name
func safeFilename(f string) string {
  return regexp.MustCompile(`[\/\\]+`).ReplaceAllString(f, "_")
}

// replace whitespaces with one space
func fixWhitespace(s string) string {
  return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(s, " "))
}

/*
  regexpBetween makes use of 'between' syntax:
  ie, (?s)BEFORE.*?AFTER

  trims BEFORE & AFTER from result
*/

func regexpBetween(before, after, within string) string {
  r := `(?s)` + regexp.QuoteMeta(before) + `.*?` +
    regexp.QuoteMeta(after)

  s := regexp.MustCompile(r).FindString(within)

  // trim before & after
  if len(s) >= len(before)+len(after) {
    s = s[len(before):len(s)-len(after)]
  }

  return s
}

// remove HTML aka stuff betwee < & > from string
func removeHtml(body string) string {
  return strings.TrimSpace(
    regexp.MustCompile(`<[^<>]+>`).ReplaceAllString(body, ""))
}
