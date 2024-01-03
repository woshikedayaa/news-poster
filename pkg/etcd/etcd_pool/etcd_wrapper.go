package etcd_pool

import (
	etcdc "go.etcd.io/etcd/client/v3"
)

type EtcdClientWrapper struct {
	*etcdc.Client
	wrapperID int // index
	isUsed    bool
	unWrapped bool
	pushed    bool // pushed 代表了这个链接是否已经被push到了 EtcdPool.freeChannel 中
}

func NewEtcdClientWrapper(cfg etcdc.Config, id int) (*EtcdClientWrapper, error) {
	client, err := etcdc.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdClientWrapper{
		Client: client,

		wrapperID: id,
		isUsed:    false,
		unWrapped: false,
		pushed:    false,
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
