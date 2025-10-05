package events

import (
	"time"
	"github.com/google/uuid"
)

const SchemaVersion = 1

type NetIntelEvent struct {
	ID         uuid.UUID       `json:"id"`
	Version    int             `json:"version"`
	Source     string          `json:"source"`   // dns/http/whois...
	Subject    string          `json:"subject"`  // domain/host
	ObservedAt time.Time       `json:"observed_at"`
	Meta       map[string]any  `json:"meta,omitempty"`
}

func New(source, subject string, meta map[string]any) NetIntelEvent {
	return NetIntelEvent{
		ID: uuid.New(), Version: SchemaVersion,
		Source: source, Subject: subject, ObservedAt: time.Now().UTC(),
		Meta: meta,
	}
}
