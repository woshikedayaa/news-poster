package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Watch(ctx context.Context, key string, option ...clientv3.OpOption) clientv3.WatchChan {
	conn := getConn()
	watchChan := conn.Watch(ctx, key, option...)
	return watchChan
}

func GrantLease(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	conn := getConn()
	return conn.Lease.Grant(ctx, ttl)
}

func GrandLeaseWithKeepAlive(ctx context.Context, ttl int64) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	conn := getConn()
	response, err := conn.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	alive, err := conn.KeepAlive(ctx, response.ID)
	if err != nil {
		return nil, err
	}

	return alive, nil
}
