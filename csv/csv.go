package csv

import (
  "bytes"
  "encoding/csv"
  "reflect"
)

func Unmarshal(b []byte, v interface{}) error {
  reader := bytes.NewReader(b)
  csvReader := csv.NewReader(reader)

  fc, err := csvReader.ReadAll()
  if err != nil {
    return err
  }

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

  rv := reflect.ValueOf(v).Elem()
  rv.Set(reflect.ValueOf(data))

  return nil
}
