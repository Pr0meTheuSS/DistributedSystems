package model

import "time"

type CrackHashStatus string

const (
	PENDING             CrackHashStatus = "PENDING"
	IN_PROGRESS         CrackHashStatus = "IN_PROGRESS"
	ERROR               CrackHashStatus = "ERROR"
	READY               CrackHashStatus = "READY"
	PARTICALLY_FINISHED CrackHashStatus = "PARTICALLY_FINISHED"
)

type CrackHashRecord struct {
	ID         string
	Hash       string
	Alphabet   string
	Length     int64
	Status     CrackHashStatus
	CreatedAt  time.Time
	FinishedAt *time.Time
}

type CrackHashProgress struct {
	RecordID                string
	ProgressInPercents      float64
	CurrentHandlingDuration time.Duration
}

type CrackHashResult struct {
	RecordID string
	Answers  []string
}
