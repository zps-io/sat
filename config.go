package sat

import (
	"fmt"
	"sort"
)

type Variable struct {
	name string
}

func (v *Variable) IsNegated() uint {
	return 0
}

func (v *Variable) Name() string {
	return v.name
}

func NewVariable(name string) *Variable {
	return &Variable{
		name: name,
	}
}

type NegatedVariable struct {
	name string
}

func (v *NegatedVariable) IsNegated() uint {
	return 1
}

func (v *NegatedVariable) Name() string {
	return v.name
}

func (v *Variable) Not() LiteralEncoder {
	return &NegatedVariable{
		name: v.name,
	}
}

type Solution map[string]bool

func (s Solution) Print() {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("\t%s = %t", key, s[key])
	}
}

func NewSolution() *Solution {
	return &Solution{}
}

func (s Solution) Set(variable string, value bool) {
	s[variable] = value
}

func (s Solution) Value(variable string) bool {
	return s[variable]
}
