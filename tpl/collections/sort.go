package collections

import (
  "reflect"
  "sort"

  "github.com/spf13/cast"
)

// Sort returns a sorted sequence.
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

// Credit for pair sorting method goes to Andrew Gerrand
// https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
// A data structure to hold a key/value pair.
type pair struct {
	Key   reflect.Value
	Value reflect.Value
}

// A slice of pairs that implements sort.Interface to sort by Value.
type pairList struct {
	Pairs     []pair
	SortAsc   bool
	SliceType reflect.Type
}

func (p pairList) Swap(i, j int) { p.Pairs[i], p.Pairs[j] = p.Pairs[j], p.Pairs[i] }
func (p pairList) Len() int      { return len(p.Pairs) }
func (p pairList) Less(i, j int) bool {
	iv := p.Pairs[i].Key
	jv := p.Pairs[j].Key

  left, right := compareFloat(iv, jv)
  return left < right
}

// sorts a pairList and returns a slice of sorted values
func (p pairList) sort() interface{} {
	if p.SortAsc {
		sort.Sort(p)
	} else {
		sort.Sort(sort.Reverse(p))
	}
	sorted := reflect.MakeSlice(p.SliceType, len(p.Pairs), len(p.Pairs))
	for i, v := range p.Pairs {
		sorted.Index(i).Set(v.Value)
	}

	return sorted.Interface()
}
