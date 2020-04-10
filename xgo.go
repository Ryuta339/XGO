package main

import (
	"fmt"
	"io/ioutil"
);

func main() {
	contents, ok := ioutil.ReadFile ("/dev/stdin")
	if ok != nil {
		panic (ok)
	}
	fmt.Println("\t.global _mymain")
	fmt.Println("_mymain:");
	fmt.Printf("\tmovl $%s, %%eax\n", string(contents));
	fmt.Println("\tret");
}
