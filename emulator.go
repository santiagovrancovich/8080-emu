package main

type CpuFlags struct {
	Parity   bool
	Zero     bool
	Sign     bool
	Carry    bool
	AuxCarry uint8
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
	Memory    *uint8
	IntEnable bool
	Condition CpuFlags
}

func bitParity(byte uint16) bool {
	var acc uint16 = 0
	for i := range 8 {
		acc += (byte >> i) & 0x01
	}

	return acc%2 == 1
}

func AddInstruction(state *CpuState, register *uint8) {
	var answer uint16 = uint16(state.RegA) + uint16(*register)

	// Cpu Flags
	state.Condition.Zero = ((answer & 0xff) == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = (answer > 0xff)
	state.Condition.Parity = bitParity(answer & 0xff)

	state.RegA = uint8(answer & 0xff)
}
