package template

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "strings"
  "text/template"

  "github.com/PuerkitoBio/goquery"
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/css"
  "github.com/tdewolff/minify/html"
  "github.com/tdewolff/minify/js"
)

var FuncMap = template.FuncMap {
  "ampify": Ampify,
  "file": StringifyFile,
  "minify": MinifyCode,
  "markdown": Markdown,
  "slice": Slice,
  "shim": Shim,
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
