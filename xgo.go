package main

import (
	"fmt"
	"io/ioutil"
);

func readFile (filename string) string {
	bytes, ok := ioutil.ReadFile (filename)
	if ok != nil {
		panic (ok)
	}
	return string (bytes);
}

func main () {
	str := readFile ("/dev/stdin");
	fmt.Println("\t.global _mymain")
	fmt.Println("_mymain:");
	fmt.Printf("\tmovl $%s, %%eax\n", string(str));
	fmt.Println("\tret");
}
