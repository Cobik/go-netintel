package storage

import (
	"context"
	"time"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickDB struct{ conn clickhouse.Conn }

func NewClick(dsn string) (*ClickDB, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{Addr: []string{dsn}})
	if err != nil { return nil, err }
	return &ClickDB{conn: conn}, nil
}

func (db *ClickDB) Init(ctx context.Context) error {
	q := `
	CREATE TABLE IF NOT EXISTS netintel_events (
		id UUID,
		version UInt8,
		source LowCardinality(String),
		subject String,
		observed_at DateTime,
		meta JSON
	) ENGINE=MergeTree
	ORDER BY (observed_at, source, subject)`
	return db.conn.Exec(ctx, q)
}

func (db *ClickDB) InsertJSON(ctx context.Context, id string, version int8, source, subject string, observed time.Time, meta string) error {
	const q = `INSERT INTO netintel_events (id, version, source, subject, observed_at, meta) VALUES (?, ?, ?, ?, ?, ?)`
	return db.conn.Exec(ctx, q, id, version, source, subject, observed, meta)
}
