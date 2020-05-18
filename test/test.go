package main

import "fmt"
import (
	"strings"
	"strconv"
)


func f1 () {
	// this is a comment
	puts ("")
	puts ("hello")
	puts ("world")
}

func f2 () {
	/* 
	 * this is a block comment
	 */
	var i int
	i = 3
	printf ("%d\n", i)
}

func f3 () {
	var j int = 5
	var k int = 6
	printf ("%d\n", j + k) // this is a comment
}


var gi int = 30200

func f4 () {
	gi = gi + 3
	printf ("%d\n", gi)
}

func f5 (a int, b int) {
	printf ("%d\n", a + b)
}

func f6 (a int, b int) {
	var c int = 3
	var d int = 4
	printf ("%d\n", a + b + c + d)
}

func f7 (a int) {
	var b int = 2
	{
		var c int = 3
		printf("%d\n", a + b * c)
	}
	printf("%d\n", a + b)
	{
		var d int
		printf("%d\n", a + b + d)
	}
}

func main () {
	printf ("%d\n", 2 + 5)
	printf ("%d\n", 10 - 4)
	printf ("%d\n", 4 * 3)
	printf ("%d\n", 1 * 2 + 3 * 4)
	printf ("%d\n", 1 + 2 * 3 + 4)
	printf ("%d\n", 6 - 3 - 2)
	f1 ()
	f2 ()
	f3 ()
	f4 ()
	f5 (1, 2)
	f6 (1, 2)
	f7 (1)
}
