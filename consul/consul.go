// Copyright 2018 yejiantai Authors
//
// package consul consul连接操作
package consul

import (
	"errors"
	"fmt"
	"time"

	tm "github.com/liudanyejiantai/gosdk/time"

	api "github.com/hashicorp/consul/api"
)

// 创建一个consul本机客户端,通过consul自带的客户端实现集群连接
// hostaddr默认为127.0.0.1:8500
// consul agent -dev启动开发者模式
func CreateConsulClient(hostaddr string) (*api.Client, error) {
	var (
		client *api.Client
		err    error
	)

	if client, err = api.NewClient(&api.Config{Address: hostaddr}); err != nil {
		return nil, errors.New(fmt.Sprintf("create consul client faild, connect info[%s], %s", hostaddr, err.Error()))
	}

	// 添加一个状态判断，不然连接状态没法确定
	if _, err = client.Status().Leader(); err != nil {
		return nil, errors.New(fmt.Sprintf("client connect consul faild, connect info[%s], %s", hostaddr, err.Error()))
	}
	return client, nil
}

// 从consul集群中根据服务名称service_name获得一个可用的服务
// 可以使用命令curl测试返回
func GetService(client *api.Client, service_name string) (string, error) {
	var (
		services     map[string][]string
		err          error
		healthyOnly  bool
		name         string
		servicesData []*api.ServiceEntry
		entry        *api.ServiceEntry
		health       *api.HealthCheck
	)

	// 获得全部的服务
	if services, _, err = client.Catalog().Services(&api.QueryOptions{}); err != nil {
		return "", errors.New(fmt.Sprintf("list all services error, %s", err.Error()))
	}

	for name = range services {
		// 这个发现的是大的名称
		if name != service_name {
			continue
		}
		healthyOnly = true
		//第三个参数为true表示只发现健康的服务，false表示全部服务都发现
		if servicesData, _, err = client.Health().Service(name, "", healthyOnly, &api.QueryOptions{}); err != nil {
			return "", errors.New(fmt.Sprintf(" found service name[%s] health status faild, %s", name, err.Error()))
		}

		for _, entry = range servicesData {
			if service_name != entry.Service.Service {
				continue
			}
			for _, health = range entry.Checks {
				if health.ServiceName != service_name {
					continue
				}
				// str_info := fmt.Sprintf(`节点名称为[%s],服务名称[%s],服务ID[%s],健康状态[%s],服务地址信息[%s:%d]`,
				//	health.Node, health.ServiceName, health.ServiceID,
				//	health.Status, entry.Service.Address, entry.Service.Port)
				// fmt.Println(str_info)
			}
		}
	}
	return "", errors.New(fmt.Sprintf("found service name[%s] faild", service_name))
}

// 注册服务 reg表示一个服务 chk表示服务的检查属性
// reg.ID = "1234567890123"
// reg.Name = "http_svr"
// reg.Port = 8081
// reg.Address = "127.0.0.1"
// reg.Tags = []string{"http下载", "golang编写", "跨平台"}
// chk.HTTP = fmt.Sprintf("http://%s:%d", reg.Address, reg.Port)
// chk.Timeout = "5s"  //5秒的超时时间
// chk.Interval = "2s" //2秒的间隔时间
func RegistService(reg *api.AgentServiceRegistration, chk *api.AgentServiceCheck, client *api.Client) error {
	var (
		err error
	)
	reg.Check = chk
	if err = client.Agent().ServiceRegister(reg); err != nil {
		return errors.New(fmt.Sprintf("register service [%s] to consul faild, %s", reg.Name, err.Error()))
	}
	return nil
}

// 取消一个已经注册的服务
func DeRegistService(strRegid string, client *api.Client) error {
	var (
		err error
	)
	if err = client.Agent().ServiceDeregister(strRegid); err != nil {
		return errors.New(fmt.Sprintf("DeRegistServiceId [%s] from consul faild, %s", strRegid, err.Error()))
	}
	return nil
}

// consul的分布式锁
type ConsulLock struct {
	keyName string //分布式锁的名称
	stopCh  chan struct{}
	client  *api.Client //consul客户端
	locker  *api.Lock   //锁
	isOk    bool        //是否已经创建成功
}

// 创建一个新的consul分布式lock
func NewConsulLock(strKey string, client *api.Client) (*ConsulLock, error) {
	var (
		err error
		lk  ConsulLock
	)
	lk.keyName, lk.client = strKey, client
	if lk.locker, err = client.LockKey(lk.keyName); err != nil {
		return nil, errors.New(fmt.Sprintf("DistributedLock key [%s] faild, %s", lk.keyName, err.Error()))
	}
	lk.isOk = true
	return &lk, nil
}

// 是否创建成果
func (lk *ConsulLock) IsCreateOk() bool {
	return lk.isOk
}

// 上锁
func (lk *ConsulLock) Lock() error {
	var (
		err error
	)

	if _, err = lk.locker.Lock(lk.stopCh); err != nil {
		return errors.New(fmt.Sprintf("DistributedLock lock [%s] faild, %s", lk.keyName, err.Error()))
	}
	return nil
}

// 解锁
func (lk *ConsulLock) UnLock() error {
	var (
		err error
	)

	if err = lk.locker.Unlock(); err != nil {
		return err
	}

	return nil
}

// 释放锁,只有全部都不在使用的时候才释放成功
func (lk *ConsulLock) ReleaseLock() error {
	var (
		err error
	)

	if err = lk.locker.Destroy(); err != nil {
		return err
	}

	return nil
}

// 创建分布式锁
func DistributedLock(strKey string, client *api.Client) error {
	var (
		lock *api.Lock
		err  error
	)

	if lock, err = client.LockKey("DistributedLock/yjt"); err != nil {
		return errors.New(fmt.Sprintf("DistributedLock faild, %s", err.Error()))
	}
	for i := 0; i < 100; i++ {

		stopCh := make(chan struct{})
		_, err = lock.Lock(stopCh)
		if err != nil {
			//glog.Errorf("DistributedLock:%s", err.Error())
			return err
		}
		for j := 0; j < 500; j++ {
			fmt.Println(tm.GetFormatTime(), fmt.Sprintf("%s-%d", strKey, j))
			time.Sleep(time.Millisecond * 10)
		}
		err = lock.Unlock()
		if err == nil {
			fmt.Println("lock already unlocked")
		}

	}
	//必须全部没有再使用才能释放
	err = lock.Destroy()
	if err != nil {
		fmt.Println(tm.GetFormatTime(), "释放分布式锁失败,原因:"+err.Error())
	} else {
		fmt.Println(tm.GetFormatTime(), "释放分布式锁成功")
	}
	//<-lockCh //2018-11-09

	return nil
}

// 设置key值为strValue
func SetKey(strKey, strValue string, client *api.Client) error {
	var (
		err error
		kv  *api.KVPair
	)

	kv = &api.KVPair{
		Key:   strKey,
		Flags: 0,
		Value: []byte(strValue),
	}

	if _, err = client.KV().Put(kv, nil); err != nil {
		return errors.New(fmt.Sprintf("consul set key[%s] value[%s] faild, %s", strKey, strValue, err.Error()))
	}

	return nil
}

// 从consul获取key为strKey的键值,没有数据返回错误
func GetKey(strKey string, client *api.Client) (string, error) {
	var (
		err error
		kv  *api.KVPair
	)

	if kv, _, err = client.KV().Get(strKey, nil); err != nil {
		return "", errors.New(fmt.Sprintf("consul get key[%s] value faild, %s", strKey, err.Error()))
	}
	// 没有数据
	if kv == nil {
		return "", errors.New(fmt.Sprintf("key [%s] not exists", strKey))
	}

	return string(kv.Value), nil
}

// 从consul删除key为strKey的记录
func DeleteKey(strKey string, client *api.Client) error {
	var (
		err error
	)

	if _, err = client.KV().Delete(strKey, nil); err != nil {
		return errors.New(fmt.Sprintf("consul delete key[%s] value faild, %s", strKey, err.Error()))
	}
	return nil
}
