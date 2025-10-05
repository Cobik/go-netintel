package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
	"github.com/yourname/go-netintel/internal/config"
	"github.com/yourname/go-netintel/internal/events"
)

func main() {
	cfg := config.FromEnv()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  "go-netintel-ingester",
		MinBytes: 1, MaxBytes: 10e6,
	})
	defer r.Close()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.CHAddr},
		Auth: clickhouse.Auth{
			Database: cfg.CHDatabase,
			Username: cfg.CHUser,
			Password: cfg.CHPassword,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS netintel_events (
    id UUID, version UInt8, source LowCardinality(String), subject String,
    observed_at DateTime, meta String
	) ENGINE=MergeTree ORDER BY (observed_at, source, subject)`)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var evt events.NetIntelEvent
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Println("bad message:", err)
			continue
		}

		metaBytes, _ := json.Marshal(evt.Meta)
		if err := conn.Exec(ctx, `INSERT INTO netintel_events (id, version, source, subject, observed_at, meta) VALUES (?, ?, ?, ?, ?, ?)`,
			evt.ID.String(), uint8(evt.Version), evt.Source, evt.Subject, evt.ObservedAt.Truncate(time.Second), string(metaBytes)); err != nil {
			log.Println("insert error:", err)
		}
	}
}
