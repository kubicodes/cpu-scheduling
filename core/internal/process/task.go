package process

type Task[T any] struct {
	executeFn func() (T, error)
}

func NewTask[T any](fn func() (T, error)) *Task[T] {
	return &Task[T]{executeFn: fn}
}

func (t *Task[T]) Execute() (T, error) {
	return t.executeFn()
}
