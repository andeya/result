package result

import (
	"fmt"
)

// Result is a type that represents either success (T) or failure (error).
type Result[T any] struct {
	ok  T
	err error
}

func Wrap[T any](some T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(some)
}

func Ok[T any](ok T) Result[T] {
	return Result[T]{ok: ok}
}

func Err[T any](err any) Result[T] {
	return Result[T]{err: newAnyError(err)}
}

// IsOk returns true if the result is Ok.
func (r Result[T]) IsOk() bool {
	return !r.IsErr()
}

// IsOkAnd returns true if the result is Ok and the value inside of it matches a predicate.
func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.ok)
	}
	return false
}

// IsErr returns true if the result is error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// IsErrAnd returns true if the result is Err and the value inside of it matches a predicate.
func (r Result[T]) IsErrAnd(f func(error) bool) bool {
	if r.IsErr() {
		return f(r.err)
	}
	return false
}

// Ok returns T, and returns empty if it is an error.
func (r Result[T]) Ok() *T {
	if r.IsOk() {
		return &r.ok
	}
	return nil
}

// Err returns error.
func (r Result[T]) Err() error {
	return r.err
}

// ErrVal returns error inner value.
func (r Result[T]) ErrVal() any {
	if r.IsErr() {
		return nil
	}
	if ev, _ := r.err.(*errorWithVal); ev != nil {
		return ev.val
	}
	return r.err
}

// Map maps a Result[T] to Result[U] by applying a function to a contained Ok value, leaving an Err value untouched.
// This function can be used to compose the results of two functions.
func Map[T any, U any](r Result[T], op func(T) U) Result[U] {
	if r.IsOk() {
		return Ok[U](op(r.ok))
	}
	return Err[U](r.err)
}

// MapOr returns the provided default (if Err), or applies a function to the contained value (if Ok),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use map_or_else, which is lazily evaluated.
func MapOr[T any, U any](r Result[T], defaultOk U, op func(T) U) U {
	if r.IsOk() {
		return op(r.ok)
	}
	return defaultOk
}

func (r Result[T]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%s)", r.err.Error())
	}
	return fmt.Sprintf("Ok(%v)", r.ok)
}
