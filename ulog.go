// Copyright 2018 yejiantai Authors
//
// package ulog 基于uber zap的日志库
// 默认记录的详细信息模式，如果要做debug，需要创建的时候开启
package ulog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/liudanyejiantai/gosdk/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DEBUG = 1 // 调试类型
	INFO  = 2 // 信息类型
	WARN  = 3 // 警告类型
	ERROR = 4 // 错误类型
	FATAL = 5 // 致命错误
)

type Ulog struct {
	Logger  *zap.Logger
	Dir     string // 日志文件存储的路径
	Prefix  string // 日志文件的前缀
	logName string // 日志文件完整的名称
	level   int    // 日志文件记录的级别，值越小记录的越详细
	bInit   bool   // 是否已经初始化过
	err     error
}

func NewULog(dir, prefix string) *Ulog {
	if dir == "" {
		dir = file.GetCurDir() + "logs/"
	} else {
		dir = strings.Replace(dir, "\\", "/", -1)
		dir = strings.TrimRight(dir, "/") + "/"
	}

	var ulog = Ulog{Dir: dir, Prefix: prefix, level: INFO}
	return &ulog
}

func (p *Ulog) SetWriteLevel(level int) {
	p.level = level
}

func (p *Ulog) SetLogFile(fileName string) {
	p.logName = fileName
}

func (p *Ulog) initLogger(lp string, lv string, isDebug bool) error {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	var js string
	if isDebug {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stdout"]
      }`, lv)
	} else {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["%s"],
      "errorOutputPaths": ["%s"]
      }`, lv, lp, lp)
	}

	var cfg zap.Config
	if p.err = json.Unmarshal([]byte(js), &cfg); p.err != nil {
		return p.err
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if p.Logger, p.err = cfg.Build(); p.err != nil {
		log.Fatal("init logger error: ", p.err)
		return p.err
	}
	p.bInit = true
	return nil
}

// 2017-03-15 16:07:32.236
func curTime() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d", time.Now().Year(),
		time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond()/1e6)
}

// 写日志,只是单纯的写文件内容,格式内容类似
// 只有log_type的值大于设定的默认级别才可以记录
// {"level":"info","ts":"2018-09-13T13:28:06.659+0800","caller":"log/log.go:72","msg":"信息:","logtext":"试试怎么样[writelog:2]"}
func (p *Ulog) WriteLog(log_type int, str string, fmtArgs ...interface{}) error {
	if p.level > log_type {
		return nil
	}
	os.MkdirAll(p.Dir, 0777)

	strType, bDebug := "", false
	switch log_type {
	case INFO:
		strType = "info"
	case DEBUG:
		strType = "debug"
		bDebug = true
	case WARN:
		strType = "warn"
	case ERROR:
		strType = "error"
	case FATAL:
		strType = "fatal"
	default:
		return fmt.Errorf("日志类型%d不对", strType)
	}

	//initLogger("logs/"+strType+time.Now().Format("20060102")+".log", "DEBUG", false)
	p.logName = p.Dir + p.Prefix + strType + time.Now().Format("20060102") + ".log"
	err := p.initLogger(p.logName, strings.ToUpper(strType), bDebug)
	if err != nil {
		return err
	}

	msg := str
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(str, fmtArgs...)
	}

	fmt.Println(curTime()+" "+strType, msg)

	switch log_type {
	case INFO:
		p.Logger.Info("Info:", zap.String("logtext", msg))
	case DEBUG:
		p.Logger.Debug("Debug:", zap.String("logtext", msg))
	case WARN:
		p.Logger.Warn("Warn:", zap.String("logtext", msg))
	case ERROR:
		p.Logger.Error("Error:", zap.String("logtext", msg))
	case FATAL:
		p.Logger.Fatal("Fatal:", zap.String("logtext", msg))
	default:
		return errors.New(fmt.Sprintf("日志类型%d不对", log_type))
	}

	return nil
}
