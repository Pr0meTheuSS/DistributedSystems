package model

import "time"

type SubTask struct {
	ID         string
	TaskID     string
	Hash       string
	Length     int64
	Alphabet   string
	PartNumber int64
	PartCount  int64

	Status CrackHashStatus

	Answers  []string
	Progress float64

	CreatedAt time.Time
}
