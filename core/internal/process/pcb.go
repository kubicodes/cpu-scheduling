package process

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
	"time"
)

type PCB struct {
	pid       int
	state     types.ProcessState
	createdAt time.Time
	context   types.ProcessContext
}

func NewPCB(pid int) *PCB {
	return &PCB{
		pid:       pid,
		state:     types.NEW,
		createdAt: time.Now(),
		context:   NewProcessContext(),
	}
}

func (p *PCB) GetPID() int {
	return p.pid
}

func (p *PCB) GetState() types.ProcessState {
	return p.state
}

func (p *PCB) SetState(state types.ProcessState) error {
	if p.state == state {
		return fmt.Errorf("process is already in state %d", state)
	}

	if state == types.NEW {
		return fmt.Errorf("cannot set existing process to state NEW")
	}

	if p.state == types.TERMINATED {
		return fmt.Errorf("cannot change state of terminated process")
	}

	if p.state == types.NEW && state == types.RUNNING {
		return fmt.Errorf("process cannot be set to RUNNING state before being in READY state, current state is %d", p.state)
	}

	if p.state == types.READY && state == types.WAITING {
		return fmt.Errorf("process cannot be set from READY state to WAITING state, must go through RUNNING state first, current state is %d", p.state)
	}

	if p.state == types.WAITING && state == types.RUNNING {
		return fmt.Errorf("process cannot be set from WAITING state to RUNNING state, must go through READY state first, current state is %d", p.state)
	}

	p.state = state

	return nil
}

func (p *PCB) GetCreationTime() time.Time {
	return p.createdAt
}

func (p *PCB) GetContext() types.ProcessContext {
	return p.context
}
