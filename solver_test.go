package sat

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureDir = "./test-fixtures"

func ExampleSimpleSatisfiableSolutions() {
	variableA := NewVariable("A")
	negatedVariableA := variableA.Not()

	solver := NewSolver()
	solver.AddClause(variableA, negatedVariableA)

	satisfiable, solutions := solver.Satisfiable()

	fmt.Printf("Satisfiable?: %t\nSolutions:", satisfiable)
	for i, soln := range solutions {
		fmt.Printf("\n%d:\n", i)
		soln.Print()
	}
	// Output: Satisfiable?: true
	// Solutions:
	// 0:
	// 	A = false
	// 1:
	// 	A = true
}

func ExampleSimpleSatisfiableSolutionsOr() {
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
	// Output: Satisfiable?: true
	// Solutions:
	// 0:
	// 	A = false	B = true
	// 1:
	// 	A = true	B = false
	// 2:
	// 	A = true	B = true
}

func ExampleSimpleSatisfiableDependency() {
	variableA := NewVariable("A")

	variableB1 := NewVariable("B1")
	variableB2 := NewVariable("B2")
	variableB3 := NewVariable("B3")

	solver := NewSolver()
	solver.AddClause(variableA)
	solver.AddClause(variableA.Not(), variableB1, variableB2, variableB3)
	solver.AddClause(variableB1.Not(), variableB2.Not())
	solver.AddClause(variableB1.Not(), variableB3.Not())
	solver.AddClause(variableB2.Not(), variableB3.Not())

	satisfiable, solutions := solver.Satisfiable()

	fmt.Printf("Satisfiable?: %t\nSolutions:", satisfiable)
	for i, soln := range solutions {
		fmt.Printf("\n%d:\n", i)
		soln.Print()
	}
	// Output: Satisfiable?: true
	// Solutions:
	// 0:
	// 	A = true	B1 = false	B2 = false	B3 = true
	// 1:
	// 	A = true	B1 = false	B2 = true	B3 = false
	// 2:
	//	A = true	B1 = true	B2 = false	B3 = false
}

func ExampleSimpleNonSatisfiable() {
	variableA := NewVariable("A")
	negatedVariableA := variableA.Not()

	solver := NewSolver()
	solver.AddClause(variableA)
	solver.AddClause(negatedVariableA)

	satisfiable, _ := solver.Satisfiable()

	fmt.Printf("Satisfiable?: %t", satisfiable)
	// Output: Satisfiable?: false
}

func ExampleSimpleSatisfiable() {
	variableA := NewVariable("A")
	variableB := NewVariable("B")

	solver := NewSolver()
	solver.AddClause(variableA, variableB.Not())

	satisfiable, _ := solver.Satisfiable()

	fmt.Printf("Satisfiable?: %t", satisfiable)
	// Output: Satisfiable?: true
}

func TestSATSolverCases_simple(t *testing.T) {
	cases := []struct {
		filename    string
		satisfiable bool
	}{
		{
			"simple/01.txt",
			true,
		},
		{
			"simple/02.txt",
			true,
		},
		{
			"simple/03.txt",
			true,
		},
		{
			"simple/04.txt",
			false,
		},
	}

	for _, tc := range cases {
		solver, err := configureFromTestFixture(t, tc.filename)
		if err != nil {
			t.Fatalf("Error configuring test: %s", err)
		}

		if satisfiable, _ := solver.Satisfiable(); satisfiable != tc.satisfiable {
			t.Fatalf("Expected: %t, Got %t for problem: \n%s", tc.satisfiable, satisfiable, solver)
		}
	}
}

func TestSATSolverCases_colouring(t *testing.T) {
	cases := []struct {
		filename    string
		satisfiable bool
	}{
		{
			"colouring/01.txt",
			true,
		},
	}

	for _, tc := range cases {
		solver, err := configureFromTestFixture(t, tc.filename)
		if err != nil {
			t.Fatalf("Error configuring test: %s", err)
		}

		if satisfiable, _ := solver.Satisfiable(); satisfiable != tc.satisfiable {
			t.Fatalf("Expected: %t, Got %t for problem: \n%s", tc.satisfiable, satisfiable, solver)
		}
	}
}

func TestSATSolverCases_fsm(t *testing.T) {
	cases := []struct {
		filename    string
		satisfiable bool
	}{
		{
			"fsm/even-ones-3.txt",
			true,
		},
		{
			"fsm/even-ones-4.txt",
			true,
		},
		{
			"fsm/even-zeros-3.txt",
			true,
		},
		{
			"fsm/even-zeros-4.txt",
			true,
		},
	}

	for _, tc := range cases {
		solver, err := configureFromTestFixture(t, tc.filename)
		if err != nil {
			t.Fatalf("Error configuring test: %s", err)
		}

		if satisfiable, _ := solver.Satisfiable(); satisfiable != tc.satisfiable {
			t.Fatalf("Expected: %t, Got %t for problem: \n%s", tc.satisfiable, satisfiable, solver)
		}
	}
}

func TestSATSolverCases_w44(t *testing.T) {
	cases := []struct {
		filename    string
		satisfiable bool
	}{
		{
			"w44/w44-008.txt",
			true,
		},
		{
			"w44/w44-009.txt",
			true,
		},
		{
			"w44/w44-010.txt",
			true,
		},
		{
			"w44/w44-011.txt",
			true,
		},
		{
			"w44/w44-012.txt",
			true,
		},
		{
			"w44/w44-013.txt",
			true,
		},
		{
			"w44/w44-014.txt",
			true,
		},
		{
			"w44/w44-015.txt",
			true,
		},
		{
			"w44/w44-016.txt",
			true,
		},
		{
			"w44/w44-017.txt",
			true,
		},
		{
			"w44/w44-018.txt",
			true,
		},
		{
			"w44/w44-019.txt",
			true,
		},
		{
			"w44/w44-020.txt",
			true,
		},
		{
			"w44/w44-021.txt",
			true,
		},
		{
			"w44/w44-022.txt",
			true,
		},
		{
			"w44/w44-023.txt",
			true,
		},
		{
			"w44/w44-024.txt",
			true,
		},
		{
			"w44/w44-025.txt",
			true,
		},
		{
			"w44/w44-026.txt",
			true,
		},
		{
			"w44/w44-027.txt",
			true,
		},
		{
			"w44/w44-028.txt",
			true,
		},
		{
			"w44/w44-029.txt",
			true,
		},
		{
			"w44/w44-030.txt",
			true,
		},
		{
			"w44/w44-031.txt",
			true,
		},
		{
			"w44/w44-032.txt",
			true,
		},
		{
			"w44/w44-033.txt",
			true,
		},
		{
			"w44/w44-034.txt",
			true,
		},
		{
			"w44/w44-035.txt",
			false,
		},
		{
			"w44/w44-036.txt",
			false,
		},
		{
			"w44/w44-037.txt",
			false,
		},
		{
			"w44/w44-038.txt",
			false,
		},
		{
			"w44/w44-039.txt",
			false,
		},
	}

	for i, tc := range cases {
		solver, err := configureFromTestFixture(t, tc.filename)
		if err != nil {
			t.Fatalf("Error configuring test: %s", err)
		}

		if satisfiable, _ := solver.Satisfiable(); satisfiable != tc.satisfiable {
			t.Fatalf("Expected: %t, Got %t for problem %d: \n%s", tc.satisfiable, satisfiable, i, solver)
		}
	}
}

func configureFromTestFixture(t *testing.T, name string) (*Solver, error) {
	fileName := filepath.Join(fixtureDir, name)
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Error opening %s: %s", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	satSolver := NewSolver()
	for scanner.Scan() {
		text := scanner.Text()

		//Ignore comments
		if strings.HasPrefix(text, "#") {
			continue
		}

		// Ignore empty lines
		if strings.TrimSpace(text) == "" {
			continue
		}

		literalStrings := strings.Split(text, " ")
		clauseLiteralMakers := make([]LiteralEncoder, 0)
		for _, literal := range literalStrings {
			trimmedLiteral := strings.TrimSpace(literal)
			if trimmedLiteral == "" {
				continue
			}

			if strings.HasPrefix(trimmedLiteral, "~") {
				variable := NewVariable(trimmedLiteral[1:])
				clauseLiteralMakers = append(clauseLiteralMakers, variable.Not())
			} else {
				clauseLiteralMakers = append(clauseLiteralMakers, NewVariable(trimmedLiteral))
			}
		}

		satSolver.AddClause(clauseLiteralMakers...)
	}

	return satSolver, nil
}
