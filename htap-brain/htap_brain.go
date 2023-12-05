package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
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

	c chan struct{}
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
	if w.wal, err = wlog.NewSize(nil, nil, "wal", 10*wlog.DefaultSegmentSize, wlog.CompressionNone); err != nil {
		return nil, err
	}
	segment, err := wlog.OpenReadSegment(wlog.SegmentName("wal", 0))
	if err != nil {
		return nil, err
	}
	w.walReaderCloser = segment

	w.walLiveReader = wlog.NewLiveReader(nil, wlog.NewLiveReaderMetrics(nil), segment)

	w.c = make(chan struct{}, 100000)

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
		fmt.Println(err)
		return err
	}

	if err := w.wal.Log([]byte(query)); err != nil {
		return err
	}

	w.c <- struct{}{}

	return nil
}

func (w *HTAPBrain) Query(query string) error {
	olapQueries := map[string]bool{
		"select AVG(user_age) from htap_table": true,
	}
	if olapQueries[query] {
		fmt.Println("Reading clickhouse:", query)
		return w.queryClickhouse(query)
	}
	fmt.Println("Reading postgres:", query)
	return w.queryPostgres(query)
}

func (w *HTAPBrain) writeToClickhouseLoop() error {
	for {
		select {
		case <-w.c:
		}
		if !w.walLiveReader.Next() {
			break
		}
		query := string(w.walLiveReader.Record())

		s := strings.Split(query, "VALUES")
		s = strings.Split(s[1], "),(")
		for i := range s {
			s[i] = strings.Replace(s[i], "(", "", -1)
			s[i] = strings.Replace(s[i], ")", "", -1)
			spl := strings.Split(s[i], ",")
			s[i] = fmt.Sprintf("(%s, %s)", spl[0], spl[1])
		}

		query = fmt.Sprintf(`
			INSERT INTO htap_table
			(user_name, user_age)
			VALUES
			%s	    
		`, strings.Join(s, ","))

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
	}
	return nil
}

func (w *HTAPBrain) queryPostgres(query string) error {
	rows, err := w.psqlClientRead.Query(query)
	if err != nil {
		return err
	}
	for rows.Next() {
	}
	return nil
}

func (w *HTAPBrain) Close() error {
	close(w.c)
	err1 := w.psqlClientWrite.Close()
	err2 := w.chClientWrite.Close()
	if err1 != nil || err2 != nil {
		return fmt.Errorf("postgres error = %q, clickhouse error = %q", err1, err2)
	}
	return w.walReaderCloser.Close()
}

func makePostgresClient() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 dbname=htap sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func makeClickhouseClient() (cdriver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "htap",
			Username: "default",
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return conn, nil
}
