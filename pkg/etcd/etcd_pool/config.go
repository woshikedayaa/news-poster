package etcd_pool

import (
	etcdc "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

type Option func(pool *EtcdPool)

func (o Option) apply(pool *EtcdPool) {
	o(pool)
}

func WithLogger(logger *zap.Logger) Option {
	return func(pool *EtcdPool) {
		pool.logger = logger
	}
}

func WithEtcdClientConfig(config etcdc.Config) Option {
	return func(pool *EtcdPool) {
		pool.config = config
	}
}

func WithMaxConnUseTime(d time.Duration) Option {
	return func(pool *EtcdPool) {
		pool
	}
}
