package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AjithPanneerselvam/task-etcd/store"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

const (
	keyTaskFormat  = "task:%v:%v"
	keyTasksFormat = "task:%v"
)

// ErrTaskStore implements Error interface
type ErrTaskStore string

const (
	ErrTaskStoreNoRecord ErrTaskStore = "error no task record"
)

func (e ErrTaskStore) Error() string {
	return string(e)
}

type taskStore struct {
	clientv3.KV
}

func New(db clientv3.KV) store.TaskStore {
	return &taskStore{
		db,
	}
}

func (t *taskStore) UpsertTask(ctx context.Context, userID string, task store.Task) error {
	taskInBytes, err := json.Marshal(task)
	if err != nil {
		return errors.Wrap(err, "error marshalling task")
	}

	key := fmt.Sprintf(keyTaskFormat, userID, task.ID)

	_, err = t.Put(ctx, key, string(taskInBytes))
	if err != nil {
		return errors.Wrapf(err, "error creating task in the store")
	}

	return nil
}

func (t *taskStore) ReadTask(ctx context.Context, userID string, taskID string) (*store.Task, error) {
	key := fmt.Sprintf(keyTaskFormat, userID, taskID)

	resp, err := t.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) != 1 {
		return nil, ErrTaskStoreNoRecord
	}

	taskInBytes := resp.Kvs[0].Value

	var task store.Task
	err = json.Unmarshal(taskInBytes, &task)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling task response from store")
	}

	return &task, nil
}

func (t *taskStore) ReadAllTasks(ctx context.Context, userID string) ([]store.Task, error) {

	key := fmt.Sprintf(keyTasksFormat, userID)

	resp, err := t.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	tasks := make([]store.Task, 0)

	for _, val := range resp.Kvs {
		var task store.Task
		err := json.Unmarshal(val.Value, &task)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *taskStore) DeleteTask(ctx context.Context, userID string, taskID string) error {
	key := fmt.Sprintf(keyTaskFormat, userID, taskID)
	_, err := t.Delete(ctx, key)
	return err
}
