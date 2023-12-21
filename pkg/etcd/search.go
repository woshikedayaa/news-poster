package etcd

import (
	"context"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"slices"
	"time"
)

const (
	scheme = "etcd"
)

type ResolverBuilder struct {
}

type Resolver struct {
	cc        resolver.ClientConn
	target    resolver.Target
	addresses []resolver.Address

	cli *clientv3.Client

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
	return r.sync()
}

func (r *Resolver) sync() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
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
	if len(addresses) == 0 {
		return errors.New(fmt.Sprintf("no service=%s discovered", r.key))
	}

	r.addresses = addresses

	// 更新cc
	return r.cc.UpdateState(
		resolver.State{
			Addresses: r.addresses,
			//// 把这个配置改到配置文件去
			//ServiceConfig: r.cc.ParseServiceConfig("{\"loadBalancingPolicy\":\"round_robin\"}"),
			// 写于 2023-12-21 19:51
			// 这个方案不好使 需要在初始化客户端连接的时候配置负载均衡
			// 如果使用这个方案 每次配置都要解析一个json 降低效率
		},
	)
}

func (r *Resolver) watch() {
	watchChan := r.cli.Watch(context.Background(), r.key, clientv3.WithPrefix())
	ticker := time.NewTicker(time.Duration(serviceExpireTime*2) * time.Second)
	for {
		select {
		case response, ok := <-watchChan:
			if ok {
				err := r.handleEvent(response.Events)
				if err != nil {
					//TODO logger
				}
			}
		case <-ticker.C:
			err := r.sync()
			if err != nil {
				//TODO logger
			}
		}
	}
}

func (r *Resolver) handleEvent(events []*clientv3.Event) error {
	var err error
	for _, v := range events {
		addr := string(v.Kv.Value)
		// ...
		switch v.Type {
		// 有新服务加入
		case clientv3.EventTypePut:
			if serviceInAddresses(r.addresses, addr) {
				continue
			}
			r.addresses = append(r.addresses, resolver.Address{Addr: addr})

		// 有服务被删除
		case clientv3.EventTypeDelete:
			if !serviceInAddresses(r.addresses, addr) {
				continue
			}

			// 更新 addresses
			r.addresses = slices.DeleteFunc(r.addresses, func(a resolver.Address) bool {
				if a.Addr == addr {
					return true
				}
				return false
			})

			err = r.cc.UpdateState(
				resolver.State{
					Addresses: r.addresses,
				},
			)

			if err != nil {
				return err
			}
		}
	}
	return nil
}
