package db

import (
    "github.com/intdxdt/mbr"
    "github.com/intdxdt/rtree"
)

type DBI interface {
    Search(*mbr.MBR) []*rtree.Node
    Insert(rtree.BoxObj) DBI
    Load([]rtree.BoxObj) DBI
    Remove(rtree.BoxObj) DBI
    KNN(
        query rtree.BoxObj,
        limit int,
        score func(rtree.BoxObj, rtree.BoxObj) float64,
        predicates ...func(*rtree.KObj) (bool, bool),
    ) []rtree.BoxObj
}
