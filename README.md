sat: a simple sat solver for golang
===================================

A simple SAT solver that utilizes an iterative algorithm. The algorithm uses a watchlist to keep track of all the assignments.

Most of the code is based on the Python implementation by sahands, which can be downloaded [here](https://github.com/sahands/simple-sat). Some checking was done utilizing another port found [here](https://raw.githubusercontent.com/marcvanzee/go-sat).

The primary motivation for this port was to create a code base that could support the future work and requirements of [ZPS](https://github.com/solvent-io/zps) while also maintaining a generic interface that others may find useful.

### Installing

go get github.com/fezz-io/sat

### Basic Usage

```
import (
	"github.com/solvent-io/sat"
)

variableA := NewVariable("A")
variableB := NewVariable("B")

solver := NewSolver()
solver.AddClause(variableA, variableB)

satisfiable, solutions := solver.Satisfiable()

fmt.Printf("Satisfiable?: %t\nSolutions:", satisfiable)
for i, soln := range solutions {
	fmt.Printf("\n%d:\n", i)
	soln.Print()
}
```
