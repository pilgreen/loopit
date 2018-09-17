package partial

import (
  "bytes"
  "io/ioutil"
  "net/http"
  "net/url"
  "path"
  "text/template"

  "github.com/pilgreen/loopit/tpl/strings"
  "github.com/pilgreen/loopit/tpl/collections"
  "github.com/pilgreen/loopit/tpl/scratch"
)

/**
 * Returns a template as a string
 * Partial has no scope other than itself including all loopit functions.
 */

func Partial(name string, data interface{}) (string, error) {
  tmpl := template.New(path.Base(name))
  tmpl.Funcs(strings.FuncMap)
  tmpl.Funcs(collections.FuncMap)
  tmpl.Funcs(scratch.FuncMap)
  template.Must(tmpl.ParseFiles(name))

  var b bytes.Buffer
  err := tmpl.Execute(&b, data)
  if err != nil { return "", err }

  return b.String(), nil
}

/**
 * Returns a file or url as a string
 */

func File(name string) (string, error) {
  var fc []byte

  _, err := url.ParseRequestURI(name);
  if err != nil {
    fc, err = ioutil.ReadFile(name)
    if err != nil { return "", err }
  } else {
    resp, err := http.Get(name)
    if err != nil { return "", err }
    fc, err = ioutil.ReadAll(resp.Body)
    if err != nil { return "", err }
  }

  return bytes.NewBuffer(fc).String(), nil
}
