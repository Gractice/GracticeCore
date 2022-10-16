package util

type CastResult[T any] struct {
	ok  bool
	val T
}

func (r *CastResult[T]) Is(f func(T)) *CastResult[T] {
	if r.ok {
		f(r.val)
	}
	return r
}

func (r *CastResult[T]) Else(f func()) *CastResult[T] {
	if !r.ok {
		f()
	}
	return r
}

func Cast[T any](v interface{}) *CastResult[T] {
	if vv, ok := v.(T); ok {
		return &CastResult[T]{
			ok:  true,
			val: vv,
		}
	}
	return &CastResult[T]{
		ok: false,
	}
}
