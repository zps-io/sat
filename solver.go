package sat

import (
	"fmt"
	"strings"
	"sync"
)

const conjunction = " ∧ "
const disjunction = " ∨ "

type LiteralEncoder interface {
	IsNegated() uint
	Name() string
}

type Solver struct {
	VerboseDebug bool

	variableLock sync.Mutex

	variableNames        []string
	variableNameToNumber map[string]uint
	variableNumberToName map[uint]string

	clauses [][]uint

	watchList map[uint][][]uint
}

func NewSolver() *Solver {
	return &Solver{
		variableNameToNumber: make(map[string]uint),
		variableNumberToName: make(map[uint]string),
	}
}

func (s *Solver) setUpWatchList() {
	s.watchList = make(map[uint][][]uint)
	for i, clause := range s.clauses {
		if s.VerboseDebug {
			fmt.Printf("i is %d, Clause is: %s", i, clause)
		}
		s.watchList[clause[0]] = append(s.watchList[clause[0]], clause)
	}

	if s.VerboseDebug {
		s.printWatchList()
	}
}

func (s *Solver) printableLiteral(literal uint) string {
	if literal%2 == 0 {
		if varName, ok := s.variableNumberToName[(literal|1)>>1]; ok {
			return fmt.Sprintf("¬%s", varName)
		}
	} else {
		if varName, ok := s.variableNumberToName[literal>>1]; ok {
			return varName
		}
	}
	panic("Problem statement invariant is not maintained")
}

func (s *Solver) componentLiteral(literal uint) (string, bool) {
	if literal%2 == 0 {
		if varName, ok := s.variableNumberToName[(literal|1)>>1]; ok {
			return varName, false
		}
	} else {
		if varName, ok := s.variableNumberToName[literal>>1]; ok {
			return varName, true
		}
	}
	panic("Problem statement invariant is not maintained")
}

func (s *Solver) printableClause(clause []uint) string {
	literals := make([]string, len(clause))
	for j, literal := range clause {
		literals[j] = s.printableLiteral(literal)
	}
	return fmt.Sprintf("(%s)", strings.Join(literals, disjunction))
}

func (s *Solver) printAssignment(assignment []*uint) {
	fmt.Println("Current assignment:")
	literals := make([]string, len(assignment))
	for i, assign := range assignment {
		if assign != nil {
			literals[i] = s.printableLiteral(uint(i)*2 + *assign)
		}
	}
	fmt.Printf("  %s\n", strings.Join(literals, " "))
}

func (s *Solver) printWatchList() {
	fmt.Println("Current watchlist:")
	for literal, w := range s.watchList {
		literalName := s.printableLiteral(uint(literal))
		clauseStrings := make([]string, len(w))
		for i, c := range w {
			clauseStrings[i] = s.printableClause(c)
		}
		clausesString := strings.Join(clauseStrings, ",")

		fmt.Printf("  %4s: %s\n", literalName, clausesString)
	}
}

func (s *Solver) updateWatchList(falseLiteral uint, assignment []*uint) bool {
	for len(s.watchList[falseLiteral]) > 0 {
		clause := s.watchList[falseLiteral][0]
		foundAlternative := false

		for _, alternative := range clause {
			v := alternative >> 1
			a := alternative & 1

			if assignment[v] == nil || *assignment[v] == uint(a^1) {
				foundAlternative = true
				if len(s.watchList[falseLiteral]) > 1 {
					s.watchList[falseLiteral] = s.watchList[falseLiteral][1:]
				} else {
					delete(s.watchList, falseLiteral)
				}
				s.watchList[alternative] = append(s.watchList[alternative], clause)
				break
			}
		}

		if !foundAlternative {
			if s.VerboseDebug {
				s.printWatchList()
				s.printAssignment(assignment)
				fmt.Printf("Clause %s contradicted.\n", s.printableClause(clause))
			}
			return false
		}
	}
	return true
}

func (s *Solver) String() string {
	disjunctions := make([]string, len(s.clauses))
	for i, clause := range s.clauses {
		disjunctions[i] = s.printableClause(clause)
	}

	return strings.Join(disjunctions, conjunction)
}

func (s *Solver) AddClause(literals ...LiteralEncoder) {
	s.variableLock.Lock()
	defer s.variableLock.Unlock()

	clause := make([]uint, len(literals))

	for i, l := range literals {
		varNumber, ok := s.variableNameToNumber[l.Name()]
		if !ok {
			s.variableNames = append(s.variableNames, l.Name())
			varNumber = uint(len(s.variableNames) - 1)
			s.variableNameToNumber[l.Name()] = varNumber
			s.variableNumberToName[varNumber] = l.Name()
		}

		clause[i] = varNumber<<1 | l.IsNegated()
	}

	UintSlice(clause).Sort()

	s.clauses = append(s.clauses, clause)
}

func (s *Solver) solve() (bool, [][]*uint) {
	results := make([][]*uint, 0)

	numberOfVariables := len(s.variableNames)
	assignment := make([]*uint, numberOfVariables)
	state := make([]int, numberOfVariables)

	var d uint = 0

	for {
		if d == uint(numberOfVariables) {
			result := make([]*uint, len(assignment))
			for i, val := range assignment {
				result[i] = val
			}

			results = append(results, result)
			d -= 1
			continue
		}

		tried := false
		for _, a := range []uint{0, 1} {
			if (state[d]>>a)&1 == 0 {
				if s.VerboseDebug {
					fmt.Printf("Trying %s = %d\n", s.variableNames[d], a)
				}
				tried = true
				state[d] |= 1 << a
				assignment[d] = &a

				if !s.updateWatchList(d<<1|a, assignment) {
					assignment[d] = nil
				} else {
					d += 1
					break
				}
			}
		}

		if !tried {
			if d == 0 {
				return len(results) > 0, results
			}

			if s.VerboseDebug {
				fmt.Printf("d: %d\n", d)
			}

			state[d] = 0
			assignment[d] = nil
			d -= 1
		}
	}

	return len(results) > 0, results
}

func (s *Solver) Satisfiable() (bool, []*Solution) {
	var result []*Solution

	s.setUpWatchList()
	satisfiable, solutions := s.solve()

	for _, solution := range solutions {
		so := NewSolution()

		for i := range solution {
			if solution[i] == nil {
				continue
			}

			so.Set(s.componentLiteral(*solution[i] + 2*uint(i)))
		}

		result = append(result, so)
	}

	return satisfiable, result
}
