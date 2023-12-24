package etcd

import (
	"bufio"
	"errors"
	"google.golang.org/grpc/resolver"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func genKey(raw string) stringWrapper {
	// 通过时间戳来补全不同服务
	return stringWrapper{
		raw:       raw,
		processed: raw + strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

func genValue(raw string) stringWrapper {
	//直接就是原来内容
	return stringWrapper{
		raw:       raw,
		processed: raw,
	}
}

// parseEndpoints 解析文件夹下的 endpoints.list
// 这个文件通过每一行来写一个 etcd 客户端地址来解析
// e.g.
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
		if len(ed) != 0 {
			res = append(res, ed)
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	//安全退出
	return res
}

func serviceInAddresses(addresses []resolver.Address, target string) bool {
	for _, address := range addresses {
		if address.Addr == target {
			return true
		}
	}
	return false
}
