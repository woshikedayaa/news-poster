package etcd_pool

import (
	etcdc "go.etcd.io/etcd/client/v3"
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
			if c.unWrapped { // TODO 添加更多检查条件 完善容错
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
		// 处理有问题的链接
		err = handleCreateConn(conn, idx, &pool.config)
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

// handleCreateConn 这里处理需要处理的连接 直接覆盖（
func handleCreateConn(conn []*EtcdClientWrapper, idx []int, cfg *etcdc.Config) error {
	var (
		v   = 0
		c   *EtcdClientWrapper
		err error
	)

	for i := 0; i < len(idx) && idx[i] != -1; i++ {
		v = idx[i]
		c, err = NewEtcdClientWrapper(*cfg, v)
		if err != nil {
			return err
		}
		conn[v] = c
	}

	return nil
}
