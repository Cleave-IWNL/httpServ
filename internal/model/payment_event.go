package model

import (
	"time"
)

type PaymentCreatedEvent struct {
	EventID       string    `json:"event_id"`
	EventType     string    `json:"event_type"`
	OccurredAt    time.Time `json:"ocсurred_at"`
	SchemaVersion int       `json:"schema_version"`
	PaymentID     string    `json:"payment_id"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
}

const PaymentCreatedType = "payment.created"
