package etcd

import (
	"github.com/woshikedayaa/news-poster/pkg/etcd/etcd_pool"
	"go.etcd.io/etcd/client/v3"
	"time"
)

// inits
var (
	globalConnPool *etcd_pool.EtcdPool

	retryMax = 5 // 最大重试次数
)

func mustRefreshConn() {
	if err := RefreshConn(); err != nil {
		panic(err)
	}
}

func getConn() *etcd_pool.EtcdClientWrapper {
	// lazy
	// double if check
	if globalConnPool == nil {
		if globalConnPool == nil {
			mustRefreshConn()
		}
	}

	return globalConnPool.GetConn()
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
	// TODO
	// globalConnPool.conn = client
	return nil
}
