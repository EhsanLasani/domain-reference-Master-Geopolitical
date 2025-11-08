// ============================================================================
// FILE: commands.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Business Logic Layer - CQRS Commands
// PURPOSE: Command definitions and handlers for write operations
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package cqrs

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// Command represents intent to change state
type Command interface {
	GetID() string
	GetCorrelationID() string
	GetRequestedBy() uuid.UUID
	GetRequestedAt() time.Time
}

// CreateCountryCommand for country creation
type CreateCountryCommand struct {
	ID            string    `json:"id"`
	CorrelationID string    `json:"correlation_id"`
	CountryCode   string    `json:"country_code" validate:"required,len=2"`
	CountryName   string    `json:"country_name" validate:"required,min=2,max=100"`
	OfficialName  string    `json:"official_name,omitempty" validate:"omitempty,max=200"`
	RequestedBy   uuid.UUID `json:"requested_by"`
	RequestedAt   time.Time `json:"requested_at"`
}

func (c *CreateCountryCommand) GetID() string           { return c.ID }
func (c *CreateCountryCommand) GetCorrelationID() string { return c.CorrelationID }
func (c *CreateCountryCommand) GetRequestedBy() uuid.UUID { return c.RequestedBy }
func (c *CreateCountryCommand) GetRequestedAt() time.Time { return c.RequestedAt }

// CommandResult represents command execution result
type CommandResult struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
	Events  []string  `json:"events"`
	Success bool      `json:"success"`
}

// CommandHandler processes commands
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) (*CommandResult, error)
}

// CreateCountryCommandHandler handles country creation
type CreateCountryCommandHandler struct {
	countryService interface{} // Would be actual service interface
	validator      interface{} // Command validator
}

func NewCreateCountryCommandHandler() *CreateCountryCommandHandler {
	return &CreateCountryCommandHandler{}
}

func (h *CreateCountryCommandHandler) Handle(ctx context.Context, cmd Command) (*CommandResult, error) {
	createCmd := cmd.(*CreateCountryCommand)
	
	// Command processing logic would go here
	// 1. Validate command
	// 2. Execute business logic
	// 3. Return result
	
	return &CommandResult{
		ID:      uuid.New(),
		Version: 1,
		Events:  []string{"CountryCreated"},
		Success: true,
	}, nil
}