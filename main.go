package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "os"

  // Remote packages
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"

  // Local packages
  "github.com/pilgreen/loopit/tpl"
  "github.com/pilgreen/loopit/csv"
)

var version = "0.3.3"

var Config struct {
  DataFile string
  Shim bool
  Minify bool
  Version bool
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
  flag.BoolVar(&Config.Minify, "minify", false, "minifies html code")
  flag.BoolVar(&Config.Version, "v", false, "print version info")
  flag.Parse()

  var tmpls = flag.Args()
  var data interface{}
  var reader io.Reader

  // check for version flag
  if Config.Version == true {
    fmt.Println(version)
    os.Exit(1)
  }

  // check for data piped to input
  fi, err := os.Stdin.Stat()
  check(err)

  if len(Config.DataFile) > 0 {
    if(tpl.IsUrl(Config.DataFile)) {
      reader = tpl.OpenRemote(Config.DataFile)
    } else {
      reader = tpl.OpenLocal(Config.DataFile)
    }
  } else if fi.Mode() & os.ModeNamedPipe != 0 {
    reader = os.Stdin
  }

  if reader != nil {
    b, err := ioutil.ReadAll(reader)
    check(err)

    if len(b) > 0 {
      data, err = csv.ConvertToInterface(b)
      if err != nil {
        json.Unmarshal(b, &data)
      }
    }
  }

  if len(tmpls) > 0 {
    var src bytes.Buffer

    tmpl := tpl.ParseFiles(tmpls...)
    err := tmpl.Execute(&src, data)
    check(err)

    if Config.Shim {
      src, err = tpl.Shim(src)
    }

    if Config.Minify {
      minifier := minify.New()
      minifier.AddFunc("text/html", html.Minify)

      m, err := minifier.Bytes("text/html", src.Bytes())
      check(err)

      src.Reset()
      src.Write(m)
    }

    src.WriteTo(os.Stdout)
  } else {
    b, err := json.Marshal(data)
    check(err)
    os.Stdout.Write(b)
    os.Stdout.WriteString("\n")
  }
}
