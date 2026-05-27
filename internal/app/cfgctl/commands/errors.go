package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ErrorKind string

type CommandKind string

const (
	CommandUnknown CommandKind = "unknown"
	CommandInit    CommandKind = "init"
	CommandGet     CommandKind = "get"
	CommandDelete  CommandKind = "delete"
	CommandSet     CommandKind = "set"
	CommandList    CommandKind = "list"
)

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
	ErrorKind    ErrorKind
	commandKind  CommandKind
	errorMessage string
	Detail       error
}

func (e *CommandError) Error() string {
	return e.ErrorMessage()
}

func (e *CommandError) Unwrap() error {
	return e.Detail
}

func (e *CommandError) DetailedErrorMessage() string {
	if e.errorMessage != "" {
		return e.errorMessage
	}

	seen := map[string]struct{}{}
	errorLines := []string{e.errorMessage}
	for err := e.Detail; err != nil; err = errors.Unwrap(err) {
		s := err.Error()
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		errorLines = append(errorLines, s)
		seen[s] = struct{}{}
	}

	return strings.Join(errorLines, "\n")
}

func (e *CommandError) ErrorMessage() string {
	if e.errorMessage != "" {
		return e.errorMessage
	}

	kindMessage := ""

	switch e.ErrorKind {
	case KindUnknown:
		kindMessage = "unknown error"
	case KindInvalidJSON:
		kindMessage = "invalid JSON"
	case KindMissingPath:
		kindMessage = "missing path"
	case KindMissingArgs:
		kindMessage = "missing arguments"
	case KindTooManyArgs:
		kindMessage = "too many arguments"
	case KindFileStoreGet:
		kindMessage = "file store get error"
	case KindFileStoreLoad:
		kindMessage = "file store load error"
		var pathError *os.PathError
		var unmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(e.Detail, &pathError):
			kindMessage = "file not found"
		case errors.As(e.Detail, &unmarshalError):
			kindMessage = "file contains invalid JSON"
		}
	case KindFileStoreDelete:
		kindMessage = "file store delete error"
	case KindFileStoreSet:
		kindMessage = "file store set error"
	}

	e.errorMessage = fmt.Sprintf("[%s] %s", e.commandKind, kindMessage)
	return e.errorMessage
}

func NewCommandError(errorKind ErrorKind, commandKind CommandKind, detail error) *CommandError {
	if errorKind == "" {
		errorKind = KindUnknown
	}

	if commandKind == "" {
		commandKind = CommandUnknown
	}

	return &CommandError{
		ErrorKind:   errorKind,
		commandKind: commandKind,
		Detail:      detail,
	}
}
