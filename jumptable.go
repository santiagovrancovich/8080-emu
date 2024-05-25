package main

import "fmt"

type Operation struct {
	Name     string
	ByteSize int
}

var Registers = []string{"B", "C", "D", "E", "H", "L", "M", "A"}
var RegisterPair = []string{"BC", "DE", "HL", "SP"}
var Conditions = []string{"NZ", "Z", "NC", "C", "PO", "PE", "P", "M"}

func getRegisterVal(state *CpuState, reg byte) *uint8 {
	switch reg {
	case 0:
		return &state.RegB
	case 1:
		return &state.RegC
	case 2:
		return &state.RegD
	case 3:
		return &state.RegE
	case 4:
		return &state.RegH
	case 5:
		return &state.RegL
	case 6:
		addr := (uint16(state.RegH) << 8) | uint16(state.RegL)
		return &state.Memory[addr]
	default:
		return &state.RegA
	}
}

// implementation pending
func UpdateState(state *CpuState, op byte) {
	switch op {
	case 0x76:
		state.PC = uint16(len(state.Memory)) + 1
	case 0xc6:
		AritmethicOperation(state, state.Memory[state.PC+1], false,
			func(a, b uint16) uint16 { return a + b })
		state.PC += 1
	case 0xce:
		AritmethicOperation(state, state.Memory[state.PC+1], true,
			func(a, b uint16) uint16 { return a + b })
		state.PC += 1
	case 0xd6:
		AritmethicOperation(state, state.Memory[state.PC+1], false,
			func(a, b uint16) uint16 { return a - b })
		state.PC += 1
	case 0xd9:
		AritmethicOperation(state, state.Memory[state.PC+1], true,
			func(a, b uint16) uint16 { return a - b })
		state.PC += 1
	case 0xe6:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a & b })
		state.PC += 1
	case 0xee:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a ^ b })
		state.PC += 1
	case 0xf6:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a | b })
		state.PC += 1
	case 0xf9:
		state.SP = uint16(state.RegL)<<8 | uint16(state.RegH)
	case 0xfe:
		CompareOperation(state, state.Memory[state.PC+1])
		state.PC += 1
	}

	switch {
	case 0x80 <= op && op <= 0x87:
		AritmethicOperation(state, *getRegisterVal(state, op&0x7), false,
			func(a, b uint16) uint16 { return a + b })
	case 0x88 <= op && op <= 0x8f:
		AritmethicOperation(state, *getRegisterVal(state, op&0x7), true,
			func(a, b uint16) uint16 { return a + b })
	case 0x90 <= op && op <= 0x97:
		AritmethicOperation(state, *getRegisterVal(state, op&0x7), false,
			func(a, b uint16) uint16 { return a - b })
	case 0x98 <= op && op <= 0x9f:
		AritmethicOperation(state, *getRegisterVal(state, op&0x7), true,
			func(a, b uint16) uint16 { return a - b })
	case 0xa0 <= op && op <= 0xa7:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a & b })
	case 0xa8 <= op && op <= 0xaf:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a ^ b })
	case 0xb0 <= op && op <= 0xb7:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a | b })
	case op&0xb8 == 0xb8:
		CompareOperation(state, *getRegisterVal(state, op&0x07))
	case (op^0xce)|0x30 == 0xff:
		if op>>4 == 3 {
			state.SP = uint16(state.Memory[state.PC+1])<<8 | uint16(state.Memory[state.PC+2])
		} else {
			lb := getRegisterVal(state, (op>>4)*2)
			hb := getRegisterVal(state, ((op>>4)*2)+1)
			*lb = state.Memory[state.PC+1]
			*hb = state.Memory[state.PC+2]
		}
		state.PC += 2
	case (op^0xc1)|0x38 == 0xff:
		reg := getRegisterVal(state, op>>3)
		*reg = state.Memory[state.PC+1]
		state.PC += 1
	case op>>6 == 0x1:
		dst := getRegisterVal(state, (op&0x3f)>>3)
		src := getRegisterVal(state, op&0x07)
		*dst = *src
	case op&0xc2 == 0xc2:
		if op&0x01 == 0x01 {
			state.PC = uint16(state.Memory[state.PC+2])<<8 | uint16(state.Memory[state.PC+1]) - 1
		} else {
			JumpOperation(state, op)
		}
	}
}

func decodeMOV(op byte) Operation {
	val := 0x40 ^ op
	dest := (val >> 3) & 0x7
	source := val & 0x7

	return Operation{fmt.Sprintf("MOV %s, %s", Registers[dest], Registers[source]), 1}
}

func DecodeOpCode(op byte) Operation {
	//Hardcoded OPs
	switch op {
	case 0x0, 0x10, 0x20, 0x30, 0x08, 0x18, 0x28, 0x38:
		return Operation{"NOP", 1}
	case 0x07:
		return Operation{"RLC", 1}
	case 0x0f:
		return Operation{"RRC", 1}
	case 0x17:
		return Operation{"RAL", 1}
	case 0x1f:
		return Operation{"RAR", 1}
	case 0x22:
		return Operation{"SHLD", 3}
	case 0x27:
		return Operation{"DAA", 1}
	case 0x2a:
		return Operation{"LHLD", 3}
	case 0x2f:
		return Operation{"CMA", 1}
	case 0x32:
		return Operation{"STA", 3}
	case 0x37:
		return Operation{"STC", 1}
	case 0x3a:
		return Operation{"LDA", 3}
	case 0x3f:
		return Operation{"CMC", 1}
	case 0x76:
		return Operation{"HLT", 1}
	case 0xc3:
		return Operation{"JMP", 3}
	case 0xc6:
		return Operation{"ADI", 2}
	case 0xc9:
		return Operation{"RET", 1}
	case 0xce:
		return Operation{"ACI", 2}
	case 0xcd, 0xdd, 0xed, 0xfd:
		return Operation{"CALL", 3}
	case 0xd3:
		return Operation{"OUT", 2}
	case 0xd6:
		return Operation{"SUI", 2}
	case 0xdb:
		return Operation{"IN", 2}
	case 0xde:
		return Operation{"SBI", 2}
	case 0xe3:
		return Operation{"XTHL", 1}
	case 0xe6:
		return Operation{"ANI", 2}
	case 0xe9:
		return Operation{"PCHL", 1}
	case 0xeb:
		return Operation{"XCHG", 1}
	case 0xee:
		return Operation{"XRI", 2}
	case 0xf3:
		return Operation{"DI", 1}
	case 0xf6:
		return Operation{"ORI", 2}
	case 0xf9:
		return Operation{"SPHL", 1}
	case 0xfb:
		return Operation{"EI", 1}
	case 0xfe:
		return Operation{"CPI", 2}
	}

	//Variable OPs
	switch {
	case (op^0xce)|0x30 == 0xff:
		return Operation{fmt.Sprintf("LXI %s,", RegisterPair[op>>4]), 3}
	case (op^0xc1)|0x38 == 0xff:
		return Operation{fmt.Sprintf("MVI %s", Registers[op>>3]), 2}
	case (op^0xc3)|0x38 == 0xff:
		return Operation{fmt.Sprintf("INR %s", Registers[op>>3]), 1}
	case (op^0xc2)|0x38 == 0xff:
		return Operation{fmt.Sprintf("DCR %s", Registers[op>>3]), 1}
	case (op^0xc5)|0x30 == 0xff:
		return Operation{fmt.Sprintf("LDAX %s", RegisterPair[op>>4]), 1}
	case (op^0xcd)|0x30 == 0xff:
		return Operation{fmt.Sprintf("STAX %s", RegisterPair[op>>4]), 1}
	case (op^0xcc)|0x30 == 0xff:
		return Operation{fmt.Sprintf("INX %s", RegisterPair[op>>4]), 1}
	case (op^0xc4)|0x30 == 0xff:
		return Operation{fmt.Sprintf("DCX %s", RegisterPair[op>>4]), 1}
	case (op^0xc6)|0x30 == 0xff:
		return Operation{fmt.Sprintf("DAD %s", RegisterPair[op>>4]), 1}
	case 0x40 <= op && op <= 0x7f:
		return decodeMOV(op)
	case 0x80 <= op && op <= 0x87:
		return Operation{fmt.Sprintf("ADD %s", Registers[op&0x7]), 1}
	case 0x88 <= op && op <= 0x8f:
		return Operation{fmt.Sprintf("ADC %s", Registers[op&0x7]), 1}
	case 0x90 <= op && op <= 0x97:
		return Operation{fmt.Sprintf("SUB %s", Registers[op&0x7]), 1}
	case 0x98 <= op && op <= 0x9f:
		return Operation{fmt.Sprintf("SBB %s", Registers[op&0x7]), 1}
	case 0xa0 <= op && op <= 0xa7:
		return Operation{fmt.Sprintf("ANA %s", Registers[op&0x7]), 1}
	case 0xa8 <= op && op <= 0xaf:
		return Operation{fmt.Sprintf("XRA %s", Registers[op&0x7]), 1}
	case 0xb0 <= op && op <= 0xb7:
		return Operation{fmt.Sprintf("ORA %s", Registers[op&0x7]), 1}
	case 0xb8 <= op && op <= 0xbf:
		return Operation{fmt.Sprintf("CMP %s", Registers[op&0x7]), 1}
	case (op^0x0a)|0x30 == 0xff:
		return Operation{fmt.Sprintf("PUSH %s", RegisterPair[(op&0x30)>>4]), 1}
	case (op^0x0e)|0x30 == 0xff:
		return Operation{fmt.Sprintf("POP %s", RegisterPair[(op&0x30)>>4]), 1}
	case (op^0x05)|0x38 == 0xff:
		return Operation{fmt.Sprintf("J%s", Conditions[(op&0x38)>>3]), 3}
	case (op^0x03)|0x38 == 0xff:
		return Operation{fmt.Sprintf("C%s", Conditions[(op&0x38)>>3]), 3}
	case (op^0x07)|0x38 == 0xff:
		return Operation{fmt.Sprintf("R%s", Conditions[(op&0x38)>>3]), 1}
	case op|0x38 == 0xff:
		return Operation{fmt.Sprintf("RST %d", (op&0x38)>>3), 1}
	default:
		return Operation{"Unkown Operation", 1}
	}
}
