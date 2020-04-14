#!/bin/bash

function test {
	expected="$2"
	expr="$1"

	echo -n "$expr" | go run xgo.go token.go ast.go > tmp.s
	gcc -o tmp.out driver.c tmp.s 
	result="`./tmp.out`"

	if [[ "$result" -eq "$expected" ]];then
		echo "ok"
		rm -rf tmp.out tmp.s
	else
		echo "Test failed: $expected expected but got $result"
		exit 1
	fi
}


test 0 0 
test 7 7
test '2 + 5' 7
test '10 - 4' 6
test '4 * 3' 12

echo "All tests passed"
