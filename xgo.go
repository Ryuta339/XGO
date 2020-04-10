package main

import "fmt"

func main() {
	fmt.Println("\t.global _main")
	fmt.Println("_main:");
	fmt.Println("\tmovl $0, %eax");
	fmt.Println("\tret");
}
