package suijin_test

import (
	"bytes"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Patagonicus/suijin"
)

type sink struct {
	msgs []suijin.Message
}

func (s *sink) Log(msg suijin.Message) {
	s.msgs = append(s.msgs, msg)
}

const concurrency = 8

var writerCases = []struct {
	messages []suijin.Message
	expected string
}{
	{
		messages: nil,
		expected: "",
	},
	{
		messages: []suijin.Message{
			{
				Level:   suijin.DebugLevel,
				Message: "foo",
			},
		},
		expected: "debug: foo\n",
	},
	{
		messages: []suijin.Message{
			{
				Level:   suijin.InfoLevel,
				Message: "bar",
			},
		},
		expected: "info: bar\n",
	},
	{
		messages: []suijin.Message{
			{
				Level:   suijin.InfoLevel,
				Message: "foobar",
			},
			{
				Level:   suijin.ErrorLevel,
				Message: "barfoo",
			},
		},
		expected: "info: foobar\nerror: barfoo\n",
	},
	{
		messages: []suijin.Message{
			{
				Level:   suijin.InfoLevel,
				Message: "foobar",
				Fields:  suijin.Fields{"a": 0},
			},
		},
		expected: "info: foobar 'a'='0'\n",
	},
	{
		messages: []suijin.Message{
			{
				Level:   suijin.InfoLevel,
				Message: "foobar",
				Fields:  suijin.Fields{"a": 0, "b": 1},
			},
		},
		expected: "info: foobar 'a'='0' 'b'='1'\n",
	},
}

func TestWriterSink(t *testing.T) {
	for i, c := range writerCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := &bytes.Buffer{}
			sink := suijin.Writer(buf)
			for _, m := range c.messages {
				sink.Log(m)
			}

			if buf.String() != c.expected {
				t.Fatalf("expected '%v' but got '%v'", c.expected, buf.String())
			}
		})
	}
}

func TestWriterSinkConcurrent(t *testing.T) {
	buf := &bytes.Buffer{}
	backend := suijin.Writer(buf)

	start, stop := &sync.WaitGroup{}, &sync.WaitGroup{}
	start.Add(concurrency)
	stop.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(nr int) {
			start.Done()
			start.Wait()

			for j := 0; j < 100; j++ {
				backend.Log(suijin.Message{
					Level:   suijin.DebugLevel,
					Message: "test",
					Fields: suijin.Fields{
						"worker": nr,
						"j":      j,
					},
				})
			}

			stop.Done()
		}(i)
	}

	stop.Wait()

	expected := &bytes.Buffer{}
	expectedBackend := suijin.Writer(expected)
	for nr := 0; nr < concurrency; nr++ {
		for j := 0; j < 100; j++ {
			expectedBackend.Log(suijin.Message{
				Level:   suijin.DebugLevel,
				Message: "test",
				Fields: suijin.Fields{
					"worker": nr,
					"j":      j,
				},
			})
		}
	}

	actualLines := strings.Split(expected.String(), "\n")
	expectedLines := strings.Split(expected.String(), "\n")
	sort.Strings(actualLines)
	sort.Strings(expectedLines)

	if !reflect.DeepEqual(expectedLines, actualLines) {
		t.Error("messages are not what I expected")
	}
}

func TestFieldSink(t *testing.T) {
	for i, c := range []struct {
		msg      suijin.Message
		fields   suijin.Fields
		expected suijin.Message
	}{
		{
			suijin.Message{
				Level:   suijin.DebugLevel,
				Message: "test",
				Fields:  suijin.Fields{},
			},
			nil,
			suijin.Message{
				Level:   suijin.DebugLevel,
				Message: "test",
				Fields:  suijin.Fields{},
			},
		},
		{
			suijin.Message{
				Level:   suijin.DebugLevel,
				Message: "test",
				Fields:  suijin.Fields{},
			},
			suijin.Fields{},
			suijin.Message{
				Level:   suijin.DebugLevel,
				Message: "test",
				Fields:  suijin.Fields{},
			},
		},
		{
			suijin.Message{
				Level:   suijin.InfoLevel,
				Message: "foo",
				Fields:  nil,
			},
			suijin.Fields{"a": 0},
			suijin.Message{
				Level:   suijin.InfoLevel,
				Message: "foo",
				Fields:  suijin.Fields{"a": 0},
			},
		},
		{
			suijin.Message{
				Level:   suijin.InfoLevel,
				Message: "bar",
				Fields:  suijin.Fields{"a": 0, "b": 1},
			},
			suijin.Fields{"b": 2, "c": 3},
			suijin.Message{
				Level:   suijin.InfoLevel,
				Message: "bar",
				Fields:  suijin.Fields{"a": 0, "b": 2, "c": 3},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := &sink{}
			fs := suijin.FieldSink{
				Sink:   s,
				Fields: c.fields,
			}
			fs.Log(c.msg)
			if len(s.msgs) != 1 {
				t.Fatalf("expected 1 message but got %d: %v", len(s.msgs), s.msgs)
			}
			if !reflect.DeepEqual(c.expected, s.msgs[0]) {
				t.Fatalf("expected %v but got %v", c.expected, s.msgs[0])
			}
		})
	}
}

func TestFieldSink_Multiple(t *testing.T) {
	s := &sink{}
	fs := suijin.FieldSink{
		Sink:   s,
		Fields: suijin.Fields{"sink": true},
	}

	fs.Log(suijin.Message{
		Level:   suijin.InfoLevel,
		Message: "test",
		Fields:  suijin.Fields{"a": 0},
	})

	fs.Log(suijin.Message{
		Level:   suijin.WarningLevel,
		Message: "bar",
		Fields:  suijin.Fields{"b": 1},
	})

	expected := []suijin.Message{
		{
			Level:   suijin.InfoLevel,
			Message: "test",
			Fields:  suijin.Fields{"a": 0, "sink": true},
		},
		{
			Level:   suijin.WarningLevel,
			Message: "bar",
			Fields:  suijin.Fields{"b": 1, "sink": true},
		},
	}

	if !reflect.DeepEqual(expected, s.msgs) {
		t.Errorf("expected %v but got %v", expected, s.msgs)
	}
}
