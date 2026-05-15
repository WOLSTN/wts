package core

import "slices"

type Arena[T any] struct {
	data []T
}

func (a *Arena[T]) New() *T {
	if len(a.data) == cap(a.data) {
		nextSize := nextArenaSize(len(a.data))
		a.data = slices.Grow[[]T](nil, nextSize)
	}
	index := len(a.data)
	a.data = a.data[:index+1]
	return &a.data[index]
}

func (a *Arena[T]) NewSlice(size int) []T {
	if size == 0 {
		return nil
	}
	if len(a.data)+size > cap(a.data) {
		nextSize := nextArenaSize(len(a.data))
		if size > nextSize {
			return make([]T, size)
		}
		a.data = slices.Grow[[]T](nil, nextSize)
	}
	newLen := len(a.data) + size
	slice := a.data[len(a.data):newLen:newLen]
	a.data = a.data[:newLen]
	return slice
}

func (a *Arena[T]) NewSlice1(t T) []T {
	slice := a.NewSlice(1)
	slice[0] = t
	return slice
}

func (a *Arena[T]) Clone(t []T) []T {
	if len(t) == 0 {
		return nil
	}
	slice := a.NewSlice(len(t))
	copy(slice, t)
	return slice
}

func (a *Arena[T]) Data() []T {
	return a.data
}

func nextArenaSize(size int) int {
	size = max(size, 1)
	size = min(size*2, 256)
	return size
}
