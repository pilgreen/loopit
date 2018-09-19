package partial

import (
  "text/template"
)

var FuncMap = template.FuncMap {
  "file": File,
  "partial": Partial,
}
