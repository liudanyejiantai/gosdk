// Copyright 2018 yejiantai Authors
//
// package env 操作环境变量，因为go自身的不会实际更改环境变量，需要通过调用处理
package env

import (
	"fmt"
	"os"
	"runtime"

	"github.com/liudanyejiantai/gosdk/exec/winCmd"
	"github.com/liudanyejiantai/gosdk/file"
)

// 获得环境变量str的值，如果没有返回为false
func GetEnv(str string) (string, bool) {
	return os.LookupEnv(str)
}

// 设置环境变量,value做追加
func SetEnvAdd(strKey, strValue string) error {
	return nil
}

// 设置环境变量,value直接做替换
func SetEnv(strKey, strValue string) error {
	var (
		osName string
	)
	osName = runtime.GOOS
	if osName == "windows" {
		return setWindowsEnv(strKey, strValue)
	} else if osName == "linux" {
	}
	return nil
}

//通过bat脚本设置windows的环境变量
func setWindowsEnv(strKey, strValue string) error {
	var (
		temp    string
		command string
		err     error
	)
	command, temp = "tmp.bat", fmt.Sprintf("@echo off \n@echo add env %s\nsetx %s %s\n@echo %%%s%%", strKey, strKey, strValue, strKey)
	if err = file.SaveFile([]byte(temp), command); err != nil {
		return err
	}

	_, err = winCmd.ExecLockCmd(command)
	return err
}
