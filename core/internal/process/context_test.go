package process

import (
	"cpu-scheduling/core/internal/types"
	"testing"
)

func TestNewProcessContext(t *testing.T) {
	t.Run("should create context with default values", func(t *testing.T) {
		var programCounter uint64 = 0
		context := NewProcessContext()

		if context.GetProgramCounter() != programCounter {
			t.Errorf("expected program counter %d, got %d", programCounter, context.GetProgramCounter())
		}

		expectedLen := len(types.ValidRegisters)
		if len(context.registers) != expectedLen {
			t.Errorf("expected %d registers, got %d", expectedLen, len(context.registers))
		}
	})

	t.Run("should initialize all valid registers to zero", func(t *testing.T) {
		context := NewProcessContext()
		for reg := range types.ValidRegisters {
			value, err := context.GetRegisterValue(reg)
			if err != nil || value != 0 {
				t.Errorf("expected register %s to be 0, got %d", reg, value)
			}
		}
	})
}

func TestProcessContext_SetProgramCounter(t *testing.T) {
	t.Run("should set program counter to valid value", func(t *testing.T) {
		context := NewProcessContext()
		newPC := uint64(100)

		err := context.SetProgramCounter(newPC)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if context.GetProgramCounter() != newPC {
			t.Errorf("expected program counter to be %d, got %d",
				newPC, context.GetProgramCounter())
		}
	})

	t.Run("should return error when program counter is not aligned", func(t *testing.T) {
		context := NewProcessContext()
		unalignedPC := uint64(3) // Not aligned to 4 bytes

		err := context.SetProgramCounter(unalignedPC)
		if err == nil {
			t.Error("expected error for unaligned program counter")
		}
	})
}

func TestProcessContext_GetRegisterValue(t *testing.T) {
	t.Run("should return error for non-existent register", func(t *testing.T) {
		context := NewProcessContext()
		_, err := context.GetRegisterValue(types.RegisterName("invalid"))
		if err == nil {
			t.Error("expected error for non-existent register")
		}
	})

	t.Run("should get value of existing register", func(t *testing.T) {
		context := NewProcessContext()
		expectedValue := uint64(42)

		// First set the value
		err := context.SetRegisterValue(types.RAX, expectedValue)
		if err != nil {
			t.Fatalf("failed to set register value: %v", err)
		}

		// Then get and verify
		value, err := context.GetRegisterValue(types.RAX)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if value != expectedValue {
			t.Errorf("expected value %d, got %d", expectedValue, value)
		}
	})
}

func TestProcessContext_SetRegisterValue(t *testing.T) {
	t.Run("should set value for valid register", func(t *testing.T) {
		context := NewProcessContext()
		err := context.SetRegisterValue(types.RAX, 42)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify value was set
		value, err := context.GetRegisterValue(types.RAX)
		if err != nil || value != 42 {
			t.Errorf("expected value 42, got %d with error: %v", value, err)
		}
	})

	t.Run("should set different values for different registers", func(t *testing.T) {
		context := NewProcessContext()

		// Set values
		err1 := context.SetRegisterValue(types.RAX, 42)
		err2 := context.SetRegisterValue(types.RBX, 100)

		if err1 != nil || err2 != nil {
			t.Errorf("unexpected errors: %v, %v", err1, err2)
		}

		// Verify values
		val1, _ := context.GetRegisterValue(types.RAX)
		val2, _ := context.GetRegisterValue(types.RBX)

		if val1 != 42 || val2 != 100 {
			t.Errorf("expected values 42 and 100, got %d and %d", val1, val2)
		}
	})

	t.Run("should return error for invalid register", func(t *testing.T) {
		context := NewProcessContext()
		err := context.SetRegisterValue(types.RegisterName("invalid"), 42)
		if err == nil {
			t.Error("expected error for invalid register")
		}
	})
}
