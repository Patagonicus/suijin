package suijin

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"sync"
)

// Sink consumes log messages. If you are writing a library that does logging you should accept a Sink
// and then wrap it in a Logger yourself.
type Sink interface {
	// Log logs a message. The message's level must be a valid, non-special log level and it's Fields attribute must not
	// be null. If either or both of these conditions is not met this method may drop the message silently.
	Log(msg Message)
}

// WriterSink writes incoming log messages to an io.Writer. It can be use concurrently.
type WriterSink struct {
	w io.Writer
	l sync.Locker
}

// Writer returns a Sink that writes incoming messages to the given Writer.
func Writer(w io.Writer) *WriterSink {
	return &WriterSink{
		w: w,
		l: &sync.Mutex{},
	}
}

// Log writes the given log message to the Writer.
func (s *WriterSink) Log(msg Message) {
	formatted := formatMessage(msg)

	s.l.Lock()
	defer s.l.Unlock()

	// we do not check for write errors here because we can't do anything about them
	// (other than panic-ing)
	io.WriteString(s.w, formatted) // nolint: errcheck
}

func formatMessage(msg Message) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s: %s", msg.Level, msg.Message)
	formatFields(buf, msg.Fields)
	fmt.Fprint(buf, "\n")
	return buf.String()
}

func formatFields(buf *bytes.Buffer, fields Fields) string {
	// formatFields takes *bytes.Buffer instead of io.Writer because it doesn't do any error checking.

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(buf, " '%s'='%v'", k, fields[k])
	}

	return buf.String()
}

// FieldSink is a sink that adds fields to a message before forwarding it to another sink.
// Its elements may not be modified after Log has been called for the first time.
type FieldSink struct {
	Sink   Sink
	Fields Fields
}

// Log adds fields to the given message before forwarding it.
func (f FieldSink) Log(msg Message) {
	if msg.Fields == nil && len(f.Fields) != 0 {
		msg.Fields = make(Fields)
	}
	msg.Fields.AddAll(f.Fields)
	f.Sink.Log(msg)
}
