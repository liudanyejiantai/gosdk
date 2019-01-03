// Copyright 2018 yejiantai Authors
//
// 默认是没有读取桶的全部数据情况，通过创建一个固定名称的sys_bdb_bk桶来确定桶是否存在
// sys_bdb_bk桶中前缀为bname_的名称都是表示具体的桶的前缀key
// 为了简化使用方便，将bucket部分隐藏掉
// boltdb的单个写入太慢了，批量写入速度是可以的
// 是否开启散列桶，开启就是分布更均匀，不开启就是填充慢6位，少了散列的计算，写入速度会快很多
// package bdb 基于boltdb的kv操作处理，对外屏蔽掉bucket的细节
package bdbkv

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/liudanyejiantai/gosdk/datatype"

	"github.com/liudanyejiantai/gosdk/public_func"

	"github.com/boltdb/bolt"
)

const (
	// 系统的桶名
	sys_boltdb_bucket string = "sys_bdb_bk"
	// 其他桶的前缀名称
	sys_bucket_name_prefix string = "bname_"
	// 系统参数key名称
	sys_param_bucket string = "sys_conf_bk"
	// 系统计数key
	sys_param_id_num int64 = 0
)

// boltdb封装
type Bdb struct {
	// 文件名称
	dbName string

	// 全部的桶数
	buckets map[string]string

	// 实例
	botdb *bolt.DB

	// 读写锁
	mtx_ sync.RWMutex

	// 是否散列桶
	isHashBucket bool

	// 是否已经打开
	isOpen bool
}

func (bdb *Bdb) GetAllBuckets() (map[string]string, error) {
	if !bdb.isOpen {
		return nil, dbNotOpeneErr()
	}
	return bdb.buckets, nil
}

// ID加1
func (bdb *Bdb) IdAdd() error {
	if !bdb.isOpen {
		return dbNotOpeneErr()
	}
	err := bdb.botdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sys_boltdb_bucket))
		by := b.Get([]byte(sys_param_bucket))
		id := datatype.BytesToInt64(by)
		return b.Put([]byte(sys_param_bucket), datatype.Int64ToBytes(id+1))

	})

	return err
}

// 获得ID
func (bdb *Bdb) GetId() (int64, error) {
	if !bdb.isOpen {
		return 0, dbNotOpeneErr()
	}
	var by []byte

	err := bdb.botdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sys_boltdb_bucket))
		by = b.Get([]byte(sys_param_bucket))
		return nil
	})

	return datatype.BytesToInt64(by), err
}

// 设置一个全局ID
func (bdb *Bdb) SetId(id int64) error {
	if !bdb.isOpen {
		return dbNotOpeneErr()
	}

	bdb.mtx_.Lock()
	defer bdb.mtx_.Unlock()
	err := bdb.botdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sys_boltdb_bucket))
		return b.Put([]byte(sys_param_bucket), []byte(datatype.Int64ToString(id)))
	})

	return err
}

// 批量插入一批数据(单条单条的插入不走事物是非常慢的)
func (bdb *Bdb) BatchPut(keys, values []string) error {
	if !bdb.isOpen {
		return dbNotOpeneErr()
	}
	if len(keys) != len(values) || len(keys) == 0 {
		return fmt.Errorf("keys values len error")
	}

	buckets := bdb.getBucketNames(keys)

	err := bdb.botdb.Batch(func(tx *bolt.Tx) error {
		var (
			b   *bolt.Bucket
			err error
		)

		for i := 0; i < len(keys); i++ {
			if b, err = tx.CreateBucketIfNotExists([]byte(buckets[i])); err != nil {
				return fmt.Errorf("create bucket [%s] faild, %s", buckets[i], err.Error())
			}
			if err = b.Put([]byte(keys[i]), []byte(values[i])); err != nil {
				return fmt.Errorf("bucket [%s] set key[%s] faild, %s", buckets[i], keys[i], err.Error())
			}
		}

		return nil
	})

	return err
}

// 获得桶里的数据,一定要确保bucketName存在，不然会抛异常
func (bdb *Bdb) GetKey(key string) (string, error) {
	if !bdb.isOpen {
		return "", dbNotOpeneErr()
	}
	bucketName := bdb.getBucketName(key)
	if !bdb.bucketExists(bucketName) {
		return "", fmt.Errorf("bucketName [%s] don't exists", bucketName)
	}

	var byt []byte
	err := bdb.botdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		byt = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("bucket [%s] key [%s] view value error,%s ", bucketName, key, err.Error())
	}
	if byt == nil {
		return "", fmt.Errorf("bucket [%s] dot have key [%s] 's value [%s] ", bucketName, key, string(byt))
	}

	return string(byt), err
}

// 往桶里插入数据
func (bdb *Bdb) PutKey(key, value string) error {
	if !bdb.isOpen {
		return dbNotOpeneErr()
	}
	bucketName := bdb.getBucketName(key)
	err := bdb.registerSysBucket(bucketName)
	if err != nil {
		return err
	}

	err = bdb.botdb.Batch(func(tx *bolt.Tx) error {
		var (
			b   *bolt.Bucket
			err error
		)
		fmt.Println("bucketName=", bucketName)
		if b, err = tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("bucket [%s] create  faild, %s", bucketName, err.Error())
		}
		if err = b.Put([]byte(key), []byte(value)); err != nil {
			return fmt.Errorf(" set key[%s] faild, %s", key, err.Error())
		}
		return nil
	})

	return err
}

func (bdb *Bdb) registerSysBucket(bucketName string) error {
	err := bdb.botdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(sys_boltdb_bucket))
		if err != nil {
			return fmt.Errorf("Bucket [%s] create faild %s ", bucketName, err.Error())
		}
		return b.Put([]byte(sys_bucket_name_prefix+bucketName), []byte(bucketName))
	})

	return err
}

func (bdb *Bdb) SetHashBucket(isHash bool) {
	bdb.isHashBucket = isHash
}

func (bdb *Bdb) getBucketName(key string) string {
	if bdb.isHashBucket {
		return public_func.GetMd5String(key)[:6]
	}
	return (key + "000000")[:6]
}

func (bdb *Bdb) getBucketNames(keys []string) []string {
	var array []string
	for i := 0; i < len(keys); i++ {
		array = append(array, bdb.getBucketName(keys[i]))
	}
	return array
}

// 遍历桶里的全部数据
func (bdb *Bdb) loopBucket(bucketName string) error {
	if !bdb.isOpen {
		return dbNotOpeneErr()
	}

	err := bdb.botdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		have := false
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%d bucket[%s],key=[%s],value=[%s]\n", time.Now().UnixNano()/1e6, bucketName, k, v)
			have = true
		}
		if !have {
			return fmt.Errorf("the bucket[%s] is empty")
		}
		return nil
	})

	return err
}

func dbNotOpeneErr() error {
	return fmt.Errorf("the boltdb is not open")
}

func (bdb *Bdb) loadSysBucket() error {
	err := bdb.botdb.Update(func(tx *bolt.Tx) error {
		var (
			b   *bolt.Bucket
			err error
		)
		if b, err = tx.CreateBucketIfNotExists([]byte(sys_boltdb_bucket)); err != nil {
			return fmt.Errorf("create system bucket [%s] faild, %s", sys_boltdb_bucket, err.Error())
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.HasPrefix(string(k), string(sys_bucket_name_prefix)) {
				bdb.buckets[string(v)] = string(k)
			}
		}
		return nil
	})

	return err
}

func (bdb *Bdb) bucketExists(bucketName string) bool {
	if _, ok := bdb.buckets[bucketName]; ok {
		return true
	}

	return false
}

// 创建一个boltdb实例
func NewDb(dbName string) (*Bdb, error) {
	var (
		db  = Bdb{dbName: dbName}
		err error
		opt bolt.Options = bolt.Options{Timeout: time.Second}
	)

	if db.botdb, err = bolt.Open(dbName, 0666, &opt); err != nil {
		return nil, fmt.Errorf("open boltdb file[%s] faild, %s", dbName, err.Error())
	}
	db.isOpen = true
	db.isHashBucket = false
	db.buckets = make(map[string]string)

	if err = db.loadSysBucket(); err != nil {
		return nil, err
	}

	return &db, nil
}
