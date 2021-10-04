package db

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

// NewEtcdClient returns a new etcd client instance
func NewEtcdClient(etcdURLs []string) (clientv3.KV, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdURLs,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return clientv3.NewKV(client), nil
}
