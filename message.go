package suijin

// Message describes a single event that has been logged.
type Message struct {
	// The level associated with this message. Sinks can use this to decicede if and how to log this message.
	Level Level
	// A text describing the event. This should be a static text, all information that is not known at compile time
	// should be put into Fields.
	Message string
	// A map containing additional information. This could be the result of a calculation, the name of the request
	// currently being handled or an error that happend.
	Fields Fields
}
