package process

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
	"sync"
)

type Manager struct {
	processes map[int]types.Process
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		processes: make(map[int]types.Process),
	}
}

func (m *Manager) Add(process types.Process) error {
	m.mu.Lock()
	// avoid race conditions
	defer m.mu.Unlock()

	pid := process.GetPID()
	if _, exists := m.processes[pid]; exists {
		return fmt.Errorf("process with PID %d already exists", pid)
	}

	m.processes[pid] = process
	return nil
}

func (m *Manager) Remove(pid int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.processes[pid]; !exists {
		return fmt.Errorf("process with PID %d not found", pid)
	}

	delete(m.processes, pid)
	return nil
}

func (m *Manager) Get(pid int) (types.Process, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	process, exists := m.processes[pid]
	if !exists {
		return nil, fmt.Errorf("process with PID %d not found", pid)
	}

	return process, nil
}

func (m *Manager) List() []types.Process {
	m.mu.RLock()
	defer m.mu.RUnlock()

	processes := make([]types.Process, 0, len(m.processes))
	for _, process := range m.processes {
		processes = append(processes, process)
	}

	return processes
}
