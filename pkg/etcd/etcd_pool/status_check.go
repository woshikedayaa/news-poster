package etcd_pool

import (
	"go.uber.org/zap"
	"time"
)

// statusChecker 用来检查连接池内的链接
// 检查一个链接是否已经 Unwrapped -> 创建一个新的抛弃原来的
// 检查一个链接是否还可以正常使用 -> 不可用就抛弃 -> checkStatus
// 检查一个链接是否超过了使用的超时 -> 抛弃 -> connTimeout
func statusChecker(pool *EtcdPool) {
	var err error
	conn := pool.conn
	ticker := time.NewTicker(time.Second)
	idx := make([]int, pool.maxConn)
	I := 0

	for {
		pool.RLock()
		for _, c := range conn {
			if c.unWrapped || checkStatus(c) || connTimeout(c, pool.maxConnUseTime) {
				idx[I] = c.wrapperID
				I++
			}
		}
		pool.RUnlock()

		// 设置一个flag表示结束
		if I != pool.maxConn {
			idx[I] = -1
		}
		// 重置相关变量
		pool.Lock()
		err = handleCreateConn(conn, idx)
		pool.Unlock()

		if err != nil {
			// 这里不sync了 交给客户端去
			pool.logger.Error("error when handleCreateConn", zap.Error(err))
		}

		for i := 0; i < I; i++ {
			idx[i] = 0
		}
		I = 0
		// 限速
		select {
		case <-ticker.C:
		}
	}
}

func checkStatus(c *EtcdClientWrapper) bool {

}

func connTimeout(c *EtcdClientWrapper, maxTime time.Duration) bool {
	if maxTime < 0 {
		return false
	}
	if c.isUsed && time.Now().Sub(c.lastUsed).Milliseconds() > int64(maxTime) {
		return true
	}
	return false
}

func handleCreateConn(conn []*EtcdClientWrapper, idx []int) error {
	for i := 0; i < len(idx) && idx[i] != -1; i++ {

	}
}
