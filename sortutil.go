package sat

import "sort"

// type UintSlice is used to implement in-place sorting of slices of uint
type UintSlice []uint

func (p UintSlice) Len() int           { return len(p) }
func (p UintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p UintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (s UintSlice) Sort() {
	sort.Sort(s)
}
