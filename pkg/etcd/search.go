package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

const (
	scheme = "etcd"
)

type ResolverBuilder struct {
}

type Resolver struct {
	cc     resolver.ClientConn
	target resolver.Target

	cli     *clientv3.Client
	closeCh chan struct{}
	key     string
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	// ResolveNow 只是一个标志作用(
}

func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

func NewResolver(target resolver.Target, cc resolver.ClientConn) *Resolver {
	return &Resolver{
		target: target,
		cc:     cc,
		cli:    getConn(),
		key:    target.Endpoint(),
	}
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error = nil
	res := NewResolver(target, cc)
	// 这里配置完
	err = res.configCC()
	if err != nil {
		return nil, err
	}

	// watch 实时更新
	go res.watch()

	return res, nil
}

func (r *ResolverBuilder) Scheme() string {
	return scheme // etcd
}

func (r *Resolver) configCC() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return r.sync(ctx)
}

func (r *Resolver) sync(ctx context.Context) error {
	// 获取服务地址
	response, err := r.cli.Get(ctx, r.key)
	if err != nil {
		return err
	}
	var addresses []resolver.Address

	// 添加地址
	for _, v := range response.Kvs {
		addresses = append(addresses, resolver.Address{
			Addr: string(v.Value),
		})
	}

	// 更新cc
	return r.cc.UpdateState(
		resolver.State{
			Addresses: addresses,
			// TODO 把这个配置改到配置文件去
			ServiceConfig: r.cc.ParseServiceConfig("{\"loadBalancingPolicy\":\"round_robin\"}"),
		},
	)
}

func (r *Resolver) watch() {
	watchChan := r.cli.Watch(context.Background(), r.key, clientv3.WithPrefix())
	for {
		select {
		case response, ok := <-watchChan:
			if ok {
				r.handleEvent(response.Events)
			}
		}
	}
}

func (r *Resolver) handleEvent(events []*clientv3.Event) {
	for _, v := range events {
		switch v.Type {
		case clientv3.EventTypePut:

		case clientv3.EventTypeDelete:

		}
	}
}
