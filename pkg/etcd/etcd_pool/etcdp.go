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
	maxConn:     25,
	minConn:     10,
	minIdleConn: 5,
	conn:        nil,
	logger:      zap.NewNop(),
	config:      etcdc.Config{Endpoints: []string{"127.0.0.1:2379"}},
}

type Pool interface {
	io.Closer
	GetConn() *EtcdClientWrapper
}

type EtcdPool struct {
	*sync.RWMutex
	// maxConnUseTime  一个链接最长使用时间
	// 超过这个时间的链接将检查是否可用
	// 如果不可用了就创建个新链接 同时抛弃这个链接
	// 设置成-1 来表示没有超时检查
	maxConnUseTime time.Duration
	config         etcdc.Config // config etcd
	maxConn        int          // maxConn 最大链接数
	minConn        int          // minConn 最小链接数
	minIdleConn    int          // minIdleConn 最小空闲链接
	logger         *zap.Logger  // logger
	conn           []*EtcdClientWrapper
}

func (ep *EtcdPool) GetConn() (*EtcdClientWrapper, error) {
	var tmp *EtcdClientWrapper
	// 先查找有没有空闲的
	for i := 0; i < len(ep.conn); i++ {
		tmp = ep.conn[i]
		if tmp.isUsed == false {
			tmp.isUsed = true
			return tmp, nil
		}
	}
	cur := len(ep.conn)
	// 没有找到空闲的 开始检查条件 是否达到最大上限
	if cur < ep.maxConn {
		// 没有达到 maxConn 创建新的链接
		c, err := NewEtcdClientWrapper(ep.config, cur)
		if err != nil {
			return nil, err
		}
		ep.conn = append(ep.conn, c)
		return c, nil
	}

	// 最后达到上限了 而且 也没空闲的 返回错误
	return nil, errors.New("the connection pool for etcd is currently in use and has reached its maximum limit")
}

type EtcdClientWrapper struct {
	*etcdc.Client
	wrapperID int // index
	lastUsed  time.Time
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

// UnWrapper 解析出来原生的etcd的client 然后这个链接就不归连接池管理
func (e *EtcdClientWrapper) UnWrapper() *etcdc.Client {
	e.unWrapped = true
	return e.Client
}

func (e *EtcdClientWrapper) Release() {
	e.isUsed = false
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
	if ep.maxConn*ep.minConn*ep.minIdleConn == 0 {
		return nil, errors.New("maxConn and minConn and minIdleConn can not be 0")
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
