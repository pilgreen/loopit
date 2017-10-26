package collections

import (
  "reflect"
  "strconv"
)

func compareFloat(a reflect.Value, b reflect.Value) (float64, float64) {
	var left, right float64
	var leftStr, rightStr *string

	switch a.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = float64(a.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = float64(a.Int())
	case reflect.Float32, reflect.Float64:
		left = a.Float()
	case reflect.String:
		var err error
		left, err = strconv.ParseFloat(a.String(), 64)
		if err != nil {
			str := a.String()
			leftStr = &str
		}
	}

	switch b.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = float64(b.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = float64(b.Int())
	case reflect.Float32, reflect.Float64:
		right = b.Float()
	case reflect.String:
		var err error
		right, err = strconv.ParseFloat(b.String(), 64)
		if err != nil {
			str := b.String()
			rightStr = &str
		}
	}

	switch {
	case leftStr == nil || rightStr == nil:
	case *leftStr < *rightStr:
		return 0, 1
	case *leftStr > *rightStr:
		return 1, 0
	default:
		return 0, 0
	}

	return left, right
}
