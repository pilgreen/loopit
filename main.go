package main

import (
  "encoding/csv"
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path"
  "path/filepath"
  "text/template"
)

/**
 * Global variables
 */

var templates *template.Template

/**
 * Private methods
 */

func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
}

func isUrl(s string) bool {
  _, err := url.ParseRequestURI(s);
  if err != nil {
    return false
  }
  return true
}

func openRemote(s string) io.ReadCloser {
  resp, err := http.Get(s)
  check(err)
  return resp.Body
}

func openLocal(s string) *os.File {
  file, err := os.Open(s)
  check(err)
  return file
}

func csvToJson(s io.Reader) []interface{} {
  reader := csv.NewReader(s)
  fc, err := reader.ReadAll()
  check(err)

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

/**
 * Main function
 */

func main() {
  var tmp = flag.String("template", "", "path to the template file")
  var jsonFile = flag.String("json", "", "path to a JSON file")
  var csvFile = flag.String("csv", "", "path or url to a csv file")
  flag.Parse()

  var data interface{}

  if len(*csvFile) > 0 {
    if isUrl(*csvFile) {
      data = csvToJson(openRemote(*csvFile))
    } else {
      data = csvToJson(openLocal(*csvFile))
    }
  } else if len(*jsonFile) > 0 {
    var fc []byte
    var err error

    if isUrl(*jsonFile) {
      fc, err = ioutil.ReadAll(openRemote(*jsonFile))
    } else {
      fc, err = ioutil.ReadFile(*jsonFile)
    }
    check(err)
    json.Unmarshal(fc, &data)
  }

  files, err := filepath.Glob(*tmp)
  check(err)

  if len(files) > 0 {
    templates = template.Must(template.New("").Funcs(funcMap).ParseGlob(*tmp))
    err := templates.ExecuteTemplate(os.Stdout, path.Base(files[0]), data)
    check(err)
  } else {
    b, err := json.Marshal(data)
    check(err)
    os.Stdout.Write(b)
  }
}
