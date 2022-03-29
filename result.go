package result

import (
	"fmt"
)

// Result is a type that represents either success (Ok) or failure (Err).
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

func (r Result[T]) IsOk() bool {
	return !r.IsErr()
}

func (r Result[T]) Ok() T {
	return r.ok
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Err() error {
	return r.err
}

func (r Result[T]) ErrVal() any {
	if r.IsErr() {
		return nil
	}
	if ev, _ := r.err.(*errorWithVal); ev != nil {
		return ev.val
	}
	return r.err.Error()
}

func (r Result[T]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%s)", r.err.Error())
	}
	return fmt.Sprintf("Ok(%v)", r.ok)
}

func Map[T any, B any](r Result[T], op func(T) B) Result[B] {
	if r.IsOk() {
		return Ok[B](op(r.ok))
	}
	return Err[B](r.err)
}
