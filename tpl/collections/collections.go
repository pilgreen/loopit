package collections

import (
  "fmt"
  "reflect"
  "regexp"
  "strings"
  "github.com/spf13/cast"
)

/**
* Dot notation to pull nested elements
*/

func PathValue(obj reflect.Value, s string) reflect.Value {
  path := strings.Split(strings.Trim(s, "."), ".")

  for _, name := range path {
    if obj.Kind() == reflect.Interface {
      obj = obj.Elem()
    }

    kv := reflect.ValueOf(name)
    obj = obj.MapIndex(kv)
  }

  switch obj.Kind() {
  case reflect.Interface:
    return obj.Elem()
  default:
    return obj
  }
}

/**
* Sort a sequence by value
*/

func Sort(seq interface{}, args ...string) (interface{}, error) {
  seqv := reflect.ValueOf(seq)

  // Create a list of pairs that will be used to do the sort
	p := pairList{SortAsc: true, SliceType: reflect.SliceOf(seqv.Type().Elem())}
	p.Pairs = make([]pair, seqv.Len())

  var sortByField string
  for i, l := range args {
    dStr, err := cast.ToStringE(l)
    switch {
    case i == 0 && err !=nil:
      sortByField = ""
    case i == 0 && err == nil:
      sortByField = dStr
    case i == 1 && err == nil && dStr == "desc":
			p.SortAsc = false
		case i == 1:
			p.SortAsc = true
    }
  }

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    for i := 0; i < seqv.Len(); i++ {
      p.Pairs[i].Value = seqv.Index(i)
      if sortByField == "" || sortByField == "value" {
        p.Pairs[i].Key = p.Pairs[i].Value
      } else {
        v := PathValue(p.Pairs[i].Value, sortByField)
        p.Pairs[i].Key = v
      }
    }

  case reflect.Map:
    keys := seqv.MapKeys()
    for i := 0; i < seqv.Len(); i++ {
      p.Pairs[i].Value = seqv.MapIndex(keys[i])
      if sortByField == "" {
        p.Pairs[i].Key = keys[i]
      } else if sortByField == "value" {
        p.Pairs[i].Key = p.Pairs[i].Value
      } else {
        v := PathValue(p.Pairs[i].Value, sortByField)
        p.Pairs[i].Key = v
      }
    }
  }

  return p.sort(), nil
}

/**
* Filters a sequence
*/

func Where(seq interface{}, path string, args ...interface{}) interface{} {
  seqv := reflect.ValueOf(seq)

  var comp, val interface{}
  switch len(args) {
  case 1:
    comp = "=="
    val = args[0]
  case 2:
    comp = args[0]
    val = args[1]
  }

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    rs := reflect.MakeSlice(seqv.Type(), 0, 0)

    for i := 0; i < seqv.Len(); i++ {
      obj := seqv.Index(i)

      pv := PathValue(obj, path)
      rv := reflect.ValueOf(val)
      left, right := compare(pv, rv)

      switch comp {
      case "=", "==", "eq":
        if left == right {
          rs = reflect.Append(rs, obj)
        }
      case "!=", "<>", "ne":
        if left != right {
          rs = reflect.Append(rs, obj)
        }
      case "<", "lt":
        if int(left) < int(right) {
          rs = reflect.Append(rs, obj)
        }
      case ">", "gt":
        if int(left) > int(right) {
          rs = reflect.Append(rs, obj)
        }
      case "<=", "lte":
        fmt.Println(int(left), int(right))
        if int(left) <= int(right) {
          rs = reflect.Append(rs, obj)
        }
      case ">=", "gte":
        if int(left) >= int(right) {
          rs = reflect.Append(rs, obj)
        }
      case "matches":
        re := regexp.MustCompile(rv.String())
        if re.MatchString(pv.String()) {
          rs = reflect.Append(rs, obj)
        }
      }
    }

    return rs.Interface()
  }

  return nil
}

/**
* Returns the first matching element of a sequence
*/

func Find(seq interface{}, path string, args ...interface{}) interface{} {
  seqv := reflect.ValueOf(seq)

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    matches := reflect.ValueOf(Where(seq, path, args...))
    if matches.Len() > 0 {
      return matches.Index(0).Interface()
    }
  }

  return nil
}


/**
* Joins an array of strings
*/

func Join(sep string, seq interface{}) string {
  seqv := reflect.ValueOf(seq)
  str := make([]string, seqv.Len())

  for i := 0; i < seqv.Len(); i++ {
    str[i] = fmt.Sprintf("%s", seqv.Index(i))
  }

  return strings.Join(str, sep)
}

/**
 * Pulls a subset from a sequence 
 * Good for range to pull things like top 10
 */

func Slice(start, end int, seq interface{}) interface{} {
  seqv := reflect.ValueOf(seq)

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice, reflect.String:
    return seqv.Slice(start, end).Interface()
  }

  return nil
}

/**
* Creates a slice of values from a sequence
*/

func Pluck(path string, seq interface{}) interface{} {
  seqv := reflect.ValueOf(seq)

  switch seqv.Kind() {
  case reflect.Array, reflect.Slice:
    slice := make([]interface{}, seqv.Len())
    for i := 0; i < seqv.Len(); i++ {
      obj := seqv.Index(i)
      slice[i] = PathValue(obj, path)
    }
    return slice
  case reflect.Map:
    return PathValue(seqv, path)
  }

  return nil
}
