package task

import (
	"time" //для работы со временем

	"github.com/google/uuid" //для генерации уникальных ID
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          string                 `json:"ID"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      Status                 `json:"status"`
	Priority    int                    `json:"priority"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func NewTask(title, description string, priority int, tags []string) *Task {
	now := time.Now() //текущее время

	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      StatusPending,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        tags,
		Metadata:    make(map[string]interface{}),
	}
}

func (t *Task) Complete() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

func (t *Task) UpdatePriority(newPriority int) {
	t.Priority = newPriority
	t.UpdatedAt = time.Now()
}
