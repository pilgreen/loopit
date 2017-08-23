package main

import (
  // "encoding/csv"
  "encoding/json"
  "flag"
  // "fmt"
  // "io"
  "io/ioutil"
  // "net/http"
  // "net/url"
  "os"
  "path"
  "path/filepath"
  "text/template"

  // Local packages
  "github.com/pilgreen/loopit/helpers"
  "github.com/pilgreen/loopit/tmpl"
)

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
    if helpers.IsUrl(*csvFile) {
      data = helpers.ParseCSV(helpers.OpenRemote(*csvFile))
    } else {
      data = helpers.ParseCSV(helpers.OpenLocal(*csvFile))
    }
  } else if len(*jsonFile) > 0 {
    var fc []byte
    var err error

    if helpers.IsUrl(*jsonFile) {
      b := helpers.OpenRemote(*jsonFile)
      fc, err = ioutil.ReadAll(b)
    } else {
      fc, err = ioutil.ReadFile(*jsonFile)
    }
    helpers.Check(err)
    json.Unmarshal(fc, &data)
  }

  files, err := filepath.Glob(*tmp)
  helpers.Check(err)

  if len(files) > 0 {
    templates := template.Must(template.New("").Funcs(tmpl.FuncMap).ParseGlob(*tmp))
    err := templates.ExecuteTemplate(os.Stdout, path.Base(files[0]), data)
    helpers.Check(err)
  } else {
    b, err := json.Marshal(data)
    helpers.Check(err)
    os.Stdout.Write(b)
  }
}
