package partial

import (
  "text/template"
)

var FuncMap = template.FuncMap {
  "partial": Partial,
  "file": File,
}
