// Copyright 2018 yejiantai Authors
//
// package windows 操作command
package winCmd

import (
	"bufio"
	"context"

	"fmt"
	"io"
	"os/exec"
	"time"
)

// 执行命令,先进cmd,然后在里面输入执行命令
// fullCmd表示完整的执行命令，该操作是非阻塞的，能够快速返回的或者不会阻塞的操作不要调用该接口
// start D:/file_svr.exe -port=8689
// timeOutMs表示超时时间,单位为毫秒，超时会跳出处理，不会将界面一直卡住
// 返回为读取到的信息
func ExecAysncCmd(path, fullCmd string, timeOutMs int64) error {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeOutMs))
	defer cancel()
	defer time.Sleep(time.Second * 3)

	var (
		cmd *exec.Cmd
		err error
	)
	cmd = exec.CommandContext(ctx, "cmd", "/C", fullCmd)
	if path != "" {
		cmd.Dir = path
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("ExecAysncCmd start [%s] faild, %s", fullCmd, err.Error())
	}

	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("ExecAysncCmd Wait [%s] faild, %s", fullCmd, err.Error())
	}
	return nil
}

// 执行阻塞方式的命令，能够获得管道返回的数据内容
// 调用的时候需要考虑阻塞出了进程
func ExecLockCmd(fullCmd string) (string, error) {
	var (
		out    string
		line   string
		err    error
		stdout io.ReadCloser
		reader *bufio.Reader
	)
	cmd := exec.Command("cmd", "/C", fullCmd)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return "", fmt.Errorf("ExecLockCmd StdoutPipe [%s] faild,%s", fullCmd, err.Error())
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("ExecLockCmd Start [%s] faild,%s", fullCmd, err.Error())
	}

	// 实时循环读取输出流中的一行内容
	reader = bufio.NewReader(stdout)
	for {
		line, err = reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		out += line
	}

	if err = cmd.Wait(); err != nil {
		return out, fmt.Errorf("ExecLockCmd Wait [%s] faild,%s", fullCmd, err.Error())
	}
	return out, nil
}
