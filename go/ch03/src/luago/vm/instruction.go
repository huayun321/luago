package vm

type Instruction uint32

func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}

func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}
