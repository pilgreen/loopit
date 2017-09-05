package funcs

import (
  "fmt"
  "os"
)

func Check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e)
    os.Exit(1)
  }
}
