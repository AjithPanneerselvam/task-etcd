package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AjithPanneerselvam/task-etcd/auth"
	"github.com/AjithPanneerselvam/task-etcd/store"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type TaskHandler struct {
	taskStore store.TaskStore
}

func NewTaskHandler(taskStore store.TaskStore) *TaskHandler {
	return &TaskHandler{
		taskStore: taskStore,
	}
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	userID, err := fetchUserIDFromCtx(ctx)
	if err != nil {
		log.Errorf("error fetching user id from ctx: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var task store.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Errorf("error unmarshalling task from request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.ID = uuid.NewString()
	task.IsCompleted = false

	err = t.taskStore.CreateTask(ctx, userID, task)
	if err != nil {
		log.Error("error storing task %v in the store: %v", task.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("task of id %v is successfully stored", task.ID)

	taskCreatedResponse := struct {
		TaskID string `json:"taskId"`
	}{
		task.ID,
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(taskCreatedResponse)
	if err != nil {
		log.Errorf("error encoding the task created response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := fetchUserIDFromCtx(ctx)
	if err != nil {
		log.Errorf("error fetching user id from ctx: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	taskID := chi.URLParam(r, "task-id")
	log.Debugf("task id: %v", taskID)

	task, err := t.taskStore.ReadTask(ctx, userID, taskID)
	if err != nil {
		log.Errorf("error reading task %v from store", taskID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("task of id %v: %v retrieved from store", task.ID, task)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Errorf("error encoding the task response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := fetchUserIDFromCtx(ctx)
	if err != nil {
		log.Errorf("error fetching user id from ctx: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tasks, err := t.taskStore.ReadAllTasks(ctx, userID)
	if err != nil {
		log.Errorf("error reading tasks from store: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("tasks %v retrieved from store", tasks)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Errorf("error encoding the tasks response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := fetchUserIDFromCtx(ctx)
	if err != nil {
		log.Errorf("error fetching user id from ctx: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	taskID := chi.URLParam(r, "task-id")

	w.WriteHeader(http.StatusOK)
	err = t.taskStore.DeleteTask(ctx, userID, taskID)
	if err != nil {
		log.Errorf("error encoding the tasks response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := fetchUserIDFromCtx(ctx)
	if err != nil {
		log.Errorf("error fetching user id from ctx: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	taskID := chi.URLParam(r, "task-id")

	var task store.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Errorf("error unmarshalling task from request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task.ID = taskID

	err = t.taskStore.CreateTask(ctx, userID, task)
	if err != nil {
		log.Error("error storing task %v in the store: %v", task.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func fetchUserIDFromCtx(ctx context.Context) (string, error) {
	userIDVal, err := auth.FetchClaimValFromCtx(ctx, auth.ClaimsKeyUserID)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching claim %v value", auth.ClaimsKeyUserID)
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return "", fmt.Errorf("error type asserting value %v of %v key", userIDVal, auth.ClaimsKeyUserID)
	}

	return userID, nil
}
