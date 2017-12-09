package suijin_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Patagonicus/suijin"
	"github.com/Patagonicus/suijin/legacy"
)

const goroutineCount = 20
const iterationCount = 1000

type mockFormatter struct{}

func (f mockFormatter) Format(m suijin.Message) ([]byte, error) {
	return []byte(m.Msg), nil
}

func TestBackend(t *testing.T) {
	buf := new(bytes.Buffer)
	b := suijin.NewWriterBackend(buf, mockFormatter{})
	wg := new(sync.WaitGroup)
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(index int) {
			for i := 0; i < iterationCount; i++ {
				b.Log(suijin.Message{
					Msg: fmt.Sprintf("%d %d", index, i),
				})
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	if buf.Bytes()[buf.Len()-1] != '\n' {
		t.Fatalf("buffer does not end with a newline")
	}

	buf.Truncate(buf.Len() - 1)
	lines, err := extractMessages(buf.String())
	if err != nil {
		t.Fatal(err)
	}

	legacy.Slice(lines, func(a, b int) bool {
		if lines[a].goroutine == lines[b].goroutine {
			return lines[a].index < lines[b].index
		}
		return lines[a].goroutine < lines[b].goroutine
	})

	for i := 0; i < goroutineCount; i++ {
		for j := 0; j < iterationCount; j++ {
			index := i*iterationCount + j
			expected := line{i, j}
			if lines[index] != expected {
				t.Fatalf("line %d is wrong, expected '%v' but got '%v'", index, expected, lines[index])
			}
		}
	}
}

type line struct {
	goroutine int
	index     int
}

func extractMessages(log string) ([]line, error) {
	var lines []line
	for _, l := range strings.Split(log, "\n") {
		s := strings.Split(l, " ")
		if len(s) != 2 {
			return nil, fmt.Errorf("invalid line: %s", l)
		}

		r, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid line: %s", l)
		}

		i, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid line: %s", l)
		}

		lines = append(lines, line{r, i})
	}
	return lines, nil
}
