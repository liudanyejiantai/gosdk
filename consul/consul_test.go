// Copyright 2018 yejiantai Authors
//
// package consul consul的测试程序
package consul

import (
	"fmt"
	"testing"
	"time"

	tm "github.com/liudanyejiantai/gosdk/time"
)

// 测试consul的分布式锁
// go test -v
func Test_ConsulDistributedLock(t *testing.T) {

	client, er := CreateConsulClient("127.0.0.1:8500")
	if er != nil {
		fmt.Println(tm.GetFormatTime(), "CreateConsulClient", er)
		return
	}

	consulLock, er := NewConsulLock("DistributedLock", client)
	if er != nil {
		fmt.Println("NewConsulLock", er)
		return
	}
	for i := 0; i < 100; i++ {
		fmt.Println("consulLock.Lock()", consulLock.Lock())
		for j := 0; j < 200; j++ {
			fmt.Println(fmt.Sprintf("%s %s-%d-%d", tm.GetFormatTime(), consulLock.KeyName, i, j))
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Println(tm.GetFormatTime(), "consulLock.UnLock()", consulLock.UnLock())
	}
	fmt.Println(tm.GetFormatTime(), "consulLock.ReleaseLock()", consulLock.ReleaseLock())
	return
}
