package cmd

import "fmt"

// Define custom error types

type removeError struct {
	message string
}

func (e removeError) Error() string {
	return fmt.Sprintf("RemoveError: %s", e.message)
}

type startError struct {
	message string
}

func (e startError) Error() string {
	return fmt.Sprintf("StartError: %s", e.message)
}

type stopError struct {
	message string
}

func (e stopError) Error() string {
	return fmt.Sprintf("StopError: %s", e.message)
}

type setupError struct {
	message string
}

func (e setupError) Error() string {
	return fmt.Sprintf("SetupError: %s", e.message)
}

func NewSetupError(err error) error {
	return &setupError{message: err.Error()}
}

func NewRemoveError(err error) error {
	return &removeError{message: err.Error()}
}

func NewStartError(err error) error {
	return &startError{message: err.Error()}
}

func NewStopError(err error) error {
	return &stopError{message: err.Error()}
}
