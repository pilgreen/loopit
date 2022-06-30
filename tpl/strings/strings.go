package strings

import (
  "bytes"
  "encoding/json"
  "net/url"
  "regexp"
  "strings"
  "time"

  "github.com/dustin/go-humanize"
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/css"
  "github.com/tdewolff/minify/html"
  "github.com/tdewolff/minify/js"
)

/**
 * Returns the minified version of a string
 * You must set the mimetype manually
 */

func MinifyCode(mimetype string, code string) (string, error) {
  m := minify.New()
  m.AddFunc("text/css", css.Minify)
  m.AddFunc("text/html", html.Minify)
  m.AddFunc("text/javascript", js.Minify)
  return m.String(mimetype, code)
}

/**
 * Runs a string through blackfriday
 */

func Markdown(s string) string {
  bits := []byte(s)
  newBits := blackfriday.MarkdownCommon(bits)
  return bytes.NewBuffer(newBits).String()
}

/**
 * Converts float64 to an int (for comparison)
 */

func FloatToInt(n float64) int {
  return int(n)
}

/**
 * Passes a date string through the time package
 * Order of arguments: layout, format, timezone
 */

func DateFormat(date string, args ...string) string {
  t, err := time.Parse(args[0], date)
  if err != nil {
    return date
  }

  if len(args) > 2 {
    loc, err := time.LoadLocation(args[2])
    if err != nil {
      return date
    }
    return t.In(loc).Format(args[1])
  }

  if args[1] == "humanize" {
    return humanize.Time(t)
  }

  return t.Format(args[1])
}

/**
 * Addition/Subtraction
 */

func Add(add int, initial int) int {
  return initial + add
}

func Subtract(sub int, initial int) int {
  return initial - sub
}

/**
 * Replace a portion of a string
 */

func Replace(from, to, input string) string {
  return strings.Replace(input, from, to, -1)
}

/**
 * Regex functions
 */

func MatchRe(pattern string, input string) bool {
  re := regexp.MustCompile(pattern)
  return re.MatchString(input)
}

func FindRe(pattern string, input string) string {
  re := regexp.MustCompile(pattern)
  return re.FindString(input)
}

func FindSubRe(pattern string, input string) []string {
  re := regexp.MustCompile(pattern)
  return re.FindStringSubmatch(input)
}

func ReplaceRe(pattern string, to string, input string) string {
  re := regexp.MustCompile(pattern)
  return re.ReplaceAllString(input, to)
}

/**
 * String conversions
 */

func ToLower(s string) string {
  return strings.ToLower(s)
}

func ToUpper(s string) string {
  return strings.ToUpper(s)
}

func Unescape(s string) (string, error) {
  return url.QueryUnescape(s)
}

/**
 * Trims characters from both sides
 * Note: the cutset is required
 */

func Trim(cutset string, s string) string {
  return strings.Trim(s, cutset)
}

/**
 * Returns a marshaled JSON string
 */

func Marshal(s string) string {
  b, _ := json.Marshal(s)
  return string(b)
}
