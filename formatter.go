package suijin

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"time"
)

// Formatter formats a log message.
type Formatter interface {
	// Format takes a message and turns it into a []byte representation.
	Format(m Message) ([]byte, error)
}

type FormatterConfig struct {
	// The clock to use. Defaults to SystemClock if nil. You should probably not change this.
	Clock Clock
}

type textFormatter struct {
	c Clock
}

var DefaultFormatter = NewTextFormatter(FormatterConfig{})

func NewTextFormatter(config FormatterConfig) Formatter {
	f := textFormatter{}

	if config.Clock == nil {
		f.c = SystemClock
	} else {
		f.c = config.Clock
	}

	return f
}

func (f textFormatter) Format(m Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	fmt.Fprintf(w, "%s %s: %s", f.c.Now().Format(time.RFC3339), m.Lvl, m.Msg)

	keys := make([]string, 0, len(m.Fds))
	for k := range m.Fds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, " %q=%q", k, m.Fds[k])
	}

	err := w.Flush()

	return buf.Bytes(), err
}
