package csv

import (
  "bytes"
  "encoding/csv"
)

func ConvertToInterface(b []byte) ([]interface{}, error) {
  var data []interface{}
  reader := bytes.NewReader(b)

  csvReader := csv.NewReader(reader)
  fc, err := csvReader.ReadAll()
  if err != nil {
    return data, err
  }

  header := fc[0]
  for _, row := range fc[1:] {
    obj := make(map[string]interface{}, len(header))
    for j, v := range row {
      key := header[j]
      obj[key] = v
    }
    data = append(data, obj)
  }

  return data, nil
}
