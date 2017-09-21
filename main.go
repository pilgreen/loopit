package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "io/ioutil"
  "os"
  "regexp"
  "text/template"

  // Local packages
  "github.com/pilgreen/loopit/funcs"
)

var Config struct {
  DataFile string
  Shim bool
}

/**
 * Main function
 */

func main() {
  flag.StringVar(&Config.DataFile, "data", "", "path or url to a JSON or CSV file")
  flag.BoolVar(&Config.Shim, "shim", false, "shims content using goquery")
  flag.Parse()

  var tmpls = flag.Args()
  var data interface{}

  if len(Config.DataFile) > 0 {
    matchCSV, err := regexp.MatchString("\\.csv$", Config.DataFile)
    funcs.Check(err)

    if matchCSV {
      if funcs.IsUrl(Config.DataFile) {
        data = funcs.ParseCSV(funcs.OpenRemote(Config.DataFile))
      } else {
        data = funcs.ParseCSV(funcs.OpenLocal(Config.DataFile))
      }
    }

    matchJSON, err := regexp.MatchString("\\.json$", Config.DataFile)
    funcs.Check(err)

    if matchJSON {
      var fc []byte
      var err error

      if funcs.IsUrl(Config.DataFile) {
        b := funcs.OpenRemote(Config.DataFile)
        fc, err = ioutil.ReadAll(b)
      } else {
        fc, err = ioutil.ReadFile(Config.DataFile)
      }

      funcs.Check(err)
      json.Unmarshal(fc, &data)
    }
  }

  if len(tmpls) > 0 {
    var src bytes.Buffer

    tmp := template.New(tmpls[0]).Funcs(funcs.FuncMap)
    templates := template.Must(tmp.ParseFiles(tmpls...))
    err := templates.Execute(&src, data)
    funcs.Check(err)

    if Config.Shim {
      src, err = funcs.Shim(src)
    }

    src.WriteTo(os.Stdout)
  } else {
    b, err := json.Marshal(data)
    funcs.Check(err)
    os.Stdout.Write(b)
    os.Stdout.WriteString("\n")
  }
}
