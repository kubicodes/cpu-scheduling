package types

type ProcessContext interface {
	GetProgramCounter() uint64
	SetProgramCounter(pc uint64) error
	GetRegisterValue(register RegisterName) (uint64, error)
	SetRegisterValue(register RegisterName, value uint64) error
	SaveState()
	LoadState()
}
