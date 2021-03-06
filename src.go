package db

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type DataSrc struct {
	Src    *sql.DB
	Config Config
	SRID   int
	Dim    int
	Table  string
}

type GeomCol struct {
	Table     string
	GeoColumn string
	GeomType  string
	SRID      int
}

func NewDataSrc(configToml string) *DataSrc {
	var cfg = NewConfig(configToml)
	if cfg.Ignore {
		return nil
	}
	var dsrc = &DataSrc{Config: cfg}

	var sqlSrc, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database,
	))
	sqlSrc.SetMaxOpenConns(-1)
	if err == nil {
		dsrc.Src = sqlSrc
		dsrc.Dim = dsrc.CoordDim()
		dsrc.SRID = dsrc.GetSRID()
	} else {
		log.Fatalln(err)
	}

	return dsrc
}

func (dbsrc *DataSrc) Close() *DataSrc {
	if err := dbsrc.Src.Close(); err != nil {
		log.Panic(err)
	}
	return dbsrc
}

func (dbsrc *DataSrc) AlterAsMultiLineString(tableName, geomColumn string, srid int) *DataSrc {
	var query = fmt.Sprintf(`
			ALTER TABLE %v
			ALTER COLUMN %v TYPE geometry(MULTILINESTRING, %v)
			USING ST_Multi(%v);
		`, tableName, geomColumn, srid, geomColumn,
	)

	var _, err = dbsrc.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
	return dbsrc
}

func (dbsrc *DataSrc) DuplicateTable(newTableName string) *DataSrc {
	var sq = fmt.Sprintf(`
		DROP TABLE IF EXISTS %v CASCADE;
		CREATE TABLE %v AS TABLE %v;
		`,
		newTableName, newTableName, dbsrc.Config.Table,
	)
	var _, err = dbsrc.Exec(sq)
	if err != nil {
		log.Fatalln(err)
	}
	return dbsrc
}

func (dbsrc *DataSrc) DeleteTable(table string) *DataSrc {
	var sq = fmt.Sprintf("DROP TABLE IF EXISTS %v CASCADE;", table)
	var _, err = dbsrc.Exec(sq)
	if err != nil {
		log.Fatalln(err)
	}
	return dbsrc
}

func (dbsrc *DataSrc) Exec(query string) (sql.Result, error) {
	return dbsrc.Src.Exec(query)
}

func (dbsrc *DataSrc) Query(query string) (*sql.Rows, error) {
	return dbsrc.Src.Query(query)
}

func (dbsrc *DataSrc) CoordDim() int {
	h, err := dbsrc.Query(fmt.Sprintf(
		`SELECT ST_CoordDim(%v) as dim FROM %v LIMIT 1;`,
		dbsrc.Config.GeometryColumn, dbsrc.Config.Table,
	))

	if err != nil {
		log.Panic(err)
	}
	defer h.Close()

	var dim int
	for h.Next() {
		h.Scan(&dim)
	}
	return dim
}

func (dbsrc *DataSrc) GetSRID() int {
	h, err := dbsrc.Query(fmt.Sprintf(
		`SELECT ST_SRID(%v) as srid FROM %v LIMIT 1;`,
		dbsrc.Config.GeometryColumn, dbsrc.Config.Table,
	))
	if err != nil {
		log.Fatalln(err)
	}
	defer h.Close()

	var srid int
	for h.Next() {
		h.Scan(&srid)
	}
	return srid
}

func SQLInsertIntoTable(table string, columns string, values [][]string) string {
	var n = len(values) - 1
	var buf bytes.Buffer
	if len(values) < 0 {
		log.Fatalln("no values provided")
	}
	var v = values[0]
	var c = strings.Split(columns, ",")
	if len(c) != len(v) {
		log.Fatalln("inconsistent number of columns")
	}

	for i, row := range values {
		buf.WriteString("(" + strings.Join(row, ",") + ")")
		if i < n {
			buf.WriteString(",\n")
		}
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES \n%s;",
		table, columns, buf.String(),
	)
}
