package legacy_test

import (
	"testing"
	"testing/quick"

	"github.com/Patagonicus/suijin/legacy"
)

func TestSlice(t *testing.T) {
	f := func(slice []int) bool {
		legacy.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})

		for i := 0; i < len(slice)-1; i++ {
			if slice[i] > slice[i+1] {
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
