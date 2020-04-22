#!/bin/bash

s_file="./out/tmp.s"
prog_name="xgo"
go build -o ${prog_name} *.go

function test {
	expected="$2"
	expr="$1"

	echo "$expr" | go run *.go > tmp.s
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

function run_test_go {
	./${prog_name} < ./test/test.go > $s_file
	gcc -o ./out/tmp.out $s_file
	./out/tmp.out > ./out/actual.txt
	diff ./out/actual.txt ./test/expected.txt
}

run_test_go

echo "All tests passed"
