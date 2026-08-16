package analytics

import (
	"todo-cli/task"
)

type Analytics struct {
	TotalTasks     int            `json:"total_tasks"`
	CompletedTasks int            `json:"completed_tasks"`
	PendingTasks   int            `json:"pending_tasks"`
	InProgress     int            `json:"in_progress"`
	AvgPriority    float64        `json:"avg_priority"`
	TasksByTag     map[string]int `json:"tasks_by_tag"`
	CompletionRate float64        `json:"completion_rate"`
}

func Calculate(tasks []task.Task) Analytics {
	if len(tasks) == 0 {
		return Analytics{
			TasksByTag: make(map[string]int),
		}
	}

	a := Analytics{
		TasksByTag: make(map[string]int),
	}

	var totalPriority int

	for _, t := range tasks {
		a.TotalTasks++
		totalPriority += t.Priority

		switch t.Status {
		case task.StatusDone:
			a.CompletedTasks++
		case task.StatusInProgress:
			a.InProgress++
		case task.StatusPending:
			a.PendingTasks++
		}

		for _, tag := range t.Tags {
			a.TasksByTag[tag]++
		}
	}

	if a.TotalTasks > 0 {
		a.AvgPriority = float64(totalPriority) / float64(a.TotalTasks)
	}

	if a.TotalTasks > 0 {
		a.CompletionRate = float64(a.CompletedTasks) / float64(a.TotalTasks) * 100
	}

	return a
}
