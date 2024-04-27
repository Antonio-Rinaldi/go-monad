package function

type Predicate[T any] func(x T) bool
type Function[I, O any] func(x I) O

func Identity[T any]() Function[T, T] {
	return func(x T) T {
		return x
	}
}

func True[T any]() Predicate[T] {
	return func(x T) bool {
		return true
	}
}

func False[T any]() Predicate[T] {
	return func(x T) bool {
		return false
	}
}

func Compose[I any, M any, O any](f Function[M, O], g Function[I, M]) Function[I, O] {
	return func(x I) O {
		return f(g(x))
	}
}

func AndThen[I any, M any, O any](f Function[I, M], g Function[M, O]) Function[I, O] {
	return func(x I) O {
		return g(f(x))
	}
}

func Not[T any](p Predicate[T]) Predicate[T] {
	return func(x T) bool {
		return !p(x)
	}
}

func And[T any](p1 Predicate[T], p2 Predicate[T]) Predicate[T] {
	return func(x T) bool {
		return p1(x) && p2(x)
	}
}

func Or[T any](p1 Predicate[T], p2 Predicate[T]) Predicate[T] {
	return func(x T) bool {
		return p1(x) || p2(x)
	}
}
