package main

import (
  "html/template"
)

var funcMap = template.FuncMap {
  "slice": slice,
}

func slice(a []interface{}, start, end int) []interface{} {
  return a[start:end]
}
