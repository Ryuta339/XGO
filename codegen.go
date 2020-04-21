package main

import "fmt"

var frameHeight int = 8;

func emitCode (code string, v ...interface{}) {
	fmt.Printf (code+"\n", v...)
}

func emitDataSection () {
	emitCode (".data")

	// put stinrgs first
	for _, ast := range stringList {
		emitCode (".%s:", ast.slabel)
		emitCode (".string \"%s\"", ast.sval)
	}
}

func emitFuncMainPrologue () {
	// これ後で修正したい
	emitCode (".text")
	emitCode ("\t.global _main")
	emitCode ("_main:");
	emitCode ("\tpushq\t%%rbp")
	emitCode ("\tmovq\t%%rsp, %%rbp")
	frameHeight += 8
}

func emitFuncMainEpilogue () {
	emitCode ("\tleave")
	emitCode ("\tret")
}

func generate (ast Ast) {
	emitDataSection ()
	emitFuncMainPrologue ()

	ast.emit ()

	emitCode ("\tmovl\t$0, %%eax")  // return 0
	emitFuncMainEpilogue ()
}
