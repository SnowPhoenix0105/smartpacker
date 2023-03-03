package result

import (
	"fmt"

	"github.comn/SnowPhoenix0105/smartpacker/pkg/utils"
)

type Of[T any] struct {
	Val T
	Err error
}

func create[T any](val T, err error) Of[T] {
	return Of[T]{Val: val, Err: err}
}

func createPtr[T any](val T, err error) *Of[T] {
	res := create(val, err)
	return &res
}

func OfValue[T any](val T) Of[T] {
	return create(val, nil)
}

func OfError[T any](err error) Of[T] {
	return create(utils.ZeroOf[T](), err)
}

func PtrOfValue[T any](val T) *Of[T] {
	return createPtr(val, nil)
}

func PtrOfError[T any](err error) *Of[T] {
	return createPtr(utils.ZeroOf[T](), nil)
}

func (r Of[T]) Unwrap() (T, error) {
	return r.Val, r.Err
}

func (r Of[T]) TryAssignTo(target *T) {
	if r.Err != nil {
		return
	}
	*target = r.Val
}

func (r Of[T]) Error() error {
	return r.Err
}

func (r Of[T]) Value() T {
	return r.Val
}

func (r Of[T]) Interface() any {
	return r.Val
}

func (r Of[T]) ValueOr(def T) T {
	if r.Err == nil {
		return r.Val
	}
	return def
}

func (r Of[T]) String() string {
	val, err := r.Unwrap()
	if err != nil {
		return fmt.Sprintf("Error(%v)", err)
	}
	return fmt.Sprintf("Success(%v)", val)
}
