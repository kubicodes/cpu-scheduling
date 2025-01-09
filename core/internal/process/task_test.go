package process

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	t.Run("should execute the task", func(t *testing.T) {
		task := NewTask(func() (any, error) {
			return 1, nil
		})

		result, err := task.Execute()

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if result != 1 {
			t.Errorf("expected result to be 1, got %v", result)
		}
	})

	t.Run("should return error if the task returns an error", func(t *testing.T) {
		task := NewTask(func() (any, error) {
			return nil, fmt.Errorf("mock error")
		})

		_, err := task.Execute()

		if err == nil {
			t.Error("expected error, got nil")
		}

		expectedErrorMsg := "mock error"
		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message to be %s, got %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("should work with void functions", func(t *testing.T) {
		task := NewTask(func() (any, error) {
			return nil, nil
		})

		result, err := task.Execute()

		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}

		if result != nil {
			t.Errorf("expected result to be nil, got %v", result)
		}

	})
}
