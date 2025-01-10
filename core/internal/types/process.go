package types

import (
	"time"
)

type ProcessState int

const (
	NEW ProcessState = iota
	READY
	RUNNING
	WAITING
	TERMINATED
)

type Process interface {
	GetPID() int
	GetState() ProcessState
	SetState(state ProcessState) error
	GetCreationTime() time.Time
	GetContext() ProcessContext
	ExecuteTask() (any, error)
	// Time tracking
	GetTimeInState() time.Duration // How long the process has been in current state
	GetTotalTime() time.Duration   // Total time since process creation
}
