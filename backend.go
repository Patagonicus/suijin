package suijin

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

// Backend is something that can store messages.
type Backend interface {
	// Writes a message to the backend. This must be safe to use concurrently.
	Log(m Message)
}

// DefaultBackend writes messages to stderr.
var DefaultBackend = NewWriterBackend(os.Stderr, DefaultFormatter)

type writerBackend struct {
	formatter Formatter
	writer    io.Writer
	lock      *sync.Mutex
}

// NewWriterBackend returns a Backend that writes messages to an io.Writer.
func NewWriterBackend(w io.Writer, f Formatter) Backend {
	return writerBackend{
		f,
		w,
		new(sync.Mutex),
	}
}

func (w writerBackend) Log(m Message) {
	b, err := w.formatter.Format(m)
	if err != nil {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "error formatting log message: %s", err)
		b = buf.Bytes()
	}

	w.lock.Lock()
	defer w.lock.Unlock()
	w.writer.Write(b)
	w.writer.Write([]byte{'\n'})
}
