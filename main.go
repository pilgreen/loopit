package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "path/filepath"
  "regexp"
  "strings"

  // Remote packages
  "github.com/fsnotify/fsnotify"
  "github.com/russross/blackfriday"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"
  "github.com/yosssi/gohtml"

  // Local packages
  "github.com/pilgreen/loopit/tpl"
  "github.com/pilgreen/loopit/csv"
  "github.com/pilgreen/loopit/rss"
)

var version = "0.8.0"

type Config struct {
  DataFile string
  Minify bool
  Markdown bool
  Output string
  Shim bool
  Tidy bool
  Version bool
  Watch bool
  Port string
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

  // Check for data piped to input
  fi, err := os.Stdin.Stat()
  check(err)

  // Figure out where the data is coming from
  if len(config.DataFile) > 0 {
    if(tpl.IsUrl(config.DataFile)) {
      reader = tpl.OpenRemote(config.DataFile)
    } else {
      reader = tpl.OpenLocal(config.DataFile)
    }
  } else if fi.Mode() & os.ModeNamedPipe != 0 {
    reader = os.Stdin
  }

  // Unmarshal the data: accepts json, csv, rss, atom
  if reader != nil {
    b, err := ioutil.ReadAll(reader)
    check(err)

    // See if we have an rss or atom feed
    mime := http.DetectContentType(b)

    if strings.Contains(mime, "text/xml") {
      rss.Unmarshal(b, &data)
    } else {
      // JSON and XML come in as "text/plain"
      if json.Valid(b) {
        json.Unmarshal(b, &data)
      } else {
        csv.Unmarshal(b, &data)
      }
    }
  }

  if len(templates) > 0 {
    var src bytes.Buffer

    template := tpl.ParseFiles(templates...)
    err := template.Execute(&src, data)
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

    if config.Tidy {
      gohtml.Condense = true
      t := gohtml.FormatBytes(src.Bytes())
      src.Reset()
      src.Write(t)
    }

    if len(config.Output) > 0 {
      ioutil.WriteFile(config.Output, src.Bytes(), 0644)
    } else {
      src.WriteTo(os.Stdout)
    }
  } else {
    // No template let's kick out the JSON for developing
    b, _ := json.Marshal(data)
    os.Stdout.Write(b)
  }

}

/**
 * Main function
 */

func main() {
  config := Config{}

  flag.StringVar(&config.DataFile, "data", "", "path or url to a JSON or CSV file")
  flag.BoolVar(&config.Markdown, "markdown", false, "run output through BlackFriday")
  flag.BoolVar(&config.Minify, "minify", false, "minifies html code")
  flag.StringVar(&config.Output, "out", "", "output file")
  flag.BoolVar(&config.Shim, "shim", false, "shims content using goquery")
  flag.BoolVar(&config.Tidy, "tidy", false, "cleanup the output")
  flag.BoolVar(&config.Version, "version", false, "version info")
  flag.BoolVar(&config.Watch, "watch", false, "rebuilds on file changes and starts a server (out flag required)")
  flag.StringVar(&config.Port, "port", "1313", "port to use for the local server")
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
  if config.Watch == true && len(config.Output) > 0 {
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

    // watch template files
    for _, t := range templates {
      watcher.Add(t);
    }

    // also watch css and js files
    css := regexp.MustCompile(".*\\.css$")
    js := regexp.MustCompile(".*\\.js$")

    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
      name := info.Name()
      if css.MatchString(name) || js.MatchString(name) {
        watcher.Add(path)
      }
      return nil
    })

    // start the server
    dir, _ := os.Getwd()
    fmt.Printf("Server started at http://localhost:%s\n", config.Port)
    http.ListenAndServe(fmt.Sprintf(":%s", config.Port), http.FileServer(http.Dir(dir)))

    <-done
  }
}
