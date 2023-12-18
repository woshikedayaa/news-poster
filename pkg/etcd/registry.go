package etcd

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var (
	serviceExpireTime = 10 // 10s
)

type stringWrapper struct {
	raw       string
	processed string
}

type Register struct {
	mu *sync.Mutex

	key       stringWrapper
	value     stringWrapper
	leaseId   clientv3.LeaseID
	kc        <-chan *clientv3.LeaseKeepAliveResponse
	closeChan chan struct{}
}

func (r *Register) Close() error {
	if r.closeChan != nil {
		return errors.New("register has closed")
	}
	r.closeChan <- struct{}{}
	return nil
}

func (r *Register) close() error {
	conn := getConn()
	_, err := conn.Revoke(context.Background(), r.leaseId)
	if err != nil {
		return err
	}
	_, err = conn.Delete(context.Background(), r.key.processed)
	if err != nil {
		return err
	}
	// set empty
	r.kc = nil
	r.closeChan = nil
	return nil
}

func RegistryNewService(key, value string) (*Register, error) {
	var err error
	//
	res, err := registryNewService(key, value)
	if err != nil {
		return nil, err
	}

	go keepAlive(res, 0)
	return res, nil
}

// registryNewService put 相关操作
func registryNewService(key, value string) (*Register, error) {
	conn := getConn()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(serviceExpireTime)*time.Second)
	defer cancelFunc()

	leases, err := conn.Grant(ctx, int64(serviceExpireTime))
	if err != nil {
		return nil, err
	}

	kc, err := conn.KeepAlive(context.Background(), leases.ID)
	if err != nil {
		return nil, err
	}

	res := &Register{
		mu:        new(sync.Mutex),
		key:       genKey(key),
		value:     genValue(value),
		leaseId:   leases.ID,
		kc:        kc,
		closeChan: make(chan struct{}),
	}
	//TODO put
	_, err = conn.Put(
		context.Background(),
		res.key.processed,
		res.value.processed,
		clientv3.WithLease(res.leaseId),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func keepAlive(r *Register, retryCount int) {
	// check
	if retryCount > retryMax {
		panic(errors.New("service retry has arrive retryMax -> panic"))
	}

	timeoutDuration := time.Duration(serviceExpireTime*2) * time.Second
	timer := time.NewTimer(timeoutDuration)
	for {
		select {
		case <-r.closeChan:
			err := r.close()
			//TODO logger
			panic(err)
			return

		case <-r.kc:
			timer.Reset(timeoutDuration)
			continue

		//超时 尝试重连
		case <-timer.C:
			//TODO logger
			register, err := registryNewService(r.key.raw, r.value.raw)
			if err != nil {
				return
			}
			r = register
			go keepAlive(register, retryCount+1)
			return
		}
	}
}
