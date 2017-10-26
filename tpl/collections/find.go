package collections

import (
  "reflect"
)

func Find(seq interface{}, path string, comp interface{}) interface{} {
  seqv := reflect.ValueOf(seq)

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    for i := 0; i < seqv.Len(); i++ {
      obj := seqv.Index(i)

      iv := PathValue(obj, path)
      jv := reflect.ValueOf(comp)

      left, right := compareFloat(iv, jv)
      if left == right {
        return obj.Interface()
      }
    }
  }

  return nil
}
