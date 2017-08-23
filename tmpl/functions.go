package tmpl

import (
  "bytes"
  "io/ioutil"
  "text/template"

  "github.com/pilgreen/loopit/helpers"
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/css"
  "github.com/tdewolff/minify/html"
  "github.com/tdewolff/minify/js"
)

var FuncMap = template.FuncMap {
  "slice": Slice,
  "file": StringifyFile,
  "minify": MinifyCode,
  "markdown": Markdown,
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

  if helpers.IsUrl(name) {
    b := helpers.OpenRemote(name)
    fc, err = ioutil.ReadAll(b)
  } else {
    fc, err = ioutil.ReadFile(name)
  }

  helpers.Check(err)
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
