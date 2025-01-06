package types

// RegisterName represents valid CPU register names
type RegisterName string

const (
	// General Purpose Registers
	RAX RegisterName = "rax" // Accumulator
	RBX RegisterName = "rbx" // Base
	RCX RegisterName = "rcx" // Counter
	RDX RegisterName = "rdx" // Data
)

// ValidRegisters is a map of all valid register names
// This needs to be public as it's used by other packages for validation
var ValidRegisters = map[RegisterName]bool{
	RAX: true,
	RBX: true,
	RCX: true,
	RDX: true,
}
