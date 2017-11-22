package db

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
)

type DB struct {
	Database *sql.DB
	Cfg      Config
	nodeTbl  string
	srs      int
}

func NewDB(cfgToml string) *DB {
	cfg := ReadConfig(cfgToml)
	return &DB{Cfg: cfg}
}

func (db *DB) Open() *DB{
	var d, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		db.Cfg.User, db.Cfg.Password, db.Cfg.Database,
	))
	if err != nil {
		log.Fatalln(err)
	}
	db.Database = d
	return db
}

func (db *DB) Close() *DB{
	db.Database.Close()
	return db
}

func (db *DB) Cleanup() *DB {
	sq := fmt.Sprintf("DROP TABLE IF EXISTS  %v CASCADE;", db.nodeTbl)
	_, err := db.Exec(sq)
	if err != nil {
		log.Fatalln(err)
	}
	db.nodeTbl = ""
	return db
}

func (db *DB) Exec(sql string) (sql.Result, error) {
	return db.Database.Exec(sql)
}

func SQLInsertIntoTable(table string, columns string, values [][]string) string {
	var n = len(values) - 1
	var buf = bytes.Buffer{}
	for i, row := range values {
		buf.WriteString("(" + strings.Join(row, ",") + ")")
		if i < n {
			buf.WriteString(",\n")
		}
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES \n%s;", table, columns, buf.String())
}
