package result

import (
	"fmt"
)

// Result is a type that represents either success (T) or failure (error).
type Result[T any] struct {
	ok  T
	err error
}

func Wrap[T any](some T, err error) *Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(some)
}

func Ok[T any](ok T) *Result[T] {
	return &Result[T]{ok: ok}
}

func Err[T any](err any) *Result[T] {
	return &Result[T]{err: newAnyError(err)}
}

// IsOk returns true if the result is Ok.
func (r *Result[T]) IsOk() bool {
	return !r.IsErr()
}

// IsOkAnd returns true if the result is Ok and the value inside of it matches a predicate.
func (r *Result[T]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.ok)
	}
	return false
}

// IsErr returns true if the result is error.
func (r *Result[T]) IsErr() bool {
	return r.err != nil
}

// IsErrAnd returns true if the result is Err and the value inside of it matches a predicate.
func (r *Result[T]) IsErrAnd(f func(error) bool) bool {
	if r.IsErr() {
		return f(r.err)
	}
	return false
}

// Ok returns T, and returns empty if it is an error.
func (r *Result[T]) Ok() *T {
	if r.IsOk() {
		return &r.ok
	}
	return nil
}

// Err returns error.
func (r *Result[T]) Err() error {
	return r.err
}

// ErrVal returns error inner value.
func (r *Result[T]) ErrVal() any {
	if r.IsErr() {
		return nil
	}
	if ev, _ := r.err.(*errorWithVal); ev != nil {
		return ev.val
	}
	return r.err
}

// Map maps a Result[T] to Result[T] by applying a function to a contained Ok value, leaving an Err value untouched.
// This function can be used to compose the results of two functions.
func (r *Result[T]) Map(f func(T) T) *Result[T] {
	if r.IsOk() {
		return Ok[T](f(r.ok))
	}
	return Err[T](r.err)
}

// Map maps a Result[T] to Result[U] by applying a function to a contained Ok value, leaving an Err value untouched.
// This function can be used to compose the results of two functions.
func Map[T any, U any](r *Result[T], f func(T) U) *Result[U] {
	if r.IsOk() {
		return Ok[U](f(r.ok))
	}
	return Err[U](r.err)
}

// MapOr returns the provided default (if Err), or applies a function to the contained value (if Ok),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use map_or_else, which is lazily evaluated.
func (r *Result[T]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultOk
}

// MapOr returns the provided default (if Err), or applies a function to the contained value (if Ok),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use map_or_else, which is lazily evaluated.
func MapOr[T any, U any](r *Result[T], defaultOk U, f func(T) U) U {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultOk
}

// MapOrElse maps a Result[T] to T by applying fallback function default to a contained Err value, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r *Result[T]) MapOrElse(defaultFn func(error) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultFn(r.err)
}

// MapOrElse maps a Result[T] to U by applying fallback function default to a contained Err value, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func MapOrElse[T any, U any](r *Result[T], defaultFn func(error) U, f func(T) U) U {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultFn(r.err)
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained Err value, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r *Result[T]) MapErr(op func(error) error) *Result[T] {
	if r.IsErr() {
		r.err = op(r.err)
	}
	return r
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained Err value, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
func MapErr[T any](r *Result[T], op func(error) error) *Result[T] {
	if r.IsErr() {
		r.err = op(r.err)
	}
	return r
}

// Inspect calls the provided closure with a reference to the contained value (if Ok).
func (r *Result[T]) Inspect(f func(T)) *Result[T] {
	if r.IsOk() {
		f(r.ok)
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained error (if Err).
func (r *Result[T]) InspectErr(f func(error)) *Result[T] {
	if r.IsErr() {
		f(r.err)
	}
	return r
}

// Expect returns the contained Ok value, consuming the self value.
func (r *Result[T]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Errorf("%s: %w", msg, r.err))
	}
	return r.ok
}

// Unwrap returns the contained Ok value, consuming the self value.
// Because this function may panic, its use is generally discouraged. Instead, prefer to use pattern matching and handle the Err case explicitly, or call unwrap_or, unwrap_or_else, or unwrap_or_default.
func (r *Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Errorf("called `Result.unwrap()` on an `err` value: %w", r.err))
	}
	return r.ok
}

// UnwrapOrDefault returns the contained Ok value or a default
// Consumes the self argument then, if Ok, returns the contained value, otherwise if Err, returns the default value for that type.
func (r *Result[T]) UnwrapOrDefault() T {
	return r.ok
}

func (r *Result[T]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%s)", r.err.Error())
	}
	return fmt.Sprintf("Ok(%v)", r.ok)
}
