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
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"

  // Local packages
  "github.com/pilgreen/loopit/tpl"
  "github.com/pilgreen/loopit/csv"
)

var version = "0.6.0"

type Config struct {
  DataFile string
  Minify bool
  Markdown bool
  Output string
  Shim bool
  Version bool
  Watch string
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

    if config.Markdown {
      m := blackfriday.MarkdownCommon(src.Bytes())

      src.Reset()
      src.Write(m)
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

  flag.StringVar(&config.DataFile, "d", "", "path or url to a JSON or CSV file")
  flag.BoolVar(&config.Markdown, "markdown", false, "run output through BlackFriday")
  flag.BoolVar(&config.Minify, "minify", false, "minifies html code")
  flag.StringVar(&config.Output, "o", "", "output file")
  flag.BoolVar(&config.Shim, "shim", false, "shims content using goquery")
  flag.BoolVar(&config.Version, "v", false, "version info")
  flag.StringVar(&config.Watch, "w", "", "glob pattern to watch for changes")
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
  if len(config.Watch) > 0 {
    watcher, err := fsnotify.NewWatcher()
    check(err)
    defer watcher.Close()

    done := make(chan bool)
    go func() {
      for {
        select {
        case e := <-watcher.Events:
          if e.Op.String() == "WRITE" {
            Render(config, templates)
          }
        case err := <-watcher.Errors:
          fmt.Fprintln(os.Stderr, err)
        }
      }
    }()

    files, err := filepath.Glob(config.Watch)
    check(err)

    for _, v := range files {
      watcher.Add(v);
    }

    <-done
  }
}
