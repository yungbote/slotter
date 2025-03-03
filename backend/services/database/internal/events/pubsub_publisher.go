package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// PubSubPublisher interface for emitting events to user or company channels.
type PubSubPublisher interface {
	PublishCompanyEvent(companyID uuid.UUID, action string, payload interface{}) error
	PublishUserEvent(userID uuid.UUID, action string, payload interface{}) error
}

// redisPublisher is an implementation using go-redis and JSON messages.
type redisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher initializes the redis client, pings, and returns the publisher.
func NewRedisPublisher(redisURL string) (PubSubPublisher, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
		// Password: "", // if needed
		// DB:       0,  // if needed
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &redisPublisher{client: rdb}, nil
}

// PublishCompanyEvent sends JSON with {company_id, action, payload, timestamp} to "company:UUID".
func (p *redisPublisher) PublishCompanyEvent(companyID uuid.UUID, action string, payload interface{}) error {
	evt := struct {
		CompanyID uuid.UUID   `json:"company_id"`
		Action    string      `json:"action"`
		Payload   interface{} `json:"payload"`
		Timestamp int64       `json:"timestamp"`
	}{
		CompanyID: companyID,
		Action:    action,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal company event: %w", err)
	}
	channelName := fmt.Sprintf("company:%s", companyID.String())

	if err := p.client.Publish(context.Background(), channelName, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to redis: %w", err)
	}
	return nil
}

// PublishUserEvent sends JSON with {user_id, action, payload, timestamp} to "user:UUID".
func (p *redisPublisher) PublishUserEvent(userID uuid.UUID, action string, payload interface{}) error {
	evt := struct {
		UserID    uuid.UUID   `json:"user_id"`
		Action    string      `json:"action"`
		Payload   interface{} `json:"payload"`
		Timestamp int64       `json:"timestamp"`
	}{
		UserID:    userID,
		Action:    action,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal user event: %w", err)
	}
	channelName := fmt.Sprintf("user:%s", userID.String())

	if err := p.client.Publish(context.Background(), channelName, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to redis: %w", err)
	}
	return nil
}

