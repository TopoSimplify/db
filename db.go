//@formatter:off
package db

import (
	"sync"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/rtree"
)


type DB struct {
	sync.RWMutex
	tree *rtree.RTree
}

func NewDB(nodeCapacity int) *DB {
	return &DB{tree : rtree.NewRTree(nodeCapacity)}
}

func (db *DB) Search(bbox *mbr.MBR) []*rtree.Node {
	db.RLock(); defer db.RUnlock()
	return db.tree.Search(bbox)
}

func (db *DB) All() []*rtree.Node {
	db.RLock(); defer db.RUnlock()
	return db.tree.All()
}

func (db *DB) Insert(item rtree.BoxObj) *DB {
	db.Lock(); defer db.Unlock()
	db.tree.Insert(item)
	return db
}

func (db *DB) Load(data []rtree.BoxObj) *DB {
	db.Lock(); defer db.Unlock()
	db.tree.Load(data)
	return db
}

func (db *DB) Remove(item rtree.BoxObj) *DB {
	db.Lock(); defer db.Unlock()
	db.tree.Remove(item)
	return db
}

func (db *DB) Clear() *DB {
	db.Lock(); defer db.Unlock()
	db.tree.Clear()
	return db
}

func (db *DB) KNN (
	query rtree.BoxObj, limit int,
	score func(rtree.BoxObj, rtree.BoxObj) float64,
	predicates ...func(*rtree.KObj) (bool, bool)) []rtree.BoxObj {

	db.RLock(); defer db.RUnlock()

	if len(predicates) > 0 {
	 	return db.tree.KNN(query, limit, score, predicates[0])
	}
	return db.tree.KNN(query, limit, score)
}
