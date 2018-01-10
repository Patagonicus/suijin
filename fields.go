package suijin

// Fields stores additional information for a log message.
type Fields map[string]interface{}

// AddAll adds everything from other to this Fields object. Existing entries are overwritten if they use the same key.
func (f Fields) AddAll(other Fields) {
	for k, v := range other {
		f[k] = v
	}
}
