package main

/*** interface definitioins ***/
type ArithmeticOperator interface {
	emitOperator()
}

/* ===============================
 * Arithmetic operators implementation
 * =============================== */
type AdditiveOperator struct {
}

// implements ArithmeticOperator
func (ao *AdditiveOperator) emitOperator() {
	emitCode("\taddl\t%%ebx, %%eax")
}

type SubtractionOperator struct {
}

// implements ArithmeticOperator
func (so *SubtractionOperator) emitOperator() {
	emitCode("\tsubl\t%%ebx, %%eax")
}

type MultiplicativeOperator struct {
}

// implements ArithmeticOperator
func (mo *MultiplicativeOperator) emitOperator() {
	emitCode("\tpushq\t%%rdx")
	emitCode("\timul\t%%ebx, %%eax")
	emitCode("\tpopq\t%%rdx")
}

type DivisionOperator struct {
}

// implements AritheticOperator
func (do *DivisionOperator) emitOperator() {
	emitCode("\tidivl\t%%ebx, %%eax")
}
