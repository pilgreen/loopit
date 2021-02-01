package tpl

import (
  "fmt"
  "io"
  "net/http"
  "net/url"
  "os"
  "path"
  "text/template"

  "github.com/pilgreen/loopit/tpl/strings"
  "github.com/pilgreen/loopit/tpl/collections"
  "github.com/pilgreen/loopit/tpl/partial"
  "github.com/pilgreen/loopit/tpl/scratch"
)

/**
* Simple error checker
*/

func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
}

/**
* Appends package template functions
*/

func ParseFiles(filenames ...string) *template.Template {
  tmpl := template.New(path.Base(filenames[0]))
  tmpl.Funcs(strings.FuncMap)
  tmpl.Funcs(collections.FuncMap)
  tmpl.Funcs(partial.FuncMap)
  tmpl.Funcs(scratch.FuncMap)

  return template.Must(tmpl.ParseFiles(filenames...))
}

/**
* Remote/Local files
*/

func OpenRemote(s string) io.ReadCloser {
  resp, err := http.Get(s)
  check(err)
  return resp.Body
}

func OpenLocal(s string) *os.File {
  file, err := os.Open(s)
  check(err)
  return file
}

/**
* Simple URL check for convenience
*/

func IsUrl(s string) bool {
  u, err := url.Parse(s)
  if err != nil {
    return false
  }

  if(len(u.Host) > 0) {
    return true
  }

  return false
}

