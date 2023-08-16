package errordefs

import "fmt"

// Define custom error types

type newComposerError struct {
	message string
}

func (e newComposerError) Error() string {
	return fmt.Sprintf("ComposerError: %s", e.message)
}

type composerRemoveError struct {
	message string
}

func (e composerRemoveError) Error() string {
	return fmt.Sprintf("ComposerRemoveError: %s", e.message)
}

type composerStopError struct {
	message string
}

func (e composerStopError) Error() string {
	return fmt.Sprintf("ComposerStopError: %s", e.message)
}

type composerUpError struct {
	message string
}

func (e composerUpError) Error() string {
	return fmt.Sprintf("ComposerUpError: %s", e.message)
}

// Generic Composer Errors
func NewComposerError(err error) error {
	return &newComposerError{message: err.Error()}
}

// Errors that occur when invoking remove function
func NewComposerRemoveError(err error) error {
	return &composerRemoveError{message: err.Error()}
}

// Errors that occur when invoking stop function
func NewComposerStopError(err error) error {
	return &composerStopError{message: err.Error()}
}

// Errors that occur when invoking Up function
func NewComposerUpError(err error) error {
	return &composerUpError{message: err.Error()}
}
