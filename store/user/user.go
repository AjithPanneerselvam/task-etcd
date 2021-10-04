package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AjithPanneerselvam/todo/store"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

const (
	keyUserFormat = "user:%v"
)

// ErrUserStore implements Error interface
type ErrUserStore string

const (
	ErrUserStoreNoRecord ErrUserStore = "error no user record"
)

func (e ErrUserStore) Error() string {
	return string(e)
}

type userStore struct {
	clientv3.KV
}

func New(db clientv3.KV) store.UserStore {
	return &userStore{
		db,
	}
}

func (u *userStore) CreateUser(ctx context.Context, user store.User) error {
	userInBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "error marshalling user info")
	}

	key := fmt.Sprintf(keyUserFormat, user.ID)

	_, err = u.Put(ctx, key, string(userInBytes))
	if err != nil {
		return errors.Wrapf(err, "error creating user in the store")
	}

	return nil
}

func (u *userStore) GetInfoByID(ctx context.Context, userID string) (*store.User, error) {
	key := fmt.Sprintf(keyUserFormat, userID)

	resp, err := u.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching user info from store")
	}

	if len(resp.Kvs) != 1 {
		return nil, ErrUserStoreNoRecord
	}

	userInBytes := resp.Kvs[0].Value

	var user store.User
	err = json.Unmarshal(userInBytes, &user)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling user response from store")
	}

	return &user, nil
}
