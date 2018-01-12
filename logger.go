package suijin

// Logger is a convenience wrapper around a Sink. Logger has seperate methods for logging at the different levels
// directly. It also allows you to add fields to a log message.
//
// When logging you should use a Logger, but libraries should only depend on a Sink to make it easier to substitute
// their dependency.
type Logger struct {
	s FieldSink
}

// NewLogger wraps a Sink in a Logger.
func NewLogger(s Sink) Logger {
	var fs FieldSink
	fs, ok := s.(FieldSink)
	if !ok {
		fs = FieldSink{
			Sink:   s,
			Fields: make(Fields),
		}
	}

	return Logger{
		s: fs,
	}
}

// Debug logs a message that is only inteded for verbose logs.
func (l Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

// Info logs a message that should be printed during normal execution.
func (l Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

// Warning logs a message for an event that could lead to a problem, but that the program can recover from.
func (l Logger) Warning(msg string) {
	l.log(WarningLevel, msg)
}

// Error logs a message for something that the program cannot recover from.
func (l Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l Logger) log(lvl Level, msg string) {
	l.s.Log(Message{
		Level:   lvl,
		Message: msg,
		Fields:  make(Fields),
	})
}

// WithField returns a new Logger that will automatically add the given field to all messages logged with it.
func (l Logger) WithField(key string, value interface{}) Logger {
	return l.WithFields(Fields{key: value})
}

// WithFields returns a new Logger that will automatically add the given fields to all messages logged with it.
func (l Logger) WithFields(fds Fields) Logger {
	newFds := make(Fields)
	newFds.AddAll(l.s.Fields)
	newFds.AddAll(fds)
	return Logger{
		s: FieldSink{
			Sink:   l.s.Sink,
			Fields: newFds,
		},
	}
}
