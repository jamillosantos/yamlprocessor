package yamlprocessor

import (
	"errors"
)

var ErrEntryAlreadyRegistered = errors.New("function already registered")

// ExprEnv is a map of functions that can be registered to the processor.
type ExprEnv map[string]ExprEntry

// Add adds a function to the ExprEnv map making it available to be acessible for the Processor. If the map is
// not initialised, it will be initialised.
func (f *ExprEnv) Add(name string, function ExprEntry) error {
	if f == nil {
		*f = make(map[string]ExprEntry)
	}
	if _, ok := (*f)[name]; ok {
		return ErrEntryAlreadyRegistered
	}
	(*f)[name] = function
	return nil
}

type Context struct {
	Processor *Processor
}
