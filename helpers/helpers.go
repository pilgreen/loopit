package helpers

import (
  "encoding/csv"
  "fmt"
  "io"
  "net/http"
  "net/url"
  "os"
)

func Check(e error) {
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
  Check(err)
  return resp.Body
}

func OpenLocal(s string) *os.File {
  file, err := os.Open(s)
  Check(err)
  return file
}

func ParseCSV(s io.Reader) []interface{} {
  reader := csv.NewReader(s)
  fc, err := reader.ReadAll()
  Check(err)

  var data []interface{}
  header := fc[0]
  for _, row := range fc[1:] {
    obj := make(map[string]interface{}, len(header))
    for j, v := range row {
      key := header[j]
      obj[key] = v
    }
    data = append(data, obj)
  }
  return data
}
