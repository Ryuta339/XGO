#!/bin/bash

go run xgo.go > tmp.s
gcc -o tmp.out tmp.s && ./tmp.out

if [[ $? -eq 0 ]];then
	echo "ok"
fi

rm -rf tmp.out tmp.s
