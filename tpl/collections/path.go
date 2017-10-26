package collections

import (
  "reflect"
  "strings"
)

func PathValue(obj reflect.Value, s string) reflect.Value {
  path := splitPath(s)

  for _, elemName := range path {
    obj, _ = evaluateSubElem(obj, elemName)
  }

  switch obj.Kind() {
  case reflect.Interface:
    return obj.Elem()
  default:
    return obj
  }
}

func splitPath(s string) []string {
  return strings.Split(strings.Trim(s, "."), ".")
}

func evaluateSubElem(obj reflect.Value, elemName string) (reflect.Value, error) {
  switch obj.Kind() {
  case reflect.Interface:
    obj = obj.Elem()
  }

  kv := reflect.ValueOf(elemName)
  return obj.MapIndex(kv), nil
}
