package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello there")
	buff, err := os.ReadFile("invaders.h")

	if err != nil {
		log.Fatal("Mate you fucked")
	}

	var CpuState = CpuState{Memory: buff}

	for ; int(CpuState.PC) < len(buff); CpuState.PC++ {
		UpdateState(&CpuState, buff[CpuState.PC])
	}

	CpuState.Memory = make([]byte, 0)
	fmt.Printf("%+v\n", CpuState)
}
