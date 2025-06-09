package digital

type emitter[T any] struct {
	listeners []func(T)
}

func newEmitter[T any]() *emitter[T] {
	listeners := make([]func(T), 0)
	return &emitter[T]{listeners}
}

func (e *emitter[T]) on(listener func(T)) {
	e.listeners = append(e.listeners, listener)
}

func (e *emitter[T]) emit(value T) {
	for _, listener := range e.listeners {
		listener(value)
	}
}
