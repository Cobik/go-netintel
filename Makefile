SHELL := /bin/bash
.PHONY: tidy run-collector run-ingester compose-up compose-down

tidy:
	go mod tidy

run-collector:
	APP_NAME=collector HTTP_ADDR=:8080 METRICS_ADDR=:9090 \
	KAFKA_BROKERS=localhost:9092 KAFKA_TOPIC=netintel.events.v1 \
	go run ./cmd/collector

run-ingester:
	APP_NAME=ingester KAFKA_BROKERS=localhost:9092 \
	KAFKA_TOPIC=netintel.events.v1 CLICKHOUSE_DSN=clickhouse://default:@localhost:9000/default \
	go run ./cmd/ingester

compose-up:
	cd deploy && docker compose up -d --build

compose-down:
	cd deploy && docker compose down -v
