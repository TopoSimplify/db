package db

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
)

type DataSrc struct {
	Src       *sql.DB
	Config    Config
	SRID      int
	NodeTable string
}

func NewDataSrc(configToml string) *DataSrc {
	var cfg = NewConfig(configToml)
	var dsrc = &DataSrc{Config: cfg}
	var sqlsrc, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database,
	))

	if err == nil {
		dsrc.Src = sqlsrc
		dsrc.SRID = dsrc.GetSRID()
	} else {
		log.Fatalln(err)
	}

	return dsrc
}

func (db *DataSrc) Close() *DataSrc {
	db.Src.Close()
	return db
}

func (db *DataSrc) DeleteTable(table string) *DataSrc {
	var sq = fmt.Sprintf("DROP TABLE IF EXISTS %v CASCADE;", table)
	var _, err = db.Exec(sq)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func (db *DataSrc) Exec(query string) (sql.Result, error) {
	return db.Src.Exec(query)
}

func (db *DataSrc) Query(query string) (*sql.Rows, error) {
	return db.Src.Query(query)
}

func (db *DataSrc) GetSRID() int {
	rows, err := db.Query(fmt.Sprintf(
		`SELECT st_srid (%v) as srid FROM %v LIMIT 1;`,
		db.Config.GeometryColumn, db.Config.Table,
	))
	if err != nil {
		log.Fatalln(err)
	}
	var srid int
	for rows.Next() {
		rows.Scan(&srid)
	}
	return srid
}

func SQLInsertIntoTable(table string, columns string, values [][]string) string {
	var n = len(values) - 1
	var buf bytes.Buffer

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
