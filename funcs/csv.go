package funcs

import (
  "encoding/csv"
  "io"
)

func ParseCSV(s io.Reader) []interface{} {
  reader := csv.NewReader(s)
  fc, err := reader.ReadAll()
  Check(err)

  var data []interface{}
  header := fc[0]
  for _, row := range fc[1:] {
    obj := make(map[string]interface{}, len(header))
    for j, v := range row {
      key := header[j]
      obj[key] = v
    }
    data = append(data, obj)
  }
  return data
}
