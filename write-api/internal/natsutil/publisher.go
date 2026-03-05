package natsutil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

// Publisher wraps a NATS JetStream connection for publishing messages.
type Publisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

// PoemModerationPayload is the message sent to the moderation service.
type PoemModerationPayload struct {
	PoemID  string `json:"poem_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Connect creates a NATS connection and JetStream context.
// Returns nil Publisher (no-op) if natsURL is empty.
func Connect(natsURL string) (*Publisher, error) {
	if natsURL == "" {
		slog.Warn("NATS_URL not configured — moderation publishing disabled")
		return nil, nil
	}

	nc, err := nats.Connect(natsURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(30),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			slog.Warn("NATS disconnected", "error", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			slog.Info("NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats jetstream: %w", err)
	}

	// Ensure the POEMS stream exists
	_, err = js.StreamInfo("POEMS")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      "POEMS",
			Subjects:  []string{"POEMS.*"},
			Retention: nats.WorkQueuePolicy,
			MaxAge:    24 * time.Hour,
		})
		if err != nil {
			nc.Close()
			return nil, fmt.Errorf("nats create stream: %w", err)
		}
		slog.Info("Created NATS JetStream stream: POEMS")
	}

	slog.Info("Connected to NATS", "url", natsURL)
	return &Publisher{nc: nc, js: js}, nil
}

// PublishModeration sends a poem to the moderation queue.
func (p *Publisher) PublishModeration(payload PoemModerationPayload) error {
	if p == nil {
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal moderation payload: %w", err)
	}

	_, err = p.js.Publish("POEMS.moderate", data)
	if err != nil {
		return fmt.Errorf("publish moderation: %w", err)
	}

	slog.Info("Published poem for moderation", "poem_id", payload.PoemID)
	return nil
}

// Close gracefully drains the NATS connection.
func (p *Publisher) Close() {
	if p == nil || p.nc == nil {
		return
	}
	_ = p.nc.Drain()
	slog.Info("NATS connection closed")
}
