package model

import (
	"time"
)

type OutboxEvent struct {
	ID            string    `db:"id"`
	AggregateID   string    `db:"aggregate_id"`
	EventType     string    `db:"event_type"`
	Payload       []byte    `db:"payload"`
	Status        string    `db:"status"`
	Attempt       int       `db:"attempt"`
	NextAttemptAt time.Time `db:"next_attempt_at"`
	LastError     *string   `db:"last_error"`
}
