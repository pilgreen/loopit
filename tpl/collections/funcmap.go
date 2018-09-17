package collections

import (
  "text/template"
)

var FuncMap = template.FuncMap {
  "find": Find,
  "join": Join,
  "pluck": Pluck,
  "slice": Slice,
  "sort": Sort,
  "where": Where,
}
