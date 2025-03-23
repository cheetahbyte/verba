package registries

type Registry[T any] struct {
	contents map[string]T
}

func (r *Registry[T]) Register(name string, item T) {
	r.contents[name] = item
}

func (r *Registry[T]) Get(name string) (T, bool) {
	var zero T
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
