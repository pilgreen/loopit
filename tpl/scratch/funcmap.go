package scratch

import (
  "text/template"
)

var FuncMap = template.FuncMap {
  "scratch": NewScratch,
}
