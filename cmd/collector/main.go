package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/yourname/go-netintel/internal/config"
	"github.com/yourname/go-netintel/internal/httpserver"
	"github.com/yourname/go-netintel/internal/queue"
)

func main() {
	cfg := config.FromEnv()
	pub := queue.NewKafkaPublisher(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer pub.Close()

	s := httpserver.New(cfg.HTTPAddr, pub)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := s.Start(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown")
	}
}
