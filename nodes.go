package db

import "sort"

type Nodes []*Node

func (ns Nodes) Len() int           { return len(ns) }
func (ns Nodes) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns Nodes) Less(i, j int) bool { return ns[i].Range.I < ns[j].Range.I }
func (ns Nodes) Sort() {
	sort.Sort(ns)
}
