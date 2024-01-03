package etcd_pool

import (
	"go.uber.org/zap"
	"time"
)

// poolDaemon 用来检查连接池内的链接
// 同时来维护freeChannel 保证其中有空闲链接
func poolDaemon(pool *EtcdPool) {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		pool.checkProblemConn()
		pool.pushFreeChannel()
		// 限速
		select {
		case <-ticker.C:
		}
	}
}

func (ep *EtcdPool) pushFreeChannel() {
	conn := ep.conn
	ep.Lock()
	defer ep.Unlock()

	for _, cli := range conn {
		ep.pushChannel(cli)
	}
}

// checkProblemConn
// 1-检查一个链接是否已经 Unwrapped -> 创建一个新的抛弃原来的
func (ep *EtcdPool) checkProblemConn() {
	conn := ep.conn
	I := 0
	idx := make([]int, ep.maxConn)
	ep.RLock()
	for _, c := range conn {
		if c.unWrapped { // TODO 添加更多检查条件 完善容错
			idx[I] = c.wrapperID
			I++
		}

	}
	ep.RUnlock()

	// 设置一个flag表示结束
	if I != ep.maxConn {
		idx[I] = -1
	}
	// 重置相关变量
	ep.Lock()
	// 处理有问题的链接
	err := handleCreateConn(ep, idx)
	ep.Unlock()

	if err != nil {
		// 这里不sync了 交给客户端去
		ep.logger.Error("error when handleCreateConn", zap.Error(err))
	}
}

// handleCreateConn 这里处理需要处理的连接 直接覆盖（
func handleCreateConn(ep *EtcdPool, idx []int) error {

	var (
		v    = 0
		c    *EtcdClientWrapper
		err  error
		conn = ep.conn
	)

	for i := 0; i < len(idx) && idx[i] != -1; i++ {
		v = idx[i]
		c, err = NewEtcdClientWrapper(ep.config, v)
		if err != nil {
			return err
		}
		conn[v] = c
	}

	return nil
}
