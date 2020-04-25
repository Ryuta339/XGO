package main

import "fmt"

var frameHeight int = 8

func emitCode(code string, v ...interface{}) {
	fmt.Printf(code+"\n", v...)
}

func emitDataSection() {
	emitCode(".data")

	// put stinrgs first
	for _, ast := range stringList {
		emitCode(".%s:", ast.slabel)
		emitCode(".string \"%s\"", ast.sval)
	}
}

func emitFuncPrologue(fname string) {
	// これ後で修正したい
	emitCode(".text")
	emitCode("\t.global _%s", fname)
	emitCode("_%s:", fname)
	emitCode("\tpushq\t%%rbp")
	emitCode("\tmovq\t%%rsp, %%rbp")
	frameHeight += 8
}

func emitFuncEpilogue() {
	emitCode("\tleave")
	emitCode("\tret")
}

func generate(ast Ast) {
	emitDataSection()
	ast.emit()
}
