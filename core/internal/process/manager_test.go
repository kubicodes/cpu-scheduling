package process

import (
	"cpu-scheduling/core/internal/types"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Error("expected non-nil manager")
	}
	if manager.processes == nil {
		t.Error("expected non-nil processes map")
	}
	if manager.nextPID != 1 {
		t.Errorf("expected nextPID to be 1, got %d", manager.nextPID)
	}
}

func TestManager_CreateProcess(t *testing.T) {
	t.Run("should create process with valid task", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}

		process, err := manager.CreateProcess(task)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if process == nil {
			t.Error("expected non-nil process")
		}
		if process.GetPID() != 1 {
			t.Errorf("expected PID 1, got %d", process.GetPID())
		}
		if process.GetState() != types.NEW {
			t.Errorf("expected state NEW, got %v", process.GetState())
		}
	})

	t.Run("should return error for nil task", func(t *testing.T) {
		manager := NewManager()
		process, err := manager.CreateProcess(nil)
		if err == nil {
			t.Error("expected error for nil task")
		}
		if process != nil {
			t.Error("expected nil process")
		}
	})

	t.Run("should assign unique PIDs", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}

		p1, _ := manager.CreateProcess(task)
		p2, _ := manager.CreateProcess(task)

		if p1.GetPID() == p2.GetPID() {
			t.Error("expected unique PIDs")
		}
	})
}

func TestManager_TerminateProcess(t *testing.T) {
	t.Run("should terminate existing process", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}
		process, _ := manager.CreateProcess(task)
		pid := process.GetPID()

		// Follow valid state transitions to RUNNING before terminating
		manager.SetProcessState(pid, types.READY)
		manager.SetProcessState(pid, types.RUNNING)

		err := manager.TerminateProcess(pid)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Process should no longer exist
		_, err = manager.GetProcess(pid)
		if err == nil {
			t.Error("expected error getting terminated process")
		}
	})

	t.Run("should return error for non-existent process", func(t *testing.T) {
		manager := NewManager()
		err := manager.TerminateProcess(999)
		if err == nil {
			t.Error("expected error terminating non-existent process")
		}
	})
}

func TestManager_SetProcessState(t *testing.T) {
	t.Run("should update process state", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}
		process, _ := manager.CreateProcess(task)
		pid := process.GetPID()

		err := manager.SetProcessState(pid, types.READY)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		updatedProcess, _ := manager.GetProcess(pid)
		if updatedProcess.GetState() != types.READY {
			t.Errorf("expected state READY, got %v", updatedProcess.GetState())
		}
	})

	t.Run("should return error for non-existent process", func(t *testing.T) {
		manager := NewManager()
		err := manager.SetProcessState(999, types.READY)
		if err == nil {
			t.Error("expected error setting state of non-existent process")
		}
	})

	t.Run("should validate state transitions", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}
		process, _ := manager.CreateProcess(task)
		pid := process.GetPID()

		// Try invalid transition from NEW to RUNNING
		err := manager.SetProcessState(pid, types.RUNNING)
		if err == nil {
			t.Error("expected error for invalid state transition NEW -> RUNNING")
		}

		// Try valid transitions
		if err := manager.SetProcessState(pid, types.READY); err != nil {
			t.Errorf("unexpected error for NEW -> READY: %v", err)
		}
		if err := manager.SetProcessState(pid, types.RUNNING); err != nil {
			t.Errorf("unexpected error for READY -> RUNNING: %v", err)
		}
	})
}

func TestManager_GetProcessesByState(t *testing.T) {
	t.Run("should return processes in specified state", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}

		// Create processes in different states
		p1, _ := manager.CreateProcess(task)
		p2, _ := manager.CreateProcess(task)

		// Follow valid state transitions
		manager.SetProcessState(p1.GetPID(), types.READY)
		manager.SetProcessState(p2.GetPID(), types.READY)
		manager.SetProcessState(p2.GetPID(), types.RUNNING)

		readyProcesses := manager.GetProcessesByState(types.READY)
		if len(readyProcesses) != 1 {
			t.Errorf("expected 1 ready process, got %d", len(readyProcesses))
		}

		runningProcesses := manager.GetProcessesByState(types.RUNNING)
		if len(runningProcesses) != 1 {
			t.Errorf("expected 1 running process, got %d", len(runningProcesses))
		}
	})

	t.Run("should return empty slice for state with no processes", func(t *testing.T) {
		manager := NewManager()
		processes := manager.GetProcessesByState(types.WAITING)
		if len(processes) != 0 {
			t.Errorf("expected empty slice, got %d processes", len(processes))
		}
	})
}

func TestManager_GetProcess(t *testing.T) {
	t.Run("should return existing process", func(t *testing.T) {
		manager := NewManager()
		task := &types.SimpleTask{
			ExecuteFn: func() (any, error) { return nil, nil },
		}
		created, _ := manager.CreateProcess(task)
		pid := created.GetPID()

		process, err := manager.GetProcess(pid)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if process == nil {
			t.Error("expected non-nil process")
		}
		if process.GetPID() != pid {
			t.Errorf("expected PID %d, got %d", pid, process.GetPID())
		}
	})

	t.Run("should return error for non-existent process", func(t *testing.T) {
		manager := NewManager()
		process, err := manager.GetProcess(999)
		if err == nil {
			t.Error("expected error getting non-existent process")
		}
		if process != nil {
			t.Error("expected nil process")
		}
	})
}

func TestManager_Concurrency(t *testing.T) {
	t.Run("should handle concurrent process creation and termination", func(t *testing.T) {
		manager := NewManager()
		numProcesses := 100
		done := make(chan bool)
		terminated := make(chan int, numProcesses)

		// Concurrent process creation
		go func() {
			for i := 0; i < numProcesses; i++ {
				task := &types.SimpleTask{
					ExecuteFn: func() (any, error) { return nil, nil },
				}
				if p, err := manager.CreateProcess(task); err == nil {
					pid := p.GetPID()
					// Follow valid state transitions
					if err := manager.SetProcessState(pid, types.READY); err == nil {
						if err := manager.SetProcessState(pid, types.RUNNING); err == nil {
							terminated <- pid
						}
					}
				}
			}
			close(terminated)
			done <- true
		}()

		// Concurrent process termination
		go func() {
			for pid := range terminated {
				if err := manager.TerminateProcess(pid); err != nil {
					t.Errorf("failed to terminate process %d: %v", pid, err)
				}
			}
			done <- true
		}()

		// Wait for completion
		<-done
		<-done

		// Verify final state
		remainingProcesses := 0
		for pid := 1; pid <= numProcesses; pid++ {
			if _, err := manager.GetProcess(pid); err == nil {
				remainingProcesses++
			}
		}

		if remainingProcesses > numProcesses/2 {
			t.Errorf("expected fewer than %d processes to remain, got %d", numProcesses/2, remainingProcesses)
		}
	})
}
