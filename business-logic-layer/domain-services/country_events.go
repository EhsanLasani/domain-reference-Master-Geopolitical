// ============================================================================
// FILE: country_events.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Events - Domain Events
// PURPOSE: Event definitions and outbox pattern implementation
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package events

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

// Event types for geopolitical domain
const (
	CountryCreatedEvent    = "geo.country.created.v1"
	CountryUpdatedEvent    = "geo.country.updated.v1"
	CountryDeactivatedEvent = "geo.country.deactivated.v1"
	RegionCreatedEvent     = "geo.region.created.v1"
	LanguageCreatedEvent   = "geo.language.created.v1"
)

// DomainEvent represents a domain event
type DomainEvent struct {
	ID          uuid.UUID   `json:"id"`
	Type        string      `json:"type"`
	AggregateID string      `json:"aggregate_id"`
	Version     int         `json:"version"`
	Data        interface{} `json:"data"`
	Metadata    Metadata    `json:"metadata"`
	OccurredAt  time.Time   `json:"occurred_at"`
}

type Metadata struct {
	UserID      string `json:"user_id,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	SourceIP    string `json:"source_ip,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

// Country-specific events
type CountryCreatedEventData struct {
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	ISO3Code     string `json:"iso3_code,omitempty"`
	OfficialName string `json:"official_name,omitempty"`
}

type CountryUpdatedEventData struct {
	CountryCode  string            `json:"country_code"`
	Changes      map[string]string `json:"changes"`
	PreviousVersion int            `json:"previous_version"`
	NewVersion   int               `json:"new_version"`
}

type CountryDeactivatedEventData struct {
	CountryCode string `json:"country_code"`
	Reason      string `json:"reason,omitempty"`
}

// EventPublisher interface for outbox pattern
type EventPublisher interface {
	Publish(event *DomainEvent) error
	PublishBatch(events []*DomainEvent) error
}

// OutboxEvent represents events stored in outbox table
type OutboxEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventType   string    `gorm:"not null"`
	AggregateID string    `gorm:"not null"`
	EventData   string    `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time `gorm:"default:now()"`
	ProcessedAt *time.Time
	Status      string    `gorm:"default:'pending'"` // pending, processed, failed
	RetryCount  int       `gorm:"default:0"`
	ErrorMsg    string
}

// EventStore handles event persistence
type EventStore struct {
	// Implementation would use database connection
}

func NewEventStore() *EventStore {
	return &EventStore{}
}

// Store saves event to outbox table
func (es *EventStore) Store(event *DomainEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	
	outboxEvent := &OutboxEvent{
		EventType:   event.Type,
		AggregateID: event.AggregateID,
		EventData:   string(eventData),
		CreatedAt:   time.Now(),
		Status:      "pending",
	}
	
	// Save to database (implementation needed)
	_ = outboxEvent
	return nil
}

// CreateCountryCreatedEvent creates a country created event
func CreateCountryCreatedEvent(countryCode, countryName string, metadata Metadata) *DomainEvent {
	return &DomainEvent{
		ID:          uuid.New(),
		Type:        CountryCreatedEvent,
		AggregateID: countryCode,
		Version:     1,
		Data: CountryCreatedEventData{
			CountryCode: countryCode,
			CountryName: countryName,
		},
		Metadata:   metadata,
		OccurredAt: time.Now(),
	}
}