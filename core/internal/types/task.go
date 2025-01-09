package types

type Task interface {
	Execute() (any, error)
}
