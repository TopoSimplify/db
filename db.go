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
    return &DB{tree: rtree.NewRTree(nodeCapacity)}
}

func (db *DB) Search(bbox *mbr.MBR) []*rtree.Node {
    db.RLock()
    o := db.tree.Search(bbox)
    db.RUnlock()
    return o
}

func (db *DB) All() []*rtree.Node {
    db.RLock()
    o := db.tree.All()
    db.RUnlock()
    return o
}

func (db *DB) KNN(
    query rtree.BoxObj, limit int,
    score func(rtree.BoxObj, rtree.BoxObj) float64,
    predicates ...func(*rtree.KObj) (bool, bool)) []rtree.BoxObj {

    db.RLock()
    if len(predicates) > 0 {
        return db.tree.KNN(query, limit, score, predicates[0])
    }
    o := db.tree.KNN(query, limit, score)
    db.RUnlock()
    return o
}

func (db *DB) Insert(item rtree.BoxObj) *DB {
    db.Lock()
    db.tree.Insert(item)
    db.Unlock()
    return db
}

func (db *DB) Load(data []rtree.BoxObj) *DB {
    db.Lock()
    db.tree.Load(data)
    db.Unlock()
    return db
}

func (db *DB) Remove(item rtree.BoxObj) *DB {
    db.Lock()
    db.tree.Remove(item)
    db.Unlock()
    return db
}

func (db *DB) Clear() *DB {
    db.Lock()
    db.tree.Clear()
    db.Unlock()
    return db
}
