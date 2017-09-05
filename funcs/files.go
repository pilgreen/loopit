package funcs

import (
  "io"
  "net/http"
  "net/url"
  "os"
)

func IsUrl(s string) bool {
  _, err := url.ParseRequestURI(s);
  if err != nil {
    return false
  }
  return true
}

func OpenRemote(s string) io.ReadCloser {
  resp, err := http.Get(s)
  Check(err)
  return resp.Body
}

func OpenLocal(s string) *os.File {
  file, err := os.Open(s)
  Check(err)
  return file
}
