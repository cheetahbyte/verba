package registries

import (
	"fmt"
	"reflect"
)

type Registry[T any] struct {
	contents map[string]T
}

func (r *Registry[T]) Register(name string, item T) {
	fmt.Println("Registering", name, item)
	r.contents[name] = item
}

func (r *Registry[T]) Get(name string) (T, bool) {
	var zero T
	fmt.Println(len(r.contents), reflect.TypeOf(zero))
	if val, ok := r.contents[name]; ok {
		return val, ok
	}
	return zero, false
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		contents: make(map[string]T),
	}
}
