package service

type Task struct {
	ID              string  `json:"id"`
	Description     string  `json:"description"`
	Error           *string `json:"error,omitempty"`
	Percentage      *int    `json:"percentage,omitempty"`
	IsComplete      bool    `json:"isComplete"`
	IsIndeterminate bool    `json:"isIndeterminate"`
}

type TaskService struct {
}

func NewProcessService() *TaskService {
	return &TaskService{}
}
