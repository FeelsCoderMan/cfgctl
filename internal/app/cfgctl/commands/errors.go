package commands

import "fmt"

type ErrorKind string

const (
	KindUnknown         ErrorKind = "unknown"
	KindInvalidJSON     ErrorKind = "invalid_json"
	KindMissingPath     ErrorKind = "missing_path"
	KindMissingArgs     ErrorKind = "missing_args"
	KindTooManyArgs     ErrorKind = "too_many_args"
	KindFileStoreGet    ErrorKind = "file_store_get"
	KindFileStoreLoad   ErrorKind = "file_store_load"
	KindFileStoreDelete ErrorKind = "file_store_delete"
	KindFileStoreSet    ErrorKind = "file_store_set"
)

type CommandError struct {
	Kind  ErrorKind
	Cause error
}

func (e *CommandError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Kind, e.Cause)
	}
	return string(e.Kind)
}

func (e *CommandError) Unwrap() error {
	return e.Cause
}

func NewCommandError(kind ErrorKind, cause error) *CommandError {
	return &CommandError{
		Kind:  kind,
		Cause: cause,
	}
}
