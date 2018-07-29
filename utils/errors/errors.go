package errors

import (
	"sync"
)

type MultiError struct {
	Errors []error
	mux sync.Mutex
}

func (err *MultiError) Error() string {
	str := ""

	for _, e := range err.Errors {
		str += e.Error()
	}

	return str
}

func (err *MultiError) Add(e error) {
	if e == nil {
		return
	}

	err.mux.Lock()
	err.Errors = append(err.Errors, e)
	err.mux.Unlock()
}

func (err *MultiError) Return() error {
	if len(err.Errors) > 0 {
		return err
	}

	return nil
}
