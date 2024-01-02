package etcd_pool

import (
	"github.com/woshikedayaa/news-poster/pkg/utils/structutil"
	etcdc "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

var defaultEmptyPool = &EtcdPool{
	Mutex:       new(sync.Mutex),
	maxConn:     25,
	minConn:     10,
	minIdleConn: 5,
	conn:        nil,
	logger:      nil,
	config:      etcdc.Config{Endpoints: []string{"127.0.0.1:2379"}},
}

type EtcdClientWrapper struct {
	*etcdc.Client
	wrapperID int // index
	isUsed    bool
	unWrapped bool
}

func NewEtcdClientWrapper(cfg etcdc.Config, id int) (*EtcdClientWrapper, error) {
	var err error
	client, err := etcdc.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdClientWrapper{
		Client:    client,
		wrapperID: id,
		isUsed:    false,
		unWrapped: false,
	}, nil
}

func (e *EtcdClientWrapper) UnWrapper() *etcdc.Client {
	e.unWrapped = true
	return e.Client
}

func (e *EtcdClientWrapper) Release() {
	e.isUsed = false
}

type EtcdPool struct {
	*sync.Mutex
	// maxConnUseTime  一个链接最长使用时间
	// 超过这个时间的链接将检查是否可用
	// 如果不可用了就创建个新链接 同时抛弃这个链接
	maxConnUseTime time.Duration
	config         etcdc.Config // config etcd
	maxConn        int          // 最大链接数
	minConn        int          // 最小链接数
	minIdleConn    int          // 最小空闲链接
	logger         *zap.Logger  // logger
	conn           []*EtcdClientWrapper
}

// Clone 创建一个新的对象 其中 conn 不参与 clone
func (ep *EtcdPool) Clone() *EtcdPool {
	// 这里赋值私有属性
	// 如果后面有新的私有属性 需要来这里加入
	dst := &EtcdPool{
		Mutex:       ep.Mutex,
		maxConn:     ep.maxConn,
		minConn:     ep.minConn,
		minIdleConn: ep.minIdleConn,
		config:      ep.config,
		logger:      ep.logger,
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

func Create(options ...Option) (*EtcdPool, error) {
	var err error = nil
	ep := defaultEmptyPool.Clone()
	// 解析配置参数
	for _, option := range options {
		option.apply(ep)
	}
	// 创建连接
	for i := 0; i < max(ep.minConn, ep.minIdleConn); i++ {
		c, err := NewEtcdClientWrapper(ep.config, i)
		if err != nil {
			return nil, err
		}
		ep.conn = append(ep.conn, c)
	}
	// 开启检查
	go statusChecker(ep)
	// ...
	return ep, err
}
