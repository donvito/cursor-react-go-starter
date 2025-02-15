package repository

import "fmt"

// ErrNotFound is returned when a todo item is not found
type ErrNotFound struct {
	ID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("todo with id %s not found", e.ID)
}

// ErrInvalidInput is returned when the input is invalid
type ErrInvalidInput struct {
	Message string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Message)
}

// ErrDatabase is returned when a database error occurs
type ErrDatabase struct {
	Op  string
	Err error
}

func (e *ErrDatabase) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Op, e.Err)
}

func (e *ErrDatabase) Unwrap() error {
	return e.Err
}
