package errordefs

import "fmt"

// Define custom error types

type newConduitFromProjectError struct {
	message string
}

func (e newConduitFromProjectError) Error() string {
	return fmt.Sprintf("NewConduitProjectError: %s", e.message)
}

type envFileError struct {
	message string
}

func (e envFileError) Error() string {
	return fmt.Sprintf("EnvFileError: %s", e.message)
}

type yamlFileError struct {
	message string
}

func (e yamlFileError) Error() string {
	return fmt.Sprintf("YamlFileError: %s", e.message)
}

type newConduitBootstrapperError struct {
	message string
}

func (e newConduitBootstrapperError) Error() string {
	return fmt.Sprintf("NewConduitBootstrapperError: %s", e.message)
}

// Errors that occur when bootstraping a new conduit project
func NewConduitBootstrapperError(err error) error {
	return &newConduitBootstrapperError{message: err.Error()}
}

// Errors that occur while in conduit project directory
func NewNewConduitFromProjectError(err error) error {
	return &newConduitFromProjectError{message: err.Error()}
}

// Errors env file related
func NewEnvFileError(err error) error {
	return &envFileError{message: err.Error()}
}

// Errors yaml file related
func NewYamlFileError(err error) error {
	return &yamlFileError{message: err.Error()}
}
