package pipe

type Pipe interface {
	// 设置超时时间
	SetTimeOutMs(timeOutMs int64)

	// 进程processName是否有在运行
	IsProcessRun(processName string) (bool, error)

	// 带参数cmd开启path可执行程序
	RunProcess(path string, params []string) (bool, error)

	// 获得进程名对应的进程ID
	GetPid(processName string) ([]int, error)

	// 根据进程ID关闭程序
	KillProcess(pid int) (bool, error)

	// 根据进程名称关闭程序
	KillProcessName(processName string) (bool, error)
}
