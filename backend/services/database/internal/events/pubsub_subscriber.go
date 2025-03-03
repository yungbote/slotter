package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// PubSubSubscriber allows subscribing to user:xxx or company:xxx channels in Redis.
type PubSubSubscriber struct {
	client              *redis.Client
	mu                  sync.Mutex
	userSubscriptions   map[uuid.UUID]*redis.PubSub
	companySubscriptions map[uuid.UUID]*redis.PubSub
}

// NewPubSubSubscriber constructs a PubSubSubscriber with the given Redis client.
func NewPubSubSubscriber(client *redis.Client) *PubSubSubscriber {
	return &PubSubSubscriber{
		client:              client,
		userSubscriptions:   make(map[uuid.UUID]*redis.PubSub),
		companySubscriptions: make(map[uuid.UUID]*redis.PubSub),
	}
}

// SubscribeUser returns a channel for user events. If already subscribed, returns existing channel.
func (ps *PubSubSubscriber) SubscribeUser(userID uuid.UUID) (<-chan *redis.Message, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if sub, ok := ps.userSubscriptions[userID]; ok {
		return sub.Channel(), nil
	}

	channelName := fmt.Sprintf("user:%s", userID.String())
	sub := ps.client.Subscribe(context.Background(), channelName)
	ps.userSubscriptions[userID] = sub

	return sub.Channel(), nil
}

// SubscribeCompany returns a channel for company events. If already subscribed, returns existing channel.
func (ps *PubSubSubscriber) SubscribeCompany(companyID uuid.UUID) (<-chan *redis.Message, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if sub, ok := ps.companySubscriptions[companyID]; ok {
		return sub.Channel(), nil
	}

	channelName := fmt.Sprintf("company:%s", companyID.String())
	sub := ps.client.Subscribe(context.Background(), channelName)
	ps.companySubscriptions[companyID] = sub

	return sub.Channel(), nil
}

// UnsubscribeUser closes subscription for a given user if it exists.
func (ps *PubSubSubscriber) UnsubscribeUser(userID uuid.UUID) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if sub, ok := ps.userSubscriptions[userID]; ok {
		if err := sub.Close(); err != nil {
			return err
		}
		delete(ps.userSubscriptions, userID)
	}
	return nil
}

// UnsubscribeCompany closes subscription for a given company if it exists.
func (ps *PubSubSubscriber) UnsubscribeCompany(companyID uuid.UUID) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if sub, ok := ps.companySubscriptions[companyID]; ok {
		if err := sub.Close(); err != nil {
			return err
		}
		delete(ps.companySubscriptions, companyID)
	}
	return nil
}

