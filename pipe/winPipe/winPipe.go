// windows下的进程处理
package winPipe

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/liudanyejiantai/gosdk/exec/winCmd"
	"github.com/liudanyejiantai/gosdk/string_func"
)

type WindowsPipe struct {
	// 命令执行超时时间
	timeOutMs int64
}

// 确保超时时间不会小于1s
func (p *WindowsPipe) lazyTimeout() {
	if p.timeOutMs == 0 {
		p.timeOutMs = 1000
	}
}

// 设置超时时间
func (p *WindowsPipe) SetTimeOutMs(timeOutMs int64) {
	p.timeOutMs = timeOutMs
}

// 根据进程名称关闭程序
func (p *WindowsPipe) KillProcessName(processName string) (bool, error) {
	array, err := p.GetPid(processName)
	if err != nil {
		return false, err
	}
	if len(array) == 0 {
		return false, nil
	}
	return p.KillProcess(array[0])
}

// 使用精确查找的方式来处理
// wmic的排序其实是按照字母顺序排序的
// 精确查找 wmic process where name="file_svr.exe" get executablepath,name,processid
// 模糊查找 wmic process get processid,name |findstr /i "svr.exe"
// tasklist |findstr /i "svr.exe"
func (p *WindowsPipe) GetPid(processName string) ([]int, error) {
	p.lazyTimeout()
	var (
		array           []int = []int{}
		line, lowerName string
		n               int
	)
	fullCmd := "tasklist"
	// cmd中带参数不能运行很奇怪
	str, err := winCmd.ExecLockCmd(fullCmd)
	if err != nil {
		return array, fmt.Errorf("ExecLockCmd [%s] faild, %s", fullCmd, err.Error())
	}

	lowerName = strings.ToLower(processName)
	arrayStr := strings.Split(str, "\n")

	for i := 0; i < len(arrayStr); i++ {
		line = arrayStr[i]

		line = strings.ToLower(line)
		line = string_func.ConvertToSingleTrim(line)
		line = strings.TrimRight(line, " ")

		// 第一行ProcessId过滤掉
		if line == "" || !strings.Contains(line, lowerName) || i == 1 {
			continue
		}

		arraytemp := strings.Split(line, " ")
		if len(arraytemp) <= 2 {
			return array, fmt.Errorf("GetPid faild, convert [%s] to array faild，length to small", line)
		}

		if n, err = strconv.Atoi(arraytemp[1]); err != nil {
			return array, fmt.Errorf("GetPid faild, convert [%s] error,src[%s],%s", arraytemp[1], line, err.Error())
		}
		array = append(array, n)
	}
	return array, nil
}

// 开启路径path的程序进程，显示还是在当前页面中
// 参考使用 go wPipe.StartProcess("D:\\file_svr.exe", []string{"-port=8689"})
func (wp *WindowsPipe) StartProcess(path string, args []string) error {
	wp.lazyTimeout()
	if !FileExists(path) {
		return fmt.Errorf("file [%s] not exist, can't StartProcess", path)
	}

	fullCmd := fmt.Sprintf("start %s %s", path, strings.Join(args, " "))
	cmdPath := string_func.GetOnlyDir(path)
	return winCmd.ExecAysncCmd(cmdPath, fullCmd, wp.timeOutMs)

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
// 通过在cmd中使用start命令开启程序
func (p *WindowsPipe) RunProcess(path string, params []string) (bool, error) {
	err := p.StartProcess(path, params)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 进程processName是否有在运行
func (p *WindowsPipe) IsProcessRun(processName string) (bool, error) {
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
func (p *WindowsPipe) KillProcess(pid int) (bool, error) {
	command := fmt.Sprintf("taskkill /f /im %d", pid)
	err := winCmd.ExecAysncCmd("", command, p.timeOutMs)
	if err != nil {
		return false, err
	}
	return true, nil

}
