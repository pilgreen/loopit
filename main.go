package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "os"

  // Local packages
  "github.com/pilgreen/loopit/template"
  "github.com/pilgreen/loopit/csv"
)

var Config struct {
  DataFile string
  Shim bool
}

func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
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
  check(err)

  if(fi.Mode() & os.ModeNamedPipe != 0) {
    reader = os.Stdin
  } else if len(Config.DataFile) > 0 {
    if(template.IsUrl(Config.DataFile)) {
      reader = template.OpenRemote(Config.DataFile)
    } else {
      reader = template.OpenLocal(Config.DataFile)
    }
  }

  if reader != nil {
    b, err := ioutil.ReadAll(reader)
    check(err)

    data, err = csv.ConvertToInterface(b)
    if err != nil {
      json.Unmarshal(b, &data)
    }
  }

  if len(tmpls) > 0 {
    var src bytes.Buffer

    tmpl := template.ParseFiles(tmpls...)
    err := tmpl.Execute(&src, data)
    check(err)

    if Config.Shim {
      src, err = template.Shim(src)
    }

    src.WriteTo(os.Stdout)
  } else {
    b, err := json.Marshal(data)
    check(err)
    os.Stdout.Write(b)
    os.Stdout.WriteString("\n")
  }
}
