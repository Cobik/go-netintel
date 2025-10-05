package config

import (
	"os"
	"strings"
)

type Config struct {
	AppName                                string
	HTTPAddr                               string
	MetricsAddr                            string
	KafkaBrokers                           []string
	KafkaTopic                             string
	CHAddr, CHUser, CHPassword, CHDatabase string
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func FromEnv() Config {
	return Config{
		AppName:      getEnv("APP_NAME", "go-netintel"),
		HTTPAddr:     getEnv("HTTP_ADDR", ":8080"),
		MetricsAddr:  getEnv("METRICS_ADDR", ":9090"),
		KafkaBrokers: func() []string { v := getEnv("KAFKA_BROKERS", "redpanda:9092"); return strings.Split(v, ",") }(),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "netintel.events.v1"),
		CHAddr:       getEnv("CLICKHOUSE_ADDR", "clickhouse:9000"),
		CHUser:       getEnv("CH_USER", "netintel"),
		CHPassword:   getEnv("CH_PASSWORD", ""),
		CHDatabase:   getEnv("CH_DATABASE", "default"),
	}
}
