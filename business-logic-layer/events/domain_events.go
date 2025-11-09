package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DomainEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	TenantID    string                 `json:"tenant_id"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int                    `json:"version"`
}

type EventStore interface {
	Append(ctx context.Context, events []DomainEvent) error
	Load(ctx context.Context, aggregateID string) ([]DomainEvent, error)
}

type OutboxEvent struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Payload   string    `json:"payload"`
	TenantID  string    `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	Published bool      `json:"published"`
}

type OutboxPublisher struct {
	events []OutboxEvent
}

func NewOutboxPublisher() *OutboxPublisher {
	return &OutboxPublisher{
		events: make([]OutboxEvent, 0),
	}
}

func (op *OutboxPublisher) Publish(ctx context.Context, event DomainEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := OutboxEvent{
		ID:        uuid.New().String(),
		EventType: event.Type,
		Payload:   string(payload),
		TenantID:  event.TenantID,
		CreatedAt: time.Now(),
		Published: false,
	}

	op.events = append(op.events, outboxEvent)
	return nil
}

func (op *OutboxPublisher) GetUnpublishedEvents() []OutboxEvent {
	var unpublished []OutboxEvent
	for _, event := range op.events {
		if !event.Published {
			unpublished = append(unpublished, event)
		}
	}
	return unpublished
}

func (op *OutboxPublisher) MarkAsPublished(eventID string) {
	for i, event := range op.events {
		if event.ID == eventID {
			op.events[i].Published = true
			break
		}
	}
}

// Country-specific events
func NewCountryCreatedEvent(countryID, tenantID, countryCode, countryName string) DomainEvent {
	return DomainEvent{
		ID:          uuid.New().String(),
		Type:        "country.created",
		AggregateID: countryID,
		TenantID:    tenantID,
		Data: map[string]interface{}{
			"country_code": countryCode,
			"country_name": countryName,
		},
		Timestamp: time.Now(),
		Version:   1,
	}
}

func NewCountryUpdatedEvent(countryID, tenantID, countryCode, countryName string) DomainEvent {
	return DomainEvent{
		ID:          uuid.New().String(),
		Type:        "country.updated",
		AggregateID: countryID,
		TenantID:    tenantID,
		Data: map[string]interface{}{
			"country_code": countryCode,
			"country_name": countryName,
		},
		Timestamp: time.Now(),
		Version:   1,
	}
}