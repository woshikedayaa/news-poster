package etcd

import (
	"bufio"
	"errors"
	"go.etcd.io/etcd/client/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func mustRefreshConn() *clientv3.Client {
	if err := RefreshConn(); err != nil {
		panic(err)
	}
	return nil
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

// parseEndpoints 解析文件夹下的 endpoints.list
// 这个文件通过每一行来写一个 etcd 客户端地址来解析
// e.g:
// *********************
// 127.0.0.1:23791
// 127.0.0.1:23792
// 127.0.0.1:23793
// *********************
func parseEndpoints() []string {
	dir, _ := os.Getwd()
	list, err := os.Open(dir + string(filepath.Separator) + "endpoints.list")
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	reader := bufio.NewReader(list)

	for {
		// 这里说个坑 这里必须用单引号 '\n'
		ed, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		ed = strings.TrimSpace(ed)
		res = append(res, ed)
		if errors.Is(err, io.EOF) {
			break
		}
	}
	//安全退出
	return res
}
