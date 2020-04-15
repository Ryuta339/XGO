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
	emitCode ("pushq\t%%rbp")
	emitCode ("movq\t%%rsp, %%rbp")
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
	emitCode ("pushq\t%%rdi")
	emitCode ("pushq\t%%rsi")
	frameHeight += 16

	ast.emit ()

	emitCode ("lea\t.L0(%%rip), %%rdi")
	emitCode ("popq\t%%rsi")
	frameHeight -= 8
	emitCode ("movq\t$0, %%rax")
	emitCode ("call\t_printf")
	emitCode ("popq\t%%rsi")
	emitCode ("popq\t%%rdi")
	frameHeight -= 16
	
	emitCode ("movl\t$0, %%eax")  // return 0

	emitFuncMainEpilogue ()
}
