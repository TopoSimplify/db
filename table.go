package db

import (
	"fmt"
	"log"
	"time"
	"bytes"
	"math/rand"
	"text/template"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var nodeTblTemplate = `
CREATE TABLE IF NOT EXISTS {{.NodeTable}} (
    id  SERIAL NOT NULL,
    i INT NOT NULL,
    j INT NOT NULL,
    size INT CHECK (size > 0),
    fid INT NOT NULL,
    gob TEXT NOT NULL,
    geom GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    status INT DEFAULT 0,
    CONSTRAINT pid_{{.NodeTable}} PRIMARY KEY (id)
) WITH (OIDS=FALSE);
CREATE INDEX idx_i_{{.NodeTable}} ON {{.NodeTable}} (i);
CREATE INDEX idx_j_{{.NodeTable}} ON {{.NodeTable}} (j);
CREATE INDEX idx_size_{{.NodeTable}} ON {{.NodeTable}} (size);
CREATE INDEX idx_fid_{{.NodeTable}} ON {{.NodeTable}} (fid);
CREATE INDEX idx_status_{{.NodeTable}} ON {{.NodeTable}} (status);
CREATE INDEX gidx_{{.NodeTable}} ON {{.NodeTable}} USING GIST (geom);
`


var nodeTemplate *template.Template
var geomColTemplate *template.Template

func init() {
	rand.Seed(time.Now().UnixNano())
	var err error
	nodeTemplate, err = template.New("node_table").Parse(nodeTblTemplate)
	if err != nil {
		log.Fatalln(err)
	}
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
	if tbl == "" {
		src.NodeTable = fmt.Sprintf("node_%v", randTableName(10))
		tbl = src.NodeTable
		var query bytes.Buffer
		err = nodeTemplate.Execute(&query, src)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = src.Exec(query.String())
	}
	return err
}
