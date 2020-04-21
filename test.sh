#!/bin/bash


function test {
	expected="$2"
	expr="$1"

	echo -n "$expr" | go run *.go > tmp.s
	# gcc -o tmp.out driver.c tmp.s 
	gcc -o tmp.out tmp.s
	result="`./tmp.out`"

	if [[ "$result" -eq "$expected" ]];then
		echo "ok"
		rm -rf tmp.out tmp.s
	else
		echo "Test failed: $expected expected but got $result"
		exit 1
	fi
}


test 'printf ("%d\\n", 0)' 0 
test 'printf ("%d\\n", 7)' 7
test 'printf ("%d\\n", 2 + 5)' 7
test 'printf ("%d\\n", 10 - 4)' 6
test 'printf ("%d\\n", 4 * 3)' 12
test 'printf ("%d\\n", 1 * 2 + 3 * 4)' 14
test 'printf ("%d\\n", 1 + 2 * 3 + 4)' 11
test 'printf ("%d\\n", 6 - 3 - 2)' 1

echo "All tests passed"
