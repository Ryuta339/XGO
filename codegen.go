package main

import "fmt"

var frameHeight int = 0;

func emitCode (code string, v ...interface{}) {
	fmt.Printf (code+"\n", v...)
}

func emitDataSection () {
	emitCode (".data")

	// put dummy label
	emitCode (".L0:")
	emitCode (".string \"%%d\\n\"")
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

	// call printf ("%d\n", expr)
	emitCode ("\tpushq\t%%rdi")
	emitCode ("\tpushq\t%%rsi")
	frameHeight += 16

	ast.emit ()

	emitCode ("\tlea\t.L0(%%rip), %%rdi")
	emitCode ("\tpopq\t%%rsi")
	frameHeight -= 8
	emitCode ("\tmovq\t$0, %%rax")
	emitCode ("\tcall\t_printf")
	emitCode ("\tpopq\t%%rsi")
	emitCode ("\tpopq\t%%rdi")
	frameHeight -= 16
	
	emitCode ("\tmovl\t$0, %%eax")  // return 0

	emitFuncMainEpilogue ()
}
