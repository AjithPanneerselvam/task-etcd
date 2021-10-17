package store

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Handle    string    `json:"handle"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsCompleted bool   `json:"isCompleted"`
}

type TaskStore interface {
	CreateTask(ctx context.Context, userID string, task Task) error
	ReadTask(ctx context.Context, userID string, taskID string) (*Task, error)
	ReadAllTasks(ctx context.Context, userID string) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, userID string, taskID string,
		taskStatus bool) (string, error)
	DeleteTask(ctx context.Context, userID string, taskID string) error
}
