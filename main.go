package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello there")

	if len(os.Args) < 2 {
		log.Fatal("Mate you fucked")
	}

	buff, err := os.ReadFile(os.Args[1])

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
