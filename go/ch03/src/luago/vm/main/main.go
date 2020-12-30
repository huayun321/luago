package main

import (
	"fmt"
	"luago/go/ch03/src/luago/vm"
)

func main() {
	instruct := vm.Instruction(0x00400006)
	instruct.OpName()

	fmt.Println(instruct.OpName())
	fmt.Println(instruct.Opcode())
}
