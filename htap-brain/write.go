package main

import (
	"database/sql"

	cdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/lib/pq"
)

type Writer struct {
	postgresClient   *sql.DB
	clickhouseClient cdriver.Conn
}

func NewWriter() (*Writer, error) {
	w := &Writer{}
	var err error
	if w.postgresClient, err = makePostgresClient(); err != nil {
		return nil, err
	}
	if w.clickhouseClient, err = makeClickhouseClient(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) writeToClickhouse(query string) {

}

func (w *Writer) writeToPostgres(query string) {

}
