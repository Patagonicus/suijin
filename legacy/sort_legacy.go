// +build !go1.8

package legacy

import (
	"reflect"
	"sort"
)

type sliceInterface struct {
	slice reflect.Value
	len   int
	less  func(i, j int) bool
}

func (s sliceInterface) Len() int {
	return s.len
}

func (s sliceInterface) Less(i, j int) bool {
	return s.less(i, j)
}

func (s sliceInterface) Swap(i, j int) {
	valI, valJ := s.slice.Index(i), s.slice.Index(j)
	dataI, dataJ := valI.Interface(), valJ.Interface()
	valI.Set(reflect.ValueOf(dataJ))
	valJ.Set(reflect.ValueOf(dataI))
}

// Slice sorts the provided slice given the provided less function.
//
// This sort is not guaranteed to be stable.
func Slice(slice interface{}, less func(i, j int) bool) {
	value := reflect.ValueOf(slice)
	len := value.Len()
	sort.Sort(sliceInterface{
		value,
		len,
		less,
	})
}
