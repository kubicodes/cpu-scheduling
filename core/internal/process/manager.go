package process

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
	"sync"
)

// Manager handles process lifecycle and state management
type Manager struct {
	processes map[int]types.Process
	nextPID   int
	mu        sync.RWMutex
}

// NewManager creates a new process manager
func NewManager() *Manager {
	return &Manager{
		processes: make(map[int]types.Process),
		nextPID:   1,
	}
}

// CreateProcess creates a new process with the given task
func (m *Manager) CreateProcess(task types.Task) (types.Process, error) {
	if task == nil {
		return nil, fmt.Errorf("cannot create process with nil task")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	pid := m.nextPID
	m.nextPID++

	pcb := NewPCB(pid, task)
	m.processes[pid] = pcb

	return pcb, nil
}

// TerminateProcess terminates the process with the given PID
func (m *Manager) TerminateProcess(pid int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	process, exists := m.processes[pid]
	if !exists {
		return fmt.Errorf("process with PID %d not found", pid)
	}

	// Set state to terminated
	if err := process.SetState(types.TERMINATED); err != nil {
		return fmt.Errorf("failed to set process state to terminated: %v", err)
	}

	// Remove from processes map
	delete(m.processes, pid)
	return nil
}

// SetProcessState updates the state of the process with the given PID
func (m *Manager) SetProcessState(pid int, newState types.ProcessState) error {
	m.mu.Lock() // Use Lock instead of RLock since we're modifying process state
	defer m.mu.Unlock()

	process, exists := m.processes[pid]
	if !exists {
		return fmt.Errorf("process with PID %d not found", pid)
	}

	// Validate state transition
	currentState := process.GetState()
	if err := validateStateTransition(currentState, newState); err != nil {
		return err
	}

	return process.SetState(newState)
}

// validateStateTransition checks if the state transition is valid
func validateStateTransition(from, to types.ProcessState) error {
	// Cannot transition to NEW state
	if to == types.NEW {
		return fmt.Errorf("cannot transition to NEW state")
	}

	// Valid transitions
	switch from {
	case types.NEW:
		if to != types.READY {
			return fmt.Errorf("process in NEW state can only transition to READY state")
		}
	case types.READY:
		if to != types.RUNNING {
			return fmt.Errorf("process in READY state can only transition to RUNNING state")
		}
	case types.RUNNING:
		if to != types.READY && to != types.WAITING && to != types.TERMINATED {
			return fmt.Errorf("process in RUNNING state can only transition to READY, WAITING, or TERMINATED state")
		}
	case types.WAITING:
		if to != types.READY {
			return fmt.Errorf("process in WAITING state can only transition to READY state")
		}
	case types.TERMINATED:
		return fmt.Errorf("cannot transition from TERMINATED state")
	}

	return nil
}

// GetProcessesByState returns all processes in the given state
func (m *Manager) GetProcessesByState(state types.ProcessState) []types.Process {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []types.Process
	for _, p := range m.processes {
		if p.GetState() == state {
			result = append(result, p)
		}
	}
	return result
}

// GetProcess returns the process with the given PID
func (m *Manager) GetProcess(pid int) (types.Process, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	process, exists := m.processes[pid]
	if !exists {
		return nil, fmt.Errorf("process with PID %d not found", pid)
	}
	return process, nil
}
