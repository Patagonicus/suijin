package suijin_test

import (
	"strconv"
	"testing"

	"github.com/Patagonicus/suijin"
)

func TestFields_Copy(t *testing.T) {
	for i, c := range []suijin.Fields{
		{},
		{"a": 0},
		{"b": 0},
		{"a": 0, "b": 1},
		{"a": 0, "b": 1, "c": "foo"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			copy := c.Copy()

			checkFieldsEqual(t, c, copy)

			// Makes sure that this is actually a copy
			copy["copy"] = true
			if _, ok := c["copy"]; ok {
				t.Error("not actually a copy")
			}
		})
	}
}

func TestFields_Update(t *testing.T) {
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
			suijin.Fields{"a": 1},
			suijin.Fields{"a": 1},
		},
		{
			suijin.Fields{"a": 0},
			suijin.Fields{"b": 1},
			suijin.Fields{"a": 0, "b": 1},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c.a.Update(c.b)

			checkFieldsEqual(t, c.expected, c.a)
		})
	}
}

func checkFieldsEqual(t *testing.T, a, b suijin.Fields) {
	if len(a) != len(b) {
		t.Errorf("has %d entries, expected %d", len(b), len(a))
	}
	for k, v := range a {
		expected, ok := b[k]
		if !ok {
			t.Errorf("missing element %s", k)
			continue
		}
		if v != expected {
			t.Errorf("wrong value for %s, expected %d but got %d", k, expected, v)
		}
	}
}
