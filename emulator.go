package main

type CpuFlags struct {
	Parity   bool
	Zero     bool
	Sign     bool
	Carry    bool
	AuxCarry bool
}

type CpuState struct {
	RegA      uint8
	RegB      uint8
	RegC      uint8
	RegD      uint8
	RegE      uint8
	RegH      uint8
	RegL      uint8
	SP        uint16
	PC        uint16
	Memory    []byte
	IntEnable bool
	Condition CpuFlags
}

type Operator func(a uint16, b uint16) uint16

func bitParity(byte uint16) bool {
	var acc uint16 = 0
	for i := 0; i < 8; i++ {
		acc += (byte >> i) & 0x01
	}

	return acc%2 == 0
}

func AritmethicOperation(state *CpuState, value uint8, useCarry bool, op Operator) {
	var answer uint16 = op(uint16(state.RegA), uint16(value))

	if state.Condition.Carry && useCarry {
		answer = op(answer, 1)
	}

	// Cpu Flags
	state.Condition.Zero = ((answer & 0xff) == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = (answer > 0xff)
	state.Condition.Parity = bitParity(answer & 0xff)
	state.Condition.AuxCarry = (((state.RegA & value) & 0x04) == 0x04)

	state.RegA = uint8(answer)
}

func LogicalOperation(state *CpuState, value uint8, op Operator) {
	var answer uint16 = op(uint16(state.RegA), uint16(value))

	// Cpu Flags
	state.Condition.Zero = (answer == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = false
	state.Condition.Parity = bitParity(answer & 0xff)

	state.RegA = uint8(answer)
}

func CompareOperation(state *CpuState, value uint8) {
	var answer uint16 = uint16(state.RegA) - uint16(value)

	// Cpu Flags
	state.Condition.Zero = ((answer & 0xff) == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = (answer > 0xff)
	state.Condition.Parity = bitParity(answer & 0xff)
	state.Condition.AuxCarry = (((state.RegA & value) & 0x04) == 0x04)
}
