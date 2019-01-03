// linux下的进程处理
package linuxPipe

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/liudanyejiantai/gosdk/exec/linuxCmd"
	"github.com/liudanyejiantai/gosdk/string_func"
)

type LinuxPipe struct {
	// 命令执行超时时间
	timeOutMs int64
}

// 确保超时时间不会小于1s
func (p *LinuxPipe) lazyTimeout() {
	if p.timeOutMs == 0 {
		p.timeOutMs = 1000
	}
}

// 设置超时时间
func (p *LinuxPipe) SetTimeOutMs(timeOutMs int64) {
	p.timeOutMs = timeOutMs
}

// 根据进程名称关闭程序
func (p *LinuxPipe) KillProcessName(processName string) (bool, error) {
	array, err := p.GetPid(processName)
	if err != nil {
		return false, err
	}
	if len(array) == 0 {
		return false, nil
	}
	return p.KillProcess(array[0])
}

// 获得ps名称查找出进程对应的pid，可能会有多个,如果是包含test的命令如下
// ps -e |grep test | awk '{print $1}'
func (p *LinuxPipe) GetPid(processName string) ([]int, error) {
	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("ps -e | grep %s | awk '{print $1}'", processName))

	p.lazyTimeout()
	var (
		array []int = []int{}
		line  string
		n     int
	)
	// ps -e 只能搜索出15位长度的内容
	if len(processName) > 15 {
		processName = processName[:14]
	}
	fullCmd := fmt.Sprintf("ps -e | grep %s | awk '{print $1}'", processName)
	// cmd中带参数不能运行很奇怪
	str, err := linuxCmd.ExecLockCmd(fullCmd)
	if err != nil {
		return array, fmt.Errorf("ExecLockCmd [%s] faild, %s", fullCmd, err.Error())
	}

	arrayStr := strings.Split(str, "\n")

	for i := 0; i < len(arrayStr); i++ {
		line = string_func.ConvertToSingleTrim(arrayStr[i])
		line = strings.Replace(line, " ", "", -1)

		if line == "" {
			continue
		}

		if n, err = strconv.Atoi(line); err != nil {
			return array, fmt.Errorf("GetPid faild, convert [%s] error,%s", line, err.Error())
		}
		array = append(array, n)
	}
	return array, nil
}

// 开启路径path的程序进程，显示还是在当前页面中
// 参考使用 go wPipe.StartProcess("/opt/sunyard/file_svr", []string{"-port=8689"})
func (p *LinuxPipe) StartProcess(path string, args []string) error {
	p.lazyTimeout()

	path = "/" + strings.TrimLeft(path, "/")
	temp := strings.TrimRight(path, " ")
	if !FileExists(temp) {
		return fmt.Errorf("file [%s] not exist, can't StartProcess", temp)
	}

	fullCmd := fmt.Sprintf("%s %s &", path, strings.Join(args, " "))
	return linuxCmd.ExecAysncCmd(fullCmd, p.timeOutMs)
}

// 判断文件是否存在  存在返回true 不存在返回false
func FileExists(filename string) bool {
	var (
		f   os.FileInfo
		err error
	)

	if f, err = os.Stat(filename); err != nil {
		return false
	}
	if os.IsNotExist(err) || f.IsDir() {
		return false
	}

	return true
}

// 带参数cmd开启path可执行程序
func (p *LinuxPipe) RunProcess(path string, params []string) (bool, error) {
	err := p.StartProcess(path, params)
	if err != nil {
		return false, err
	}
	return true, nil

}

// 进程processName是否有在运行
func (p *LinuxPipe) IsProcessRun(processName string) (bool, error) {
	array, err := p.GetPid(processName)
	if err != nil {
		return false, err
	}
	if len(array) == 0 {
		return false, nil
	}
	return true, nil
}

// 根据进程ID关闭程序
func (p *LinuxPipe) KillProcess(pid int) (bool, error) {
	command := fmt.Sprintf("kill -9 %d", pid)
	err := linuxCmd.ExecAysncCmd(command, p.timeOutMs)
	if err != nil {
		return false, err
	}
	return true, nil
}
