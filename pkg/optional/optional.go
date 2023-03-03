package optional

import (
	"fmt"

	"github.comn/SnowPhoenix0105/smartpacker/pkg/utils"
)

type Of[T any] struct {
	Ptr *T
}

func FromPtr[T any](ptr *T) Of[T] {
	return Of[T]{Ptr: ptr}
}

func OfValue[T any](val T) Of[T] {
	return FromPtr(&val)
}

func OfNil[T any]() Of[T] {
	return FromPtr[T](nil)
}

func (r Of[T]) Unwrap() (T, bool) {
	if r.Ptr == nil {
		return utils.ZeroOf[T](), false
	}
	return *r.Ptr, true
}

func (r Of[T]) TryAssignTo(target *T) {
	if val, ok := r.Unwrap(); ok {
		target = &val
	}
}

func (r Of[T]) Set(val T) {
	r.Ptr = &val
}

func (r Of[T]) Del() {
	r.Ptr = nil
}

func (r Of[T]) Ok() bool {
	return r.Ptr != nil
}

func (r Of[T]) Value() T {
	val, _ := r.Unwrap()
	return val
}

func (r Of[T]) ValueOr(def T) T {
	if r.Ptr != nil {
		return *r.Ptr
	}
	return def
}

func (r Of[T]) AsPtr() *T {
	return r.Ptr
}

func (r Of[T]) AsInterface() any {
	if r.Ptr == nil {
		return nil
	}
	return *r.Ptr
}

func (r Of[T]) Interface() any {
	return r.Value()
}

func (r Of[T]) String() string {
	val, ok := r.Unwrap()
	if ok {
		return fmt.Sprintf("Ok(%v)", val)
	}
	return fmt.Sprintf("Nil(%v)", val)
}
