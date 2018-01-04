// general utility functions

package utils

import (
  "regexp"
  "strings"
)

// replaces \ & / from proposed file name
func SafeFilename(f string) string {
  return regexp.MustCompile(`[\/\\]+`).ReplaceAllString(f, "_")
}

// replace whitespaces with one space
func FixWhitespace(s string) string {
  return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(s, " "))
}

/*
  regexpBetween makes use of 'between' syntax:
  ie, (?s)BEFORE.*?AFTER

  trims BEFORE & AFTER from result
*/

func RegexpBetween(before, after, within string) string {
  r := `(?s)` + regexp.QuoteMeta(before) + `.*?` +
    regexp.QuoteMeta(after)

  s := regexp.MustCompile(r).FindString(within)

  // trim before & after
  s = strings.TrimPrefix(s, before)
  s = strings.TrimSuffix(s, after)

  return s
}

// remove HTML aka stuff betwee < & > from string
func RemoveHtml(body string) string {
  return strings.TrimSpace(
    regexp.MustCompile(`<[^<>]+>`).ReplaceAllString(body, ""))
}
