#!/bin/bash

function test {
	expected="$1"
	expr="$2"

	echo -n "$expr" | go run xgo.go > tmp.s
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

echo "All tests passed"
