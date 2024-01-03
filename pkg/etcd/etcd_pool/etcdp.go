package etcd_pool

import (
	"errors"
	"github.com/woshikedayaa/news-poster/pkg/utils/structutil"
	etcdc "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"io"
	"sync"
	"time"
)

var defaultEmptyPool = &EtcdPool{
	RWMutex:     new(sync.RWMutex),
	maxWaitTime: 3 * time.Minute,
	maxConn:     25,
	minConn:     10,
	minIdleConn: 5,
	conn:        nil,
	logger:      zap.NewNop(),
	config:      etcdc.Config{Endpoints: []string{"127.0.0.1:2379"}},
	freeChannel: make(chan *EtcdClientWrapper, 25),
}

type Pool interface {
	io.Closer
	GetConn() *EtcdClientWrapper
}

type EtcdPool struct {
	*sync.RWMutex
	// maxWaitTime
	// 这个指的是在所有链接都被使用
	// 且 len(conn) >= maxConn 的情况下 可以等待空闲链接的最大时间
	// 超过这个时间的将 panic
	maxWaitTime time.Duration
	freeChannel chan *EtcdClientWrapper
	config      etcdc.Config // config etcd
	maxConn     int          // maxConn 最大链接数
	minConn     int          // minConn 最小链接数
	minIdleConn int          // minIdleConn 最小空闲链接
	logger      *zap.Logger  // logger
	conn        []*EtcdClientWrapper
}

func (ep *EtcdPool) GetConn() *EtcdClientWrapper {
	// catch
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				panic("error when get a etcd client err=" + err.Error())
			}
		}
	}()

	// 先查找有没有空闲的
	if len(ep.freeChannel) != 0 {
		tmp, ok := <-ep.freeChannel
		if ok && tmp != nil && !tmp.isUsed {
			tmp.pushed = false
			tmp.isUsed = true
			return tmp
		}
	}

	// 没有找到空闲的 开始检查条件 是否达到最大上限
	ep.Lock()
	cur := len(ep.conn)
	if cur < ep.maxConn {
		// 没有达到 maxConn 创建新的链接
		c, err := NewEtcdClientWrapper(ep.config, cur)
		if err != nil {
			panic(err)
		}
		ep.conn = append(ep.conn, c)
		c.isUsed = true
		ep.Unlock()
		return c
	}
	ep.Unlock()

	// 最后达到上限了 而且 也没空闲的 就等待有空闲的
	timer := time.NewTimer(ep.maxWaitTime)
	for {
		select {
		case c, ok := <-ep.freeChannel:
			if !ok {
				panic(errors.New("EtcdPool.freeChannel has closed"))
			}
			return c
		case <-timer.C:
			panic(errors.New("EtcdPool.GetConn fail,timeout"))
		}
	}
}

// Clone 创建一个新的对象 其中 conn 不参与 clone
func (ep *EtcdPool) Clone() *EtcdPool {
	// 这里赋值私有属性
	// 如果后面有新的私有属性 需要来这里加入
	dst := &EtcdPool{
		RWMutex:     ep.RWMutex,
		maxConn:     ep.maxConn,
		minConn:     ep.minConn,
		minIdleConn: ep.minIdleConn,
		maxWaitTime: ep.maxWaitTime,
		config:      ep.config,
		logger:      ep.logger,
		freeChannel: make(chan *EtcdClientWrapper, ep.maxConn),
		// conn 不参与 clone
		conn: nil,
	}

	// 公有属性全交给了 util
	structutil.Clone(ep, dst)
	return dst
}

func (ep *EtcdPool) Close() error {
	var err error = nil
	ep.Lock()
	defer ep.Unlock()

	for i := 0; i < len(ep.conn); i++ {
		c := ep.conn[i]
		err = c.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ep *EtcdPool) pushChannel(c *EtcdClientWrapper) {
	if !c.pushed && !c.isUsed {
		c.pushed = true
		ep.freeChannel <- c
	}
}

func Create(options ...Option) (*EtcdPool, error) {
	var err error = nil
	ep := defaultEmptyPool.Clone()
	// 解析配置参数
	for _, option := range options {
		option.apply(ep)
	}

	if ep.maxConn*ep.minConn*ep.minIdleConn == 0 {
		return nil, errors.New("maxConn and minConn and minIdleConn can not be 0")
	}
	if ep.minConn < ep.minIdleConn {
		return nil, errors.New("minConn must bigger than minIdleConn")
	}
	if ep.maxConn < ep.minConn {
		return nil, errors.New("maxConn must bigger than minConn ")
	}

	// 创建连接
	for i := 0; i < max(ep.minConn, ep.minIdleConn); i++ {
		c, err := NewEtcdClientWrapper(ep.config, i)
		if err != nil {
			return nil, err
		}
		ep.conn = append(ep.conn, c)

		c.pushed = true
		ep.freeChannel <- c
	}
	// 开启守护协程
	go poolDaemon(ep)
	// ...
	return ep, err
}
