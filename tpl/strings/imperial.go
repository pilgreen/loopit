package strings

import (
  "fmt"
  "math"
)

func InchesToFeet(inches interface{}) (string, error) {
  float, ok := inches.(float64)
  if ok {
    ft := math.Floor(float/12)
    in := math.Mod(float, 12)
    return fmt.Sprintf("%d'%d\"", int(ft), int(in)), nil
  }

  return "", fmt.Errorf("cannot convert %s to a float64", inches)
}
