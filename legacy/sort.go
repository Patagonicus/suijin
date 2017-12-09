// +build go1.8

package legacy

import "sort"

// Slice sorts the provided slice given the provided less function.
//
// This sort is not guaranteed to be stable.
var Slice = sort.Slice
