package stream

import "monad/internal/monad/optional"

type Stream[T any] struct {
	elements chan T
}

func Empty[T any]() *Stream[T] {
	return &Stream[T]{}
}

func Of[T any](elements ...T) *Stream[T] {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for _, element := range elements {
			stream <- element
		}
	}()
	return &Stream[T]{elements: stream}
}

func (stream *Stream[T]) FlatMap(flatMapper func(input T) *Stream[T]) *Stream[T] {
	return FlatMap(stream, flatMapper)
}

func (stream *Stream[T]) Map(mapper func(input T) T) *Stream[T] {
	return Map(stream, mapper)
}

func (stream *Stream[T]) Filter(predicate func(input T) bool) *Stream[T] {
	return Filter(stream, predicate)
}

func (stream *Stream[T]) Skip(skip uint) *Stream[T] {
	return Skip(stream, skip)
}

func (stream *Stream[T]) Limit(limit uint) *Stream[T] {
	return Limit(stream, limit)
}

func (stream *Stream[T]) Peek(executor func(input T)) *Stream[T] {
	return Peek(stream, executor)
}

func (stream *Stream[T]) ToSlice() []T {
	return ToSlice(stream)
}

func (stream *Stream[T]) ForEach(executor func(input T)) {
	ForEach(stream, executor)
}

func (stream *Stream[T]) Reduce(reducer func(accumulator T, next T) T) *optional.Optional[T] {
	return Reduce(stream, reducer)
}

func (stream *Stream[T]) ReduceWithIdentity(identity T, reducer func(accumulator T, next T) T) T {
	return ReduceWithIdentity(stream, identity, reducer)
}

func (stream *Stream[T]) FindFirst(predicate func(input T) bool) *optional.Optional[T] {
	return FindFirst(stream, predicate)
}

func (stream *Stream[T]) AnyMatch(predicate func(input T) bool) bool {
	return AnyMatch(stream, predicate)
}

func (stream *Stream[T]) AllMatch(predicate func(input T) bool) bool {
	return AllMatch(stream, predicate)
}

func FlatMap[I, O any](inputStream *Stream[I], flatMapper func(input I) *Stream[O]) *Stream[O] {
	outputChannel := make(chan O)
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			for subElement := range flatMapper(element).elements {
				outputChannel <- subElement
			}
		}
	}()
	return &Stream[O]{elements: outputChannel}
}

func Map[I, O any](inputStream *Stream[I], mapper func(input I) O) *Stream[O] {
	outputChannel := make(chan O)
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			outputChannel <- mapper(element)
		}
	}()
	return &Stream[O]{elements: outputChannel}
}

func Filter[T any](inputStream *Stream[T], predicate func(input T) bool) *Stream[T] {
	outputChannel := make(chan T)
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			if predicate(element) {
				outputChannel <- element
			}
		}
	}()
	return &Stream[T]{elements: outputChannel}
}

func Skip[T any](inputStream *Stream[T], skip uint) *Stream[T] {
	outputChannel := make(chan T)
	var index uint = 0
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			if index >= skip {
				outputChannel <- element
			}
			index++
		}
	}()
	return &Stream[T]{elements: outputChannel}
}

func Limit[T any](inputStream *Stream[T], limit uint) *Stream[T] {
	outputChannel := make(chan T)
	var index uint = 0
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			if index < limit {
				outputChannel <- element
			}
			index++
		}
	}()
	return &Stream[T]{elements: outputChannel}
}

func Peek[T any](inputStream *Stream[T], executor func(input T)) *Stream[T] {
	outputChannel := make(chan T)
	go func() {
		defer close(outputChannel)
		for element := range inputStream.elements {
			executor(element)
			outputChannel <- element
		}
	}()
	return &Stream[T]{elements: outputChannel}
}

func ToSlice[T any](inputStream *Stream[T]) []T {
	var output []T
	for element := range inputStream.elements {
		output = append(output, element)
	}
	return output
}

func toMap[T any, K comparable, V any](
	inputStream *Stream[T],
	keyMapper func(input T) K,
	valueMapper func(input T) V,
) map[K]V {
	var output map[K]V
	for element := range inputStream.elements {
		output[keyMapper(element)] = valueMapper(element)
	}
	return output
}

func ForEach[T any](inputStream *Stream[T], executor func(input T)) {
	for element := range inputStream.elements {
		executor(element)
	}
}

func Reduce[I, O any](inputStream *Stream[I], reducer func(accumulator O, next I) O) *optional.Optional[O] {
	var output O
	for element := range inputStream.elements {
		output = reducer(output, element)
	}
	return optional.Of(output)
}

func ReduceWithIdentity[I, O any](inputStream *Stream[I], identity O, reducer func(accumulator O, next I) O) O {
	output := identity
	for element := range inputStream.elements {
		output = reducer(output, element)
	}
	return output
}

func FindFirst[T any](inputStream *Stream[T], predicate func(input T) bool) *optional.Optional[T] {
	var output T
	for output == nil {
		select {
		case value, ok := <-inputStream.elements:
			if !ok {
				break
			}
			if predicate(value) {
				output = value
			}
		}
	}
	return optional.Of(output)
}

func AnyMatch[T any](inputStream *Stream[T], predicate func(input T) bool) bool {
	output := false
	for !output {
		select {
		case value, ok := <-inputStream.elements:
			if !ok {
				break
			}
			if predicate(value) {
				output = true
			}
		}
	}
	return output
}

func AllMatch[T any](inputStream *Stream[T], predicate func(input T) bool) bool {
	output := true
	for output {
		select {
		case value, ok := <-inputStream.elements:
			if !ok {
				break
			}
			if !predicate(value) {
				output = false
			}
		}
	}
	return output
}

//func fromInputToOutput[I, O any](
//	inputStream *Stream[I],
//	fromInputChannelToOutputChannel func(input chan I, output chan O),
//) *Stream[O] {
//	outputChannel := make(chan O)
//	go func() {
//		defer close(outputChannel)
//		fromInputChannelToOutputChannel(inputStream.elements, outputChannel)
//	}()
//	return &Stream[O]{elements: outputChannel}
//}
