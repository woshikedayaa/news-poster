package etcd

import (
	"bufio"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"io"
	"os"
	"time"
)

// inits
var globalConn = struct {
	conn *clientv3.Client
}{
	conn: mustOpenNewConn(),
}

func mustOpenNewConn() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            parseEndpoints(),
		AutoSyncInterval:     0,
		DialTimeout:          15 * time.Second, // 15s
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		LogConfig:            nil,
		Context:              nil,
		PermitWithoutStream:  false,
	})
	if err != nil {
		panic(err)
	}
	return client
}

func getConn() *clientv3.Client {
	if globalConn.conn == nil {
		globalConn.conn = mustOpenNewConn()
	}
	return globalConn.conn
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
	list, err := os.Open(dir + "/pkg/etcd/endpoints.list")
	if err != nil {
		panic(err)
	}
	res := make([]string, 0)
	reader := bufio.NewReader(list)
	for {
		// 这里说个坑 这里必须用单引号 '\n'
		ed, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		res = append(res, ed)
	}
	//安全退出
	return res
}
