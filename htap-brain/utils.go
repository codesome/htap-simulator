package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	cdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func makePostgresClient() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func makeClickhouseClient() (cdriver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"<CLICKHOUSE_SECURE_NATIVE_HOSTNAME>:9440"},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: "default",
				Password: "<DEFAULT_USER_PASSWORD>",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return conn, nil
}
