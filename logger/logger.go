package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	EventClaimSubmitted EventType = "claim_submitted"
	EventClaimReversed  EventType = "claim_reversed"
)

// Event represents a logged event
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Logger handles event logging to file
type Logger struct {
	filePath string
	mutex    sync.Mutex
}

// NewLogger creates a new logger instance
func NewLogger(logDir string) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	filePath := filepath.Join(logDir, "pharmacy_events.json")

	return &Logger{
		filePath: filePath,
	}, nil
}

// LogClaimSubmission logs a claim submission event
func (l *Logger) LogClaimSubmission(claimID uuid.UUID, ndc, npi string, quantity int, price float64) error {
	event := Event{
		ID:        uuid.New().String(),
		Type:      EventClaimSubmitted,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"claim_id": claimID.String(),
			"ndc":      ndc,
			"npi":      npi,
			"quantity": quantity,
			"price":    price,
		},
	}

	return l.logEvent(event)
}

// LogClaimReversal logs a claim reversal event
func (l *Logger) LogClaimReversal(claimID uuid.UUID) error {
	event := Event{
		ID:        uuid.New().String(),
		Type:      EventClaimReversed,
		Timestamp: time.Now().UTC(),
		Data: map[string]interface{}{
			"claim_id": claimID.String(),
		},
	}

	return l.logEvent(event)
}

// logEvent writes an event to the log file
func (l *Logger) logEvent(event Event) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Read existing events
	events, err := l.readEvents()
	if err != nil {
		return fmt.Errorf("failed to read existing events: %w", err)
	}

	// Add new event
	events = append(events, event)

	// Write back to file
	return l.writeEvents(events)
}

// readEvents reads all events from the log file
func (l *Logger) readEvents() ([]Event, error) {
	// Check if file exists
	if _, err := os.Stat(l.filePath); os.IsNotExist(err) {
		// File doesn't exist, return empty slice
		return []Event{}, nil
	}

	// Read file
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// Parse JSON
	var events []Event
	if len(data) > 0 {
		if err := json.Unmarshal(data, &events); err != nil {
			return nil, fmt.Errorf("failed to parse log file: %w", err)
		}
	}

	return events, nil
}

// writeEvents writes events to the log file
func (l *Logger) writeEvents(events []Event) error {
	// Marshal to JSON with pretty formatting
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	// Write to file
	if err := os.WriteFile(l.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write log file: %w", err)
	}

	return nil
}

// GetEvents retrieves all logged events
func (l *Logger) GetEvents() ([]Event, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.readEvents()
}

// GetEventsByType retrieves events filtered by type
func (l *Logger) GetEventsByType(eventType EventType) ([]Event, error) {
	events, err := l.GetEvents()
	if err != nil {
		return nil, err
	}

	var filtered []Event
	for _, event := range events {
		if event.Type == eventType {
			filtered = append(filtered, event)
		}
	}

	return filtered, nil
}
