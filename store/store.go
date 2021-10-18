package store

import (
	"context"
)

type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsCompleted bool   `json:"isCompleted"`
}

type TaskStore interface {
	UpsertTask(ctx context.Context, userID string, task Task) error
	ReadTask(ctx context.Context, userID string, taskID string) (*Task, error)
	ReadAllTasks(ctx context.Context, userID string) ([]Task, error)
	DeleteTask(ctx context.Context, userID string, taskID string) error
}
