package dto

type SubTask struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	Hash       string `json:"hash"`
	Length     int64  `json:"length"`
	Alphabet   string `json:"alphabet"`
	PartNumber int64  `json:"part_number"`
	PartCount  int64  `json:"part_count"`
}

type WorkerResponseDto struct {
	ID       string   `json:"id"`
	TaskID   string   `json:"task_id"`
	Answers  []string `json:"answers"`
	Progress float64  `json:"progress"`
	Status   string   `json:"status"`
}

type TaskState struct {
	ID       string   `json:"id"`
	TaskID   string   `json:"task_id"`
	Status   string   `json:"status"`
	Progress float64  `json:"progress"`
	Answers  []string `json:"answers"`
}
