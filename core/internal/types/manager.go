package types

type ProcessManager interface {
	Add(process Process) error
	Remove(pid int) error
	Get(pid int) (Process, error)
	List() []Process
}
