package etcd_pool

import (
	"context"
	"fmt"
	etcdc "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

var cfg = etcdc.Config{
	Endpoints:            []string{"127.0.0.1:23791", "127.0.0.1:23792", "127.0.0.1:23793"},
	AutoSyncInterval:     0,
	DialTimeout:          0,
	DialKeepAliveTime:    0,
	DialKeepAliveTimeout: 0,
	MaxCallSendMsgSize:   0,
	MaxCallRecvMsgSize:   0,
	TLS:                  nil,
	Username:             "",
	Password:             "",
	RejectOldCluster:     false,
	DialOptions:          nil,
	Context:              nil,
	Logger:               nil,
	LogConfig:            nil,
	PermitWithoutStream:  false,
}

func TestCreate(t *testing.T) {
	_, err := Create(WithEtcdClientConfig(cfg))
	if err != nil {
		t.Fatal(err)
	}
}

func TestEtcdPool_GetConn(t *testing.T) {
	pool, err := Create(WithEtcdClientConfig(cfg))
	if err != nil {
		t.Fatal(err)
	}
	conn := pool.GetConn()
	_, err = conn.KV.Put(context.Background(), "hello", "world")
	if err != nil {
		t.Fatal(err)
	}
	get, err := conn.KV.Get(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(get.Kvs[0].Value))
	conn.KV.Delete(context.Background(), "hello")
}

// 实际上pass了
func TestDaemon(t *testing.T) {
	maxConn := 20
	pool, err := Create(WithEtcdClientConfig(cfg), WithMaxConn(maxConn), WithMaxWaitTime(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	// 先获取全部链接 并且不释放
	for i := 0; i < maxConn; i++ {
		_ = pool.GetConn()
	}
	// 10s后就panic
	_ = pool.GetConn()
}

func TestEtcdPool_Close(t *testing.T) {
	pool, err := Create(WithEtcdClientConfig(cfg))
	if err != nil {
		t.Fatal(err)
	}
	pool.Close()
}
