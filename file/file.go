// Copyright 2018 yejiantai Authors
//
// package file 文件操作
package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// func 获得当前目录,处理后带/符号
func GetCurDir() string {
	var (
		dir string
		err error
	)
	if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return ""
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	dir = strings.TrimRight(dir, "/") + "/"
	return dir
}

// 生成目录树
func CreateDirTree(strDir string) error {
	return os.MkdirAll(strDir, 0777)
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

// 确保目录存在
func MakeFileExist(filename string) bool {
	if err := os.MkdirAll(filename, 0777); err != nil {
		return false
	}

	return true
}

// 遍历文件目录获得目录
func LoopGetDirList(path string) ([]string, error) {
	var (
		arr_list []string
		err_     error
	)
	err_ = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			path = strings.Replace(path, "\\", "/", -1)
			arr_list = append(arr_list, path)
			return nil
		}
		return nil
	})

	return arr_list, err_
}

// 遍历文件目录获得文件
func LoopGetFileList(path string) ([]string, error) {
	var (
		arr_list []string
		err_     error
	)
	err_ = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		path = strings.Replace(path, "\\", "/", -1)
		arr_list = append(arr_list, path)
		return nil
	})

	return arr_list, err_
}

// 获取文件大小
func GetFileSize(strFileName string) (int64, error) {
	var (
		file *os.File
		err  error
		stat os.FileInfo
	)

	if file, err = os.OpenFile(strFileName, os.O_RDWR, 0666); err != nil {
		return int64(-1), err
	}

	defer file.Close()
	if stat, err = file.Stat(); err != nil {
		return int64(-1), err
	}

	return stat.Size(), nil
}

// 覆盖保存文件
func SaveFile(bytBuf []byte, strFileName string) error {
	var (
		file *os.File
		err  error
	)

	if file, err = os.OpenFile(strFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666); err != nil {
		return err
	}

	//这种方式才会覆盖写入文件
	defer file.Close()
	if _, err = file.Write(bytBuf); err != nil {
		return err
	}
	return nil
}

// 追加保存文件
func SaveAppendFile(bytBuf []byte, strFileName string) error {
	var (
		file     *os.File
		err      error
		uFileLen int64
	)

	if file, err = os.OpenFile(strFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666); err != nil {
		return err
	}

	defer file.Close()
	if uFileLen, err = file.Seek(0, os.SEEK_END); err != nil {
		return err
	}

	//这种方式才会覆盖写入文件
	if _, err = file.WriteAt(bytBuf, uFileLen); err != nil {
		return err
	}
	return nil
}

// 根据开始结束位置读取文件
func ReadFileBufferByBlock(strFileName string, u64StartPos, u64EndPos int64) ([]byte, error) {
	var (
		byt      = []byte("")
		err      error
		file     *os.File
		nLen     int
		nReadLen int
	)
	if !FileExists(strFileName) {
		return byt, errors.New("文件" + strFileName + "不存在")
	}

	if file, err = os.OpenFile(strFileName, os.O_RDWR, os.ModeType); err != nil {
		return byt, errors.New("读取文件" + strFileName + "失败,原因" + err.Error())
	}

	defer file.Close()
	if _, err = file.Seek(u64StartPos, 0); err != nil {

		return byt, errors.New("读取文件" + strFileName + "失败,指定位置失败原因" + err.Error())
	}

	nLen, nReadLen = int(u64EndPos-u64StartPos), 0
	buf := make([]byte, nLen)
	if nReadLen, err = file.Read(buf); err != nil || nReadLen != nLen {
		return byt, errors.New("读取文件" + strFileName + "失败,读取长度信息失败原因" + err.Error())
	}
	return buf, err
}

// 读取文件
func ReadFileBuffer(strFileName string) ([]byte, error) {
	var (
		file *os.File
		buf  = []byte("")
		err  error
	)
	if !FileExists(strFileName) {
		return buf, errors.New("文件" + strFileName + "不存在")
	}

	if file, err = os.OpenFile(strFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeType); err != nil {
		return buf, errors.New("读取文件" + strFileName + "失败,原因" + err.Error())
	}

	defer file.Close()
	if buf, err = ioutil.ReadAll(file); err != nil {
		return buf, errors.New("打开文件" + strFileName + "失败,原因" + err.Error())
	}
	return buf, err
}

// 删除文件
func RemoveFile(strFileName string) error {
	var (
		err error
	)

	if err = os.Remove(strFileName); err != nil {
		return err
	}
	return nil
}

// 清空目录树
func RemoveDirTree(strPath string) error {
	return os.RemoveAll(strPath)
}

// 指定类型遍历文件夹,strPath表示文件目录,suffixs表示后缀类型数组
func GetSuffixList(strPath string, suffixs []string) ([]string, error) {
	var (
		err      error
		arr_file []string
		suffix   string
	)
	err = filepath.Walk(strPath, func(strPath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() { //文件夹
			return nil
		}
		for _, suffix = range suffixs {
			if strings.HasSuffix(strPath, suffix) {
				arr_file = append(arr_file, strPath)
			}
		}

		return nil
	})
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("filepath.Walk(%s) faild, %s", strPath, err.Error()))
	}
	return arr_file, nil
}

// 遍历strDir目录树清除指定后缀的文件
func RemoveSuffixTree(strDir string, suffixs []string) error {
	var (
		err       error
		file_list []string
		file_name string
	)

	if file_list, err = GetSuffixList(strDir, suffixs); err != nil {
		return err
	}
	for _, file_name = range file_list {

		if err = os.Remove(file_name); err != nil {
			return err
		}
	}

	return nil
}

// 获得文件strFile的存在时间
func GetFileExistsSec(strFile string) (int64, error) {
	var (
		err      error
		fileInfo os.FileInfo
		sec      int64
		now_sec  int64
	)
	if fileInfo, err = os.Stat(strFile); err != nil {
		return -1, err
	}

	sec = fileInfo.ModTime().Unix() //创建的秒
	now_sec = time.Now().Unix()
	return (now_sec - sec), nil
}

//获得文件夹strDir下时效超过nExpireSec秒的全部文件
func GetExpireList(strDir string, nExpireSec int64) ([]string, error) {
	var (
		file_list  []string
		rt_list    []string
		err        error
		file_name  string
		nExistTime int64
	)
	if file_list, err = LoopGetFileList(strDir); err != nil {
		return []string{}, err
	}

	for _, file_name = range file_list {
		if nExistTime, err = GetFileExistsSec(file_name); err != nil {
			return []string{}, err
		}
		if nExistTime > nExpireSec {
			rt_list = append(rt_list, file_name)
		}
	}

	return rt_list, nil
}
