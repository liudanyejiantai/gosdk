// Copyright 2018 yejiantai Authors
//
// package ldb leveldb操作处理
package ldb

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/liudanyejiantai/gosdk/datatype"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// leveldb信息
type Ldb struct {
	// leveldb是否已经打开
	isOpen bool

	// 互斥锁
	mtx sync.RWMutex

	// 取自增长ID的value用
	mtxAuto sync.RWMutex

	// db
	db *leveldb.DB

	// 文件名称
	strFile string
}

const (
	// leveldb目录的后缀
	LeveldbSuffix = ".leveldb"

	// 初始值
	InitNum = int64(1000000)
)

//创建一个leveldb
func CreateLdb(strFile string) (*Ldb, error) {
	var (
		ldb Ldb
		err error
	)
	if !strings.HasSuffix(strFile, LeveldbSuffix) {
		strFile += LeveldbSuffix
	}
	ldb.strFile = strFile
	if ldb.db, err = leveldb.OpenFile(strFile, nil); err != nil {
		return nil, errors.New(fmt.Sprintf("leveldb open file [%s] faild", strFile))
	}
	ldb.isOpen = true
	return &ldb, nil
}

// 从leveldb中获取strKey对应的增长的序号value数据,支持多协程操作
func (ldb *Ldb) GetAutoIdentityKey(strKey []byte) (int64, error) {
	ldb.mtxAuto.Lock()
	defer ldb.mtxAuto.Unlock()
	if !ldb.isOpen {
		return -1, errors.New(fmt.Sprintf("leveldb file [%s] is not open", ldb.strFile))
	}
	var (
		byt []byte
		err error
		n   int64
	)
	// 如果leveldb中有就直接获得
	if byt, err = ldb.db.Get(strKey, nil); err == nil {
		if n, err = datatype.StringToInt64(string(byt)); err != nil {
			return -1, err
		}
		return n, nil
	}

	return ldb.getNewNo(string(strKey))
}

// 从leveldb中获得一个新的序号,如果没有返回一个初始值，有就在原来的基础上加1
// 为了防止数据格式太短，初始值设置为1000000
func (ldb *Ldb) getNewNo(strKey string) (int64, error) {
	var (
		LAST_MAX_NO = "LAST_MAX_NO"
		byt         []byte
		err         error
		str         string
		n           int64
	)
	// 如果原来一条记录都没有
	if byt, err = ldb.db.Get([]byte(LAST_MAX_NO), nil); err != nil {
		if err = ldb.db.Put([]byte(LAST_MAX_NO), []byte(datatype.Int64ToString(InitNum)), nil); err != nil {
			return -1, err
		}
		return InitNum, nil
	}

	// 如果有就在原来的基础上加1
	if n, err = datatype.StringToInt64(string(byt)); err != nil {
		return -1, err
	}
	str = datatype.Int64ToString(n + 1)
	if err = ldb.db.Put([]byte(LAST_MAX_NO), []byte(str), nil); err != nil {
		return -1, err
	}

	if err = ldb.db.Put([]byte(strKey), []byte(str), nil); err != nil {
		return -1, err
	}

	return n + 1, nil
}

// 往leveldb中获取设置数据
func (ldb *Ldb) PutLevelDbKey(strKey, strValue []byte) error {
	ldb.mtx.Lock()
	defer ldb.mtx.Unlock()
	if !ldb.isOpen {
		return errors.New(fmt.Sprintf("leveldb file [%s] is not open", ldb.strFile))
	}

	return ldb.db.Put(strKey, strValue, nil)
}

// 从leveldb中获取数据,err不为nil表示没有数据
func (ldb *Ldb) GetLevelDbValue(strKey []byte) ([]byte, error) {
	ldb.mtx.Lock()
	defer ldb.mtx.Unlock()
	if !ldb.isOpen {
		return []byte(""), errors.New(fmt.Sprintf("leveldb file [%s] is not open", ldb.strFile))
	}

	return ldb.db.Get(strKey, nil)
}

// 从leveldb中删除数据
func (ldb *Ldb) DeleteLevelDbKey(strKey []byte) error {
	ldb.mtx.Lock()
	defer ldb.mtx.Unlock()
	if !ldb.isOpen {
		return errors.New(fmt.Sprintf("leveldb file [%s] is not open", ldb.strFile))
	}

	return ldb.db.Delete(strKey, nil)
}

// 遍历leveldb
func (ldb *Ldb) LoopDb() {
	var (
		iter iterator.Iterator
	)
	iter = ldb.db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(fmt.Sprintf("key[%s],value[%s]", string(iter.Key()), string(iter.Value())))
	}
	iter.Release()
	fmt.Println(iter.Error())
}
