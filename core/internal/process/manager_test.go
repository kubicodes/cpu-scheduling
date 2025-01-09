package process

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	t.Run("should create empty manager", func(t *testing.T) {
		manager := NewManager()
		if len(manager.List()) != 0 {
			t.Error("new manager should have no processes")
		}
	})
}

func TestManager_Add(t *testing.T) {
	t.Run("should add process successfully", func(t *testing.T) {
		manager := NewManager()
		process := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))

		err := manager.Add(process)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(manager.List()) != 1 {
			t.Error("process was not added")
		}
	})

	t.Run("should return error when adding duplicate PID", func(t *testing.T) {
		manager := NewManager()
		process1 := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))
		process2 := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))

		manager.Add(process1)
		err := manager.Add(process2)

		if err == nil {
			t.Error("expected error when adding duplicate PID")
		}
	})
}

func TestManager_Remove(t *testing.T) {
	t.Run("should remove process successfully", func(t *testing.T) {
		manager := NewManager()
		process := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))

		manager.Add(process)
		err := manager.Remove(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(manager.List()) != 0 {
			t.Error("process was not removed")
		}
	})

	t.Run("should return error when removing non-existent process", func(t *testing.T) {
		manager := NewManager()
		err := manager.Remove(1)

		if err == nil {
			t.Error("expected error when removing non-existent process")
		}
	})
}

func TestManager_Get(t *testing.T) {
	t.Run("should get process successfully", func(t *testing.T) {
		manager := NewManager()
		process := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))

		manager.Add(process)
		result, err := manager.Get(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result.GetPID() != 1 {
			t.Errorf("expected PID 1, got %d", result.GetPID())
		}
	})

	t.Run("should return error when getting non-existent process", func(t *testing.T) {
		manager := NewManager()
		_, err := manager.Get(1)

		if err == nil {
			t.Error("expected error when getting non-existent process")
		}
	})
}

func TestManager_List(t *testing.T) {
	t.Run("should list all processes", func(t *testing.T) {
		manager := NewManager()
		process1 := NewPCB(1, NewTask(func() (any, error) { return nil, nil }))
		process2 := NewPCB(2, NewTask(func() (any, error) { return nil, nil }))

		manager.Add(process1)
		manager.Add(process2)

		processes := manager.List()
		if len(processes) != 2 {
			t.Errorf("expected 2 processes, got %d", len(processes))
		}
	})

	t.Run("should return empty list when no processes", func(t *testing.T) {
		manager := NewManager()
		processes := manager.List()

		if len(processes) != 0 {
			t.Errorf("expected empty list, got %d processes", len(processes))
		}
	})
}
