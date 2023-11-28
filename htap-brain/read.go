package main

import (
	"database/sql"

	cdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/lib/pq"
)

type Querier struct {
	postgresClient   *sql.DB
	clickhouseClient cdriver.Conn
}

func NewQuerier() (*Querier, error) {
	q := &Querier{}
	var err error
	if q.postgresClient, err = makePostgresClient(); err != nil {
		return nil, err
	}
	if q.clickhouseClient, err = makeClickhouseClient(); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Querier) queryClickhouse(query string) {

}

func (q *Querier) queryPostgres(query string) {

}
