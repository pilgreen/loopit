package collections

import (
  "reflect"
)

func Find(seq interface{}, path string, args ...interface{}) interface{} {
  seqv := reflect.ValueOf(seq)
  comp, val := processFindArgs(args)

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    for i := 0; i < seqv.Len(); i++ {
      obj := seqv.Index(i)

      iv := PathValue(obj, path)
      jv := reflect.ValueOf(val)

      left, right := compareFloat(iv, jv)

      switch comp {
      case "!=":
        if left != right {
          return obj.Interface()
        }
      default:
        if left == right {
          return obj.Interface()
        }
      }
    }
  }

  return nil
}

func processFindArgs(args []interface{}) (comp, val interface{}) {
  switch len(args) {
  case 1:
    comp = ""
    val = args[0]
  case 2:
    comp = args[0]
    val = args[1]
  }

  return comp, val
}
