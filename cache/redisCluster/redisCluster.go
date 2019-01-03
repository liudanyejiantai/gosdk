// redis的集群调用处理
// 具体的配置信息可以参考conf目录下的文件
package redisCluster

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/liudanyejiantai/gosdk/cache/redisCluster/redisClusterConf"
)

type RedisCluster struct {
	// 集群的连接信息
	Slots func() ([]redis.ClusterSlot, error)

	// redis集群的db文件
	Redisdbs *redis.ClusterClient

	// 是否初始化
	bInit bool
}

func NewRedisCluster() (*RedisCluster, error) {
	c, err := redisClusterConf.NewConfig("redisCluster.conf")
	if err != nil {
		return nil, fmt.Errorf("read config file [redisCluster.conf] faild, %s", err.Error())
	}
	var redisCluster RedisCluster

	slots, err := c.ReadAllConf()
	if err != nil {
		return nil, fmt.Errorf("read config file [redisCluster.conf] ReadAllConf faild, %s", err.Error())
	}
	//fmt.Println("slots", slots)
	redisCluster.Slots = func() ([]redis.ClusterSlot, error) {
		return slots, nil
	}

	redisCluster.Redisdbs = redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots:  redisCluster.Slots,
		RouteRandomly: true,
	})
	if _, err = redisCluster.Redisdbs.Ping().Result(); err != nil {
		return nil, fmt.Errorf("redis cluster Ping faild, %s", err.Error())
	}

	// ReloadState reloads cluster state. It calls ClusterSlots func
	// to get cluster slots information.
	if err = redisCluster.Redisdbs.ReloadState(); err != nil {
		return nil, fmt.Errorf("redis cluster ReloadState faild, %s", err.Error())
	}

	redisCluster.bInit = true
	return &redisCluster, nil
}

// 在redis报错的时候需要调用下
func (cluster *RedisCluster) ReLoadRedisDb() error {

	if !cluster.bInit {
		return fmt.Errorf("redis cluster not init")
	}
	if err := cluster.Redisdbs.ReloadState(); err != nil {
		return fmt.Errorf("redis cluster ReloadState faild, %s", err.Error())
	}
	return nil
}

// 在redis报错的时候需要调用下ReLoadRedisDb,key默认设置18个小时
// 如果没有配置主从方式，master节点挂了整个redis集群就不可用了
// 如果业务逻辑简单的，还是自己通过hash将key散列存储比较好
func (cluster *RedisCluster) SetKey(key string, value string) error {
	if !cluster.bInit {
		return fmt.Errorf("redis cluster not init")
	}
	var err error
	d, _ := time.ParseDuration(fmt.Sprintf("+%ds", 60*60*18))
	// 试两次
	for i := 0; i < 2; i++ {
		if _, err = cluster.Redisdbs.Set(key, value, d).Result(); err != nil {
			//cluster.ReLoadRedisDb()
		} else {
			return nil
		}
	}
	return fmt.Errorf("redis cluster set key [%s] faild, %s", key, err.Error())
}

// 在redis报错的时候需要调用下ReLoadRedisDb
func (cluster *RedisCluster) GetKey(key string) (string, error) {
	if !cluster.bInit {
		return "", fmt.Errorf("redis cluster not init")
	}
	var err error
	var temp string
	// 试两次
	for i := 0; i < 2; i++ {
		temp, err = cluster.Redisdbs.Get(key).Result()
		fmt.Println(key, "temp, err", temp, err)
		if err != nil {
			//cluster.ReLoadRedisDb()
		} else {
			return temp, nil
		}
	}
	return "", fmt.Errorf("redis cluster get key [%s] faild, %s", key, err.Error())
}
