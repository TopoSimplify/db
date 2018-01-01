package db

import (
	"log"
	"fmt"
	"bytes"
	"text/template"
)
const IdColumn = "id"
const GeomColumn = "geom"

var onlineTblTemplate = `
CREATE TABLE IF NOT EXISTS {{.Table}} (
    id          SERIAL NOT NULL,
    fid         INT NOT NULL,
    node        TEXT NOT NULL,
    geom        GEOMETRY(Geometry, {{.SRID}}) NOT NULL,
    i           INT NOT NULL,
    j           INT NOT NULL,
    size        INT CHECK (size > 0),
    status      INT DEFAULT 0,
    snapshot    INT DEFAULT 0,
    CONSTRAINT  pid_{{.Table}}  PRIMARY KEY (id),
	CONSTRAINT  u_constraint    UNIQUE  (fid, i, j)
) WITH (OIDS=FALSE);
CREATE INDEX idx_i_{{.Table}} ON {{.Table}} (i);
CREATE INDEX idx_j_{{.Table}} ON {{.Table}} (j);
CREATE INDEX idx_size_{{.Table}} ON {{.Table}} (size);
CREATE INDEX idx_status_{{.Table}} ON {{.Table}} (status);
CREATE INDEX idx_fid_{{.Table}} ON {{.Table}} (fid);
CREATE INDEX gidx_{{.Table}} ON {{.Table}} USING GIST (geom);
`

var onlineTemplate *template.Template

func init() {
	var err error
	onlineTemplate, err = template.New("online_table").Parse(onlineTblTemplate)
	if err != nil {
		log.Fatalln(err)
	}
}

func CreateNodeTable(Src *DataSrc) error {
	var query bytes.Buffer
	if err := onlineTemplate.Execute(&query, Src);  err != nil {
		log.Fatalln(err)
	}
	var tblSQl = fmt.Sprintf(`DROP TABLE IF EXISTS %v CASCADE;`, Src.Table)
	if _, err := Src.Exec(tblSQl); err != nil {
		log.Panic(err)
	}
	_, err := Src.Exec(query.String())
	return err
}
