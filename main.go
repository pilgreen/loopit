package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"

  // Remote packages
  "github.com/fsnotify/fsnotify"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"

  // Local packages
  "github.com/pilgreen/loopit/tpl"
  "github.com/pilgreen/loopit/csv"
)

var version = "0.5.0"

type Config struct {
  DataFile string
  Minify bool
  Output string
  Shim bool
  Version bool
  Watch bool
}

func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
}

/**
 * Render template
 */

func Render(config Config, templates []string) {
  var data interface{}
  var reader io.Reader

  // check for data piped to input
  fi, err := os.Stdin.Stat()
  check(err)

  if len(config.DataFile) > 0 {
    if(tpl.IsUrl(config.DataFile)) {
      reader = tpl.OpenRemote(config.DataFile)
    } else {
      reader = tpl.OpenLocal(config.DataFile)
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

  if len(templates) > 0 {
    var src bytes.Buffer

    tmpl := tpl.ParseFiles(templates...)
    err := tmpl.Execute(&src, data)
    check(err)

    if config.Shim {
      src, err = tpl.Shim(src)
    }

    if config.Minify {
      minifier := minify.New()
      minifier.AddFunc("text/html", html.Minify)

      m, err := minifier.Bytes("text/html", src.Bytes())
      check(err)

      src.Reset()
      src.Write(m)
    }

    if len(config.Output) > 0 {
      ioutil.WriteFile(config.Output, src.Bytes(), 0644)
    } else {
      src.WriteTo(os.Stdout)
    }
  }
}

/**
 * Main function
 */

func main() {
  config := Config{}

  flag.StringVar(&config.DataFile, "data", "", "path or url to a JSON or CSV file")
  flag.BoolVar(&config.Minify, "minify", false, "minifies html code")
  flag.StringVar(&config.Output, "out", "", "output file")
  flag.BoolVar(&config.Shim, "shim", false, "shims content using goquery")
  flag.BoolVar(&config.Version, "version", false, "print version info")
  flag.BoolVar(&config.Watch, "watch", false, "runs loopit on file changes")
  flag.Parse()

  // check for version flag
  if config.Version == true {
    fmt.Println(version)
    os.Exit(1)
  }

  var templates = flag.Args()

  // Always run the initial render
  Render(config, templates)

  // Set up fsnotify to watch the directory
  if config.Watch == true {
    watcher, err := fsnotify.NewWatcher()
    check(err)
    defer watcher.Close()

    done := make(chan bool)
    go func() {
      for {
        select {
        case <-watcher.Events:
          Render(config, templates)
        case err := <-watcher.Errors:
          fmt.Fprintln(os.Stderr, err)
        }
      }
    }()

    for _, v := range templates {
      watcher.Add(filepath.Dir(v));
    }

    if len(config.DataFile) > 0 && !tpl.IsUrl(config.DataFile) {
      watcher.Add(filepath.Dir(config.DataFile));
    }

    <-done
  }
}
