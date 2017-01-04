package attributes

// LocalException represents a local exception.
type LocalException struct {
	Message string                `json:"message,omitempty"`
	Type    string                `json:"type,omitempty"`
	Stack   []LocalExceptionStack `json:"stack,omitempty"`
}

// LocalExceptionStack represents the stack of a local exception.
type LocalExceptionStack struct {
	Path  []string `json:"path"`
	Line  int      `json:"line"`
	Label string   `json:"label"`
}

// NewLocalException creates a new local exception from an error.
func NewLocalException(err error) *LocalException {
	return &LocalException{
		Message: err.Error(),
	}
}
