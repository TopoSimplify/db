package db

import (
	"fmt"
	"time"
	"math/rand"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
    rand.Seed(time.Now().UnixNano())
}

func randTableName(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

func CreateNodeTable(src *DataSrc) error {
	var err error
	var tbl = src.NodeTable
	if src.NodeTable == "" {
		src.NodeTable = randTableName(10)
		var hullSQL = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		    id  SERIAL NOT NULL,
		    fid INT NOT NULL,
		    gob TEXT NOT NULL,
		    geom GEOMETRY(Geometry, %v) NOT NULL,
		    status int DEFAULT 0,
		    CONSTRAINT pid_%v PRIMARY KEY (id)
		) WITH (OIDS=FALSE);
		CREATE INDEX idx_status_%v ON %v (status);
		CREATE INDEX gidx_%v ON %v USING GIST (geom);
	`, tbl, src.SRID, tbl, tbl, tbl, tbl, tbl)
		_, err = src.Exec(hullSQL)
	}
	return err
}

