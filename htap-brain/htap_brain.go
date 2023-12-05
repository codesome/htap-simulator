package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	cdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/lib/pq"
	"github.com/prometheus/prometheus/tsdb/wlog"
)

type HTAPBrain struct {
	psqlClientWrite *sql.DB
	chClientWrite   cdriver.Conn
	psqlClientRead  *sql.DB
	chClientRead    cdriver.Conn

	wal             *wlog.WL
	walLiveReader   *wlog.LiveReader
	walReaderCloser io.Closer
}

func NewHTAPBrain() (*HTAPBrain, error) {
	w := &HTAPBrain{}
	var err error
	if w.psqlClientWrite, err = makePostgresClient(); err != nil {
		return nil, err
	}
	if w.chClientWrite, err = makeClickhouseClient(); err != nil {
		return nil, err
	}
	if w.psqlClientRead, err = makePostgresClient(); err != nil {
		return nil, err
	}
	if w.chClientRead, err = makeClickhouseClient(); err != nil {
		return nil, err
	}
	if w.wal, err = wlog.NewSize(nil, nil, "wal", 10*wlog.DefaultSegmentSize, false); err != nil {
		return nil, err
	}
	segment, err := wlog.OpenReadSegment(wlog.SegmentName("wal", 0))
	if err != nil {
		return nil, err
	}
	w.walReaderCloser = segment

	w.walLiveReader = wlog.NewLiveReader(nil, wlog.NewLiveReaderMetrics(nil), segment)

	go func() {
		err := w.writeToClickhouseLoop()
		if err != nil {
			panic(err)
		}
	}()

	return w, nil
}

func (w *HTAPBrain) Write(query string) error {
	if err := w.writeToPostgres(query); err != nil {
		return err
	}

	if err := w.wal.Log([]byte(query)); err != nil {
		return err
	}

	return nil
}

func (w *HTAPBrain) Query(query string) error {
	olapQueries := map[string]bool{
		"Example": true,
	}
	if olapQueries[query] {
		return w.queryClickhouse(query)
	}
	return w.queryPostgres(query)
}

func (w *HTAPBrain) writeToClickhouseLoop() error {
	for w.walLiveReader.Next() {
		query := string(w.walLiveReader.Record())
		err := w.chClientWrite.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *HTAPBrain) writeToPostgres(query string) error {
	_, err := w.psqlClientWrite.Exec(query)
	return err
}

func (w *HTAPBrain) queryClickhouse(query string) error {
	rows, err := w.chClientRead.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(); err != nil {
			panic(err) // TODO: remove this panic
			return err
		}
	}
	return nil
}

func (w *HTAPBrain) queryPostgres(query string) error {
	rows, err := w.psqlClientRead.Query(query)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(); err != nil {
			panic(err) // TODO: remove this panic
			return err
		}
	}
	return nil
}

func (w *HTAPBrain) Close() error {
	err1 := w.psqlClientWrite.Close()
	err2 := w.chClientWrite.Close()
	if err1 != nil || err2 != nil {
		return fmt.Errorf("postgres error = %q, clickhouse error = %q", err1, err2)
	}
	return w.walReaderCloser.Close()
}
