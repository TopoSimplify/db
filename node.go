package db

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"simplex/node"
	"encoding/base64"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/random"
)

//DBNode
type DBNode struct {
	Id       string
	HullType int
	Pln      []*geom.Point
	Range    [2]int
	WTK      string
	Instance string
}


func  CreateNodeTable(db *DB, srid int) error {
	if db.nodeTbl == "" {
		db.nodeTbl = random.String(10)
	}
	var hullSQL = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id SERIAL NOT NULL,
		    gob TEXT NOT NULL,
		    geom GEOMETRY(Geometry, %v),
		    CONSTRAINT pid_%v PRIMARY KEY (id)
		)  WITH (OIDS=FALSE);
	`, db.nodeTbl, srid, db.nodeTbl)
	_, err := db.Exec(hullSQL)
	return err
}


func NewDBNode(node *node.Node) *DBNode {
	return &DBNode{
		Id:       node.Id(),
		HullType: node.Geom.Type().Value(),
		Pln:      node.Polyline.Coordinates,
		Range:    node.Range.AsArray(),
		WTK:      node.Geom.WKT(),
		Instance: node.Instance.Id(),
	}
}

// go binary encoder
func Serialize(n *DBNode) string {
	var buf bytes.Buffer
	var err = gob.NewEncoder(&buf).Encode(n)
	if err != nil {
		log.Fatalln(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// go binary decoder
func Deserialize(str string) *DBNode {
	var n *DBNode
	var dat, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(`failed base64 Decode`, err)
	}
	var buf bytes.Buffer
	_, err = buf.Write(dat)
	if err != nil {
		log.Fatalln(`failed to write to buffer`)
	}
	err = gob.NewDecoder(&buf).Decode(&n)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return n
}
