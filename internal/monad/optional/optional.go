package optional

import (
	"errors"
	"monad/internal/monad/stream"
)

type Optional[T any] struct {
	value T
}

func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

func Of[T any](value T) *Optional[T] {
	return &Optional[T]{value: value}
}

func (optional *Optional[T]) FlatMap(flatMapper func(T) *Optional[T]) *Optional[T] {
	return FlatMap(optional, flatMapper)
}

func (optional *Optional[T]) Map(mapper func(T) T) *Optional[T] {
	return Map(optional, mapper)
}

func (optional *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	return Filter(optional, predicate)
}

func (optional *Optional[T]) Peek(executor func(T)) *Optional[T] {
	return Peek(optional, executor)
}

func (optional *Optional[T]) Stream() *stream.Stream[T] {
	return Stream(optional)
}

func (optional *Optional[T]) Or(optionalSupplier func() *Optional[T]) *Optional[T] {
	return Or(optional, optionalSupplier)
}

func (optional *Optional[T]) OrElse(defaultValue T) T {
	return OrElse(optional, defaultValue)
}

func (optional *Optional[T]) OrElseGet(defaultValueSupplier func() T) T {
	return OrElseGet(optional, defaultValueSupplier)
}

func (optional *Optional[T]) Get() (T, error) {
	return Get(optional)
}

func (optional *Optional[T]) IfPresent(executor func(input T)) {
	IfPresent(optional, executor)
}

func (optional *Optional[T]) OrElseFallback(fallback func()) {
	OrElseFallback(optional, fallback)
}

func (optional *Optional[T]) IfPresentOrElseFallback(executor func(input T), fallback func()) {
	IfPresentOrElseFallback(optional, executor, fallback)
}

func (optional *Optional[T]) IsPresent() bool {
	return IsPresent(optional)
}

func (optional *Optional[T]) IsEmpty() bool {
	return IsEmpty(optional)
}

func FlatMap[I, O any](optional *Optional[I], flatMapper func(input I) *Optional[O]) *Optional[O] {
	if optional.IsEmpty() {
		return Empty[O]()
	}
	return flatMapper(optional.value)
}

func Map[I, O any](optional *Optional[I], mapper func(input I) O) *Optional[O] {
	if optional.IsEmpty() {
		return Empty[O]()
	}
	return Of(mapper(optional.value))
}

func Filter[T any](optional *Optional[T], predicate func(input T) bool) *Optional[T] {
	if optional.IsEmpty() || !predicate(optional.value) {
		return Empty[T]()
	}
	return Of(optional.value)
}

func Peek[T any](optional *Optional[T], executor func(input T)) *Optional[T] {
	if optional.IsPresent() {
		executor(optional.value)
	}
	return Of(optional.value)
}

func Stream[T any](optional *Optional[T]) *stream.Stream[T] {
	if optional.IsEmpty() {
		return stream.Empty[T]()
	}
	return stream.Of(optional.value)
}

func Or[T any](optional *Optional[T], optionalSupplier func() *Optional[T]) *Optional[T] {
	if optional.IsEmpty() {
		return optionalSupplier()
	}
	return Of(optional.value)
}

func OrElse[T any](optional *Optional[T], defaultValue T) T {
	if optional.IsEmpty() {
		return defaultValue
	}
	return optional.value
}

func OrElseGet[T any](optional *Optional[T], defaultValueSupplier func() T) T {
	if optional.IsEmpty() {
		return defaultValueSupplier()
	}
	return optional.value
}

func Get[T any](optional *Optional[T]) (T, error) {
	if optional.IsEmpty() {
		return *new(T), errors.New("cannot get value from empty optional")
	}
	return optional.value, nil
}

func IfPresent[T any](optional *Optional[T], executor func(input T)) {
	if optional.IsPresent() {
		executor(optional.value)
	}
}

func OrElseFallback[T any](optional *Optional[T], fallback func()) {
	if optional.IsEmpty() {
		fallback()
	}
}

func IfPresentOrElseFallback[T any](optional *Optional[T], executor func(input T), fallback func()) {
	if optional.IsPresent() {
		executor(optional.value)
	} else {
		fallback()
	}
}

func IsPresent[T any](optional *Optional[T]) bool {
	return optional.value != nil
}

func IsEmpty[T any](optional *Optional[T]) bool {
	return optional.value == nil
}
