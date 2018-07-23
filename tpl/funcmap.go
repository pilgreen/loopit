package tpl

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "strings"
  "text/template"
  "time"

  "github.com/pilgreen/loopit/tpl/collections"
  "github.com/pilgreen/loopit/tpl/conversions"

  "github.com/PuerkitoBio/goquery"
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/css"
  "github.com/tdewolff/minify/html"
  "github.com/tdewolff/minify/js"
)

var FuncMap = template.FuncMap {
  "add": Add,
  "ampify": Ampify,
  "dateFormat": DateFormat,
  "file": StringifyFile,
  "find": collections.Find,
  "inchesToFeet": conversions.InchesToFeet,
  "int": Int,
  "lower": ToLower,
  "minify": MinifyCode,
  "markdown": Markdown,
  "replace": Replace,
  "slice": Slice,
  "sort": collections.Sort,
  "subtract": Subtract,
  "trim": Trim,
  "upper": ToUpper,
}

/**
 * Returns a subset of some data
 * Good for range to pull things like top 10
 */

func Slice(a []interface{}, start, end int) []interface{} {
  return a[start:end]
}

/**
 * Returns a file or url as a string
 */

func StringifyFile(name string) string {
  var fc []byte
  var err error

  if IsUrl(name) {
    b := OpenRemote(name)
    fc, err = ioutil.ReadAll(b)
  } else {
    fc, err = ioutil.ReadFile(name)
  }

  check(err)
  return bytes.NewBuffer(fc).String()
}

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
 * Converts iframes to amp-iframes
 */

func Ampify(s string) (string, error) {
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
  check(err)

  iframes := doc.Find("iframe")
  iframes.Each(func(i int, ele *goquery.Selection) {
    src, exists := ele.Attr("src");
    if exists == true {
      width := ele.AttrOr("width", "16")
      height := ele.AttrOr("height", "9")
      amp := fmt.Sprintf("<amp-iframe width='%s' height='%s' layout='responsive' sandbox='allow-scripts allow-same-origin' src='%s'></amp-iframe>", width, height, src)
      ele.ReplaceWithHtml(amp)
    }
  })
  return doc.Find("body").Html()
}

/**
 * Converts float64 to an int (for comparison)
 */

func Int(n float64) int {
  return int(n)
}

/**
 * Passes an RFC3339 formatted date string through the default parser
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
  return t.Format(args[1])
}

/**
 * Adds to a number
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
 * String conversions
 */

func ToLower(s string) string {
  return strings.ToLower(s)
}

func ToUpper(s string) string {
  return strings.ToUpper(s)
}

/**
 * Trims characters from both sides
 * Note: the cutset is required
 */

func Trim(cutset string, s string) string {
  return strings.Trim(s, cutset)
}
