package types

type Task interface {
	Execute() (any, error)
}

// SimpleTask is a basic implementation of Task interface
type SimpleTask struct {
	ExecuteFn func() (any, error)
}

func (t *SimpleTask) Execute() (any, error) {
	return t.ExecuteFn()
}
