package main

import "fmt"
import (
	"strings"
	"strconv"
)

func f () {
	puts ("")
	puts ("hello")
	puts ("world")
}

func f2 () {
	var i int
	i = 3
	printf ("%d\n", i)
}

func f3 () {
	var j int = 5
	var k int = 6
	printf ("%d\n", j + k)
}


func main () {
	printf ("%d\n", 2 + 5)
	printf ("%d\n", 10 - 4)
	printf ("%d\n", 4 * 3)
	printf ("%d\n", 1 * 2 + 3 * 4)
	printf ("%d\n", 1 + 2 * 3 + 4)
	printf ("%d\n", 6 - 3 - 2)
	f ()
	f2 ()
	f3 ()
}
