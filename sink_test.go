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
