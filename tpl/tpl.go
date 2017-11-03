package tpl

import (
  "fmt"
  "io"
  "net/http"
  "net/url"
  "os"
  "path"
  "text/template"
)

func ParseFiles(filenames ...string) *template.Template {
  tmp := template.New(path.Base(filenames[0])).Funcs(FuncMap)
  return template.Must(tmp.ParseFiles(filenames...))
}

func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
}

func IsUrl(s string) bool {
  _, err := url.ParseRequestURI(s);
  if err != nil {
    return false
  }
  return true
}

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
