package dto

type SubTask struct {
	TaskID     string `json:"task_id"`
	Hash       string `json:"hash"`
	Length     int64  `json:"length"`
	Alphabet   string `json:"alphabet"`
	PartNumber int64  `json:"part_number"`
	PartCount  int64  `json:"part_count"`
}
