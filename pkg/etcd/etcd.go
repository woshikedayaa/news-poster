package etcd

import (
	"go.etcd.io/etcd/client/v3"
	"time"
)

// inits
var (
	globalConn = struct {
		conn *clientv3.Client
	}{
		conn: nil,
	}

	retryMax = 5 // 最大重试次数
)

func mustRefreshConn() {
	if err := RefreshConn(); err != nil {
		panic(err)
	}
}

func getConn() *clientv3.Client {
	// lazy
	// double if check
	if globalConn.conn == nil {
		if globalConn.conn == nil {
			mustRefreshConn()
		}
	}

	return globalConn.conn
}

// RefreshConn
// auto retry see: retryMax
// TODO 实现一个etcd连接池
func RefreshConn() error {
	var client *clientv3.Client
	var err error
	ep := parseEndpoints()

	// retry
	for i := 0; i < retryMax; i++ {
		client, err = clientv3.New(clientv3.Config{
			Endpoints:            ep,
			AutoSyncInterval:     0,
			DialTimeout:          15 * time.Second, // 15s
			DialKeepAliveTime:    0,
			DialKeepAliveTimeout: 0,
			MaxCallSendMsgSize:   0,
			MaxCallRecvMsgSize:   0,
			//TODO etcd-TLS
			TLS:                 nil,
			Username:            "",
			Password:            "",
			RejectOldCluster:    false,
			DialOptions:         nil,
			LogConfig:           nil,
			Context:             nil,
			PermitWithoutStream: false,
		})

		if client != nil {
			break
		}
	}

	if client == nil {
		return err
	}

	globalConn.conn = client
	return nil
}
