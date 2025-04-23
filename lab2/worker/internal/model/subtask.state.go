package model

type TaskState struct {
	ID       string   `json:"id"`
	TaskID   string   `json:"task_id"`
	Status   string   `json:"status"`
	Progress float64  `json:"progress"`
	Answers  []string `json:"answers"`
}
