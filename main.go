package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "io"
  "io/ioutil"
  "os"
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
  var reader io.Reader

  // check for data piped to input
  fi, err := os.Stdin.Stat()
  funcs.Check(err)

  if(fi.Mode() & os.ModeNamedPipe != 0) {
    reader = os.Stdin
  } else if len(Config.DataFile) > 0 {
    if(funcs.IsUrl(Config.DataFile)) {
      reader = funcs.OpenRemote(Config.DataFile)
    } else {
      reader = funcs.OpenLocal(Config.DataFile)
    }
  }

  b, err := ioutil.ReadAll(reader)
  funcs.Check(err)

  data, err = funcs.ParseCSV(b)
  if err != nil {
    json.Unmarshal(b, &data)
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
