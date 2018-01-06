package suijin_test

import (
	"testing"

	"github.com/Patagonicus/suijin"
)

func TestLevelStringRoundtrip(t *testing.T) {
	for l := suijin.LogAll; l <= suijin.LogNone; l++ {
		t.Run(l.String(), func(t *testing.T) {
			result, err := suijin.LevelFromString(l.String())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if l != result {
				t.Errorf("expected %v but got %v", l, result)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	for l := suijin.LogAll; l <= suijin.LogNone; l++ {
		t.Run(l.String(), func(t *testing.T) {
			if !l.IsValid() {
				t.Error("expected level to be valid")
			}
		})
	}

	t.Run("none+1", func(t *testing.T) {
		if (suijin.LogNone + 1).IsValid() {
			t.Error("expected level to be invalid")
		}
	})
}

func TestIsSpecial(t *testing.T) {
	special := map[suijin.Level]bool{
		suijin.LogAll:  true,
		suijin.LogNone: true,
	}
	for l := suijin.LogAll; l <= suijin.LogNone; l++ {
		t.Run(l.String(), func(t *testing.T) {
			expected := special[l]
			result := l.IsSpecial()
			if expected != result {
				t.Errorf("expected %t but got %t", expected, result)
			}
		})
	}
}

func TestLevelFromString_Invalid(t *testing.T) {
	result, err := suijin.LevelFromString("invalid log level")
	if err == nil {
		t.Errorf("expected error but got %v instead", result)
	}
}
