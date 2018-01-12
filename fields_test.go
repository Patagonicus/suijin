package suijin_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/Patagonicus/suijin"
)

func TestFields_AddAll(t *testing.T) {
	for i, c := range []struct {
		a, b     suijin.Fields
		expected suijin.Fields
	}{
		{
			suijin.Fields{},
			suijin.Fields{},
			suijin.Fields{},
		},
		{
			suijin.Fields{"a": 0},
			suijin.Fields{},
			suijin.Fields{"a": 0},
		},
		{
			suijin.Fields{},
			suijin.Fields{"a": 0},
			suijin.Fields{"a": 0},
		},
		{
			suijin.Fields{"a": 0},
			suijin.Fields{"b": 1},
			suijin.Fields{"a": 0, "b": 1},
		},
		{
			suijin.Fields{"a": 0, "b": 1, "c": 2},
			suijin.Fields{"a": 0, "b": -1, "d": 3},
			suijin.Fields{"a": 0, "b": -1, "c": 2, "d": 3},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c.a.AddAll(c.b)
			if !reflect.DeepEqual(c.a, c.expected) {
				t.Errorf("invalid result: %v", c.a)
			}
		})
	}
}
