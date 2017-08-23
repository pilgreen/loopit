package main

import (
  "bytes"
  "text/template"

  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/css"
  "github.com/tdewolff/minify/html"
  "github.com/tdewolff/minify/js"
)

var funcMap = template.FuncMap {
  "slice": slice,
  "stringify": stringifyTemplate,
  "minify": minifyCode,
  "markdown": markdown,
}

/**
 * Returns a subset of some data
 * Good for range to pull things like top 10
 */

func slice(a []interface{}, start, end int) []interface{} {
  return a[start:end]
}

/**
 * Returns template contents as a string
 * use to pass contents to minify
 */

func stringifyTemplate(name string) string {
  var doc bytes.Buffer
  templates.ExecuteTemplate(&doc, name, nil)
  return doc.String()
}

/**
 * Returns the minified version of a string
 * You must set the mimetype manually
 */

func minifyCode(mimetype string, code string) (string, error) {
  m := minify.New()
  m.AddFunc("text/css", css.Minify)
  m.AddFunc("text/html", html.Minify)
  m.AddFunc("text/javascript", js.Minify)
  return m.String(mimetype, code)
}

/**
 * Runs a string through blackfriday
 */

func markdown(s string) string {
  bits := []byte(s)
  newBits := blackfriday.MarkdownCommon(bits)
  return bytes.NewBuffer(newBits).String()
}
