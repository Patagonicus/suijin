package suijin_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/Patagonicus/suijin"
)

func TestLogger_WithField(t *testing.T) {
	testLogger(t, func(l suijin.Logger, fds suijin.Fields) suijin.Logger {
		for k, v := range fds {
			l = l.WithField(k, v)
		}
		return l
	})
}

func TestLogger_WithFields(t *testing.T) {
	testLogger(t, func(l suijin.Logger, fds suijin.Fields) suijin.Logger {
		return l.WithFields(fds)
	})
}

func testLogger(t *testing.T, addFields func(l suijin.Logger, fds suijin.Fields) suijin.Logger) {
	for i, c := range []suijin.Message{
		{
			Level:   suijin.DebugLevel,
			Message: "debug",
			Fields:  nil,
		},
		{
			Level:   suijin.InfoLevel,
			Message: "info",
			Fields:  nil,
		},
		{
			Level:   suijin.WarningLevel,
			Message: "warning",
			Fields:  nil,
		},
		{
			Level:   suijin.ErrorLevel,
			Message: "error",
			Fields:  nil,
		},
		{
			Level:   suijin.DebugLevel,
			Message: "non-nil empty fields",
			Fields:  suijin.Fields{},
		},
		{
			Level:   suijin.DebugLevel,
			Message: "single field",
			Fields:  suijin.Fields{"a": 0},
		},
		{
			Level:   suijin.InfoLevel,
			Message: "multiple fields",
			Fields:  suijin.Fields{"a": 0, "b": 1},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := &sink{}
			l := suijin.NewLogger(s)
			if c.Fields != nil {
				l = addFields(l, c.Fields)
			}

			log(t, l, c.Level, c.Message)

			if len(s.msgs) != 1 {
				t.Fatalf("expected 1 message but got %d: %v", len(s.msgs), s.msgs)
			}

			if c.Fields == nil {
				c.Fields = make(suijin.Fields)
			}
			if !reflect.DeepEqual(c, s.msgs[0]) {
				t.Errorf("expected %v but got %v", c, s.msgs[0])
			}
		})
	}
}

func TestLogger_Multiple(t *testing.T) {
	s := &sink{}
	l := suijin.NewLogger(s)

	l.Debug("foo")
	l.WithField("test", 1).Info("bar")
	l.WithFields(suijin.Fields{
		"foo": "bar",
		"42":  42,
	}).Warning("baz")
	l.WithField("test", 2).WithField("test", 3).Error("foobar")

	expected := []suijin.Message{
		{Level: suijin.DebugLevel, Message: "foo", Fields: suijin.Fields{}},
		{Level: suijin.InfoLevel, Message: "bar", Fields: suijin.Fields{"test": 1}},
		{Level: suijin.WarningLevel, Message: "baz", Fields: suijin.Fields{"foo": "bar", "42": 42}},
		{Level: suijin.ErrorLevel, Message: "foobar", Fields: suijin.Fields{"test": 3}},
	}
	if !reflect.DeepEqual(s.msgs, expected) {
		t.Errorf("expected %v but got %v", expected, s.msgs)
	}
}

func log(t *testing.T, l suijin.Logger, lvl suijin.Level, msg string) {
	switch lvl {
	case suijin.DebugLevel:
		l.Debug(msg)
	case suijin.InfoLevel:
		l.Info(msg)
	case suijin.WarningLevel:
		l.Warning(msg)
	case suijin.ErrorLevel:
		l.Error(msg)
	default:
		t.Fatalf("unknown level: %v", lvl)
	}
}
