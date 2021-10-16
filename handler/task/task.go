package task

import (
	"net/http"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
