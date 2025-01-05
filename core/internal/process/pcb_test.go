package process

import (
	"cpu-scheduling/core/internal/types"
	"testing"
)

func TestPCB_GetPID(t *testing.T) {
	t.Run("should return the pid", func(t *testing.T) {
		pcb := &PCB{pid: 1}

		pid := pcb.GetPID()

		if pid != 1 {
			t.Errorf("expected pid to be 1, got %d", pid)
		}
	})
}

func TestPCB_GetState(t *testing.T) {
	t.Run("should return the state", func(t *testing.T) {
		pcb := &PCB{state: types.NEW}

		state := pcb.GetState()

		if state != types.NEW {
			t.Errorf("expected state to be NEW, got %d", state)
		}

		newPcb := &PCB{state: types.READY}

		state = newPcb.GetState()

		if state != types.READY {
			t.Errorf("expected state to be READY, got %d", state)
		}
	})
}

func TestPCB_SetState(t *testing.T) {
	t.Run("should return error if the state to set is the same as the current state", func(t *testing.T) {
		pcb := &PCB{state: types.READY}
		stateToSet := types.READY
		err := pcb.SetState(stateToSet)

		if err == nil {
			t.Errorf("expected error to be returned, got nil")
		}

		expectedErrorMsg := "process is already in state 1"

		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message to be %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should return error if the state to set is NEW", func(t *testing.T) {
		pcb := &PCB{state: types.RUNNING}
		stateToSet := types.NEW
		err := pcb.SetState(stateToSet)

		if err == nil {
			t.Errorf("expected error to be returned, got nil")
		}

		expectedErrorMsg := "cannot set existing process to state NEW"

		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message to be %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should return error if a terminated process is attempted to be set to a new state", func(t *testing.T) {
		pcb := &PCB{state: types.TERMINATED}
		stateToSet := types.READY
		err := pcb.SetState(stateToSet)

		if err == nil {
			t.Errorf("expected error to be returned, got nil")
		}

		expectedErrorMsg := "cannot change state of terminated process"

		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message to be %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should return error if the process wants to be set from new to running, before being in READY state", func(t *testing.T) {
		pcb := &PCB{state: types.NEW}
		stateToSet := types.RUNNING
		err := pcb.SetState(stateToSet)

		if err == nil {
			t.Errorf("expected error to be returned, got nil")
		}

		expectedErrorMsg := "process cannot be set to RUNNING state before being in READY state, current state is 0"

		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message to be %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should return error when transitioning from READY to WAITING", func(t *testing.T) {
		pcb := &PCB{state: types.READY}
		err := pcb.SetState(types.WAITING)

		if err == nil {
			t.Error("expected error when transitioning from READY to WAITING")
		}

		expectedMsg := "process cannot be set from READY state to WAITING state, must go through RUNNING state first, current state is 1"
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("should return error when transitioning from WAITING to RUNNING", func(t *testing.T) {
		pcb := &PCB{state: types.WAITING}
		err := pcb.SetState(types.RUNNING)

		if err == nil {
			t.Error("expected error when transitioning from WAITING to RUNNING")
		}

		expectedMsg := "process cannot be set from WAITING state to RUNNING state, must go through READY state first, current state is 3"
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("should set state from NEW to READY", func(t *testing.T) {
		pcb := &PCB{state: types.NEW}
		err := pcb.SetState(types.READY)

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if pcb.state != types.READY {
			t.Errorf("expected state to be READY, got %d", pcb.state)
		}
	})

	t.Run("should set state from READY to RUNNING", func(t *testing.T) {
		pcb := &PCB{state: types.READY}
		err := pcb.SetState(types.RUNNING)

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if pcb.state != types.RUNNING {
			t.Errorf("expected state to be RUNNING, got %d", pcb.state)
		}
	})

	t.Run("should set state from RUNNING to WAITING", func(t *testing.T) {
		pcb := &PCB{state: types.RUNNING}
		err := pcb.SetState(types.WAITING)

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if pcb.state != types.WAITING {
			t.Errorf("expected state to be WAITING, got %d", pcb.state)
		}
	})

	t.Run("should set state from RUNNING to TERMINATED", func(t *testing.T) {
		pcb := &PCB{state: types.RUNNING}
		err := pcb.SetState(types.TERMINATED)

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if pcb.state != types.TERMINATED {
			t.Errorf("expected state to be TERMINATED, got %d", pcb.state)
		}
	})
}
