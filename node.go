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
	"simplex/pln"
	"simplex/rng"
)

const nodeTblColumns = "gob, geom"

//DBNode
type DBNode struct {
	Id       string
	HullType int
	Pln      []*geom.Point
	Range    [2]int
	WTK      string
}

func CreateNodeTable(src *DataSrc) error {
	var err error
	if src.NodeTable == "" {
		src.NodeTable = random.String(10)
		var hullSQL = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id SERIAL NOT NULL,
		    gob TEXT NOT NULL,
		    geom GEOMETRY(Geometry, %v) NOT NULL,
		    status int DEFAULT 0,
		    CONSTRAINT pid_%v PRIMARY KEY (id)
		) WITH (OIDS=FALSE);
		CREATE INDEX %v_gidx ON %v USING GIST (geom);
	`, src.NodeTable, src.SRID, src.NodeTable, src.NodeTable, src.NodeTable)
		_, err = src.Exec(hullSQL)
	}
	return err
}

func NewDBNode(node *node.Node) *DBNode {
	return &DBNode{
		Id:       node.Id(),
		HullType: node.Geom.Type().Value(),
		Pln:      node.Polyline.Coordinates,
		Range:    node.Range.AsArray(),
		WTK:      node.Geom.WKT(),
	}
}

func NewNodeFromDB(ndb *DBNode) *node.Node {
	var n = &node.Node{
		Polyline: pln.New(ndb.Pln),
		Range:    rng.NewRange(ndb.Range[0], ndb.Range[1]),
		Geom:     geom.NewGeometry(ndb.WTK),
	}
	n.SetId(ndb.Id)
	return n
	//return &DBNode{
	//	Id:       node.Id(),
	//	HullType: node.Geom.Type().Value(),
	//	Pln:      node.Polyline.Coordinates,
	//	Range:    node.Range.AsArray(),
	//	WTK:      node.Geom.WKT(),
	//}
}

func BulkLoadNodes(src *DataSrc, nodes []*node.Node) error {
	var vals = make([][]string, 0)
	for _, h := range nodes {
		vals = append(vals, []string{
			fmt.Sprintf(`'%v'`, Serialize(NewDBNode(h))),
			fmt.Sprintf(`ST_GeomFromText('%v', %v)`, h.Geom.WKT(), src.SRID),
		})
	}
	_, err := src.Exec(SQLInsertIntoTable(src.NodeTable, nodeTblColumns, vals))
	return err
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
