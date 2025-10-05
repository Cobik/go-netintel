package config

import (
	"os"
	"strings"
)

type Config struct {
	AppName      string
	HTTPAddr     string
	MetricsAddr  string
	KafkaBrokers []string
	KafkaTopic   string
	CHDSN        string // clickhouse://default:@clickhouse:9000/default
}

func getEnv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

func FromEnv() Config {
	return Config{
		AppName:      getEnv("APP_NAME", "go-netintel"),
		HTTPAddr:     getEnv("HTTP_ADDR", ":8080"),
		MetricsAddr:  getEnv("METRICS_ADDR", ":9090"),
		KafkaBrokers: func() []string { v := getEnv("KAFKA_BROKERS", "kafka:9092"); return strings.Split(v, ",") }(),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "netintel.events.v1"),
		CHDSN:        getEnv("CLICKHOUSE_DSN", "clickhouse://default:@clickhouse:9000/default"),
	}
}
