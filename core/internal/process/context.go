package process

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
)

type ProcessContext struct {
	programCounter uint64
	registers      map[types.RegisterName]uint64
	state          *contextState
}

type contextState struct {
	programCounter uint64
	registers      map[types.RegisterName]uint64
}

func NewProcessContext() *ProcessContext {
	// Create registers map
	registers := make(map[types.RegisterName]uint64)

	// Initialize all valid registers to 0
	for reg := range types.ValidRegisters {
		registers[reg] = 0
	}

	return &ProcessContext{
		programCounter: 0,
		registers:      registers,
		state:          nil,
	}
}

func (p *ProcessContext) GetProgramCounter() uint64 {
	return p.programCounter
}

func (p *ProcessContext) SetProgramCounter(pc uint64) error {
	// Check if program counter is aligned to 4 bytes 32-bit (4 bytes) is a common instruction size
	if pc%4 != 0 {
		return fmt.Errorf("program counter must be 4-byte aligned, got %d", pc)
	}

	p.programCounter = pc
	return nil
}

func (p *ProcessContext) GetRegisterValue(register types.RegisterName) (uint64, error) {
	value, exists := p.registers[register]

	if !exists {
		return 0, fmt.Errorf("register %s not found", register)
	}

	return value, nil
}

func (p *ProcessContext) SetRegisterValue(register types.RegisterName, value uint64) error {
	if !types.ValidRegisters[register] {
		return fmt.Errorf("invalid register name: %s", register)
	}

	p.registers[register] = value
	return nil
}

func (p *ProcessContext) SaveState() {
	registersCopy := make(map[types.RegisterName]uint64)

	for reg, value := range p.registers {
		registersCopy[reg] = value
	}

	state := &contextState{
		programCounter: p.programCounter,
		registers:      registersCopy,
	}

	p.state = state
}

func (p *ProcessContext) LoadState() {
	if p.state == nil {
		return
	}

	p.programCounter = p.state.programCounter

	for reg, value := range p.state.registers {
		p.registers[reg] = value
	}

	p.state = nil
}
