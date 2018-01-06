package suijin

import "fmt"

// Level describes the severity of a log message.
type Level uint

const (
	// LogAll is not a valid log level but can be used with certain filters to keep/discard all messages.
	LogAll Level = iota
	// DebugLevel is for messages that should not be printed during regular execution.
	DebugLevel
	// InfoLevel is for messages that show normal execution. If you are writing a server these should only
	// be printed on startup (telling the user that the server has successfully started and logging information
	// like the servers address). For tools that are invoked by a user directly these can also be used to show
	// progress during long operations.
	InfoLevel
	// WarningLevel should be used for events that show a potential future problem such as a deprecated configuration
	// option being used or (repeated) invalid logins.
	WarningLevel
	// ErrorLevel should be used for problems that should be looked at as soon as possible.
	ErrorLevel
	// LogNone is not a valid log level but can be used with certain filters to keep/discard no messages.
	LogNone
)

// LevelFromString turns a string into a Level. This works only for these strings:
//
//   * "all"
//   * "debug"
//   * "info"
//   * "warning"
//   * "error"
//   * "none"
//
// All other strings will return an error instead.
func LevelFromString(s string) (Level, error) {
	switch s {
	case "all":
		return LogAll, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warning":
		return WarningLevel, nil
	case "error":
		return ErrorLevel, nil
	case "none":
		return LogNone, nil
	default:
		return LogNone, fmt.Errorf("invalid log level: %v", s)
	}
}

func (l Level) String() string {
	switch l {
	case LogAll:
		return "all"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case LogNone:
		return "none"
	default:
		return fmt.Sprintf("unknown log level %d", l)
	}
}

// IsValid returns true if the level is valid. A level is valid if it is one of the values defined in this package.
func (l Level) IsValid() bool {
	return LogAll <= l && l <= LogNone
}

// IsSpecial returns true if the level is one of the special levels. Currently these are LogAll and LogNone.
// Special levels cannot be used in Messages.
func (l Level) IsSpecial() bool {
	return l == LogAll || l == LogNone
}
