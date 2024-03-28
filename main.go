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

	var CpuState = CpuState{}

	for pc := 0; pc < len(buff); pc++ {
		UpdateState(&CpuState, buff[pc])
	}
}
