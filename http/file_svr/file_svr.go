// Copyright 2018 yejiantai Authors
//
// golang的http文件服务,可以指定路径监听发布，也可以全盘监听
// 支持指定路径上传，默认上传到当前目录
// package file_svr
package file_svr

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/liudanyejiantai/gosdk/public_func"

	"github.com/liudanyejiantai/gosdk/datatype"
	fl "github.com/liudanyejiantai/gosdk/file"
	"github.com/liudanyejiantai/gosdk/ulog"
)

var (
	log *ulog.Ulog
)

func init() {
	log = ulog.NewULog("", "")
}

var (
	g_strUpPath  string // 默认上传目录
	g_strDwPath  string // 默认下载目录
	g_nSvrPort   int    // http服务的端口
	g_strUpRoute string // 上传的路由
	g_strDwRoute string // 下载路由
)

const (
	maxUploadSize = int64(1024 * 1024 * 4) // 单个文件最大4Mb
)

// 设置文件服务的连接信息
func SetFileSvrInfo(nSvrPort int, strUploadPath, strDownPath, strUpRoute, strDwRoute string) {
	g_strUpPath, g_strDwPath, g_nSvrPort = strUploadPath, strDownPath, nSvrPort
	g_strUpRoute, g_strDwRoute = strUpRoute, strDwRoute
}

// 开启监听服务
func StartSvr() {
	hostAndPort := ":" + datatype.IntToString(g_nSvrPort)

	log.WriteLog(ulog.INFO, "开启文件上传服务成功,监听端口为[%s]", hostAndPort)
	http.HandleFunc(g_strUpRoute, UpFileHandle)
	http.HandleFunc(g_strDwRoute, DwFileHandle)

	http.ListenAndServe(hostAndPort, nil)
}

// 返回错误信息
func writeError(w http.ResponseWriter, r *http.Request, err_msg string, err_code int, err error) {
	w.WriteHeader(err_code)
	w.Write([]byte(err_msg))
	if err != nil {
		log.WriteLog(ulog.ERROR, "[%s]http返回[%s],返回码[%d],原因[%s]", r.RemoteAddr, err_msg, err_code, err.Error())
	} else {
		log.WriteLog(ulog.ERROR, "[%s]http返回[%s],返回码[%d]", r.RemoteAddr, err_msg, err_code)
	}
}

// 上传路由,
// curl http://110.110.110.47:12346/up/ -F "local_path=@SunTransClient.zip" -F "save_path=D:\var" -F "save_name=new.zip" -v
// 页面上传处理 <input type
func UpFileHandle(w http.ResponseWriter, r *http.Request) {
	var (
		err                        error
		filetype, uuid, local_name string

		fileBytes                       []byte
		save_path, save_name, save_full string
		newFile                         *os.File
		file                            multipart.File
	)
	if strings.Split(r.RemoteAddr, ":")[0] == "127.0.0.1" {
		writeError(w, r, "remote ip error", http.StatusBadRequest, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err = r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, r, "FILE_TOO_BIG, max size is"+datatype.Int64ToString(maxUploadSize), http.StatusBadRequest, err)
		return
	}

	if file, _, err = r.FormFile("local_path"); err != nil {
		writeError(w, r, "get file path error", http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	if fileBytes, err = ioutil.ReadAll(file); err != nil {
		writeError(w, r, "read file error", http.StatusBadRequest, err)
		return
	}
	filetype = http.DetectContentType(fileBytes)
	// 检测文件类型
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" &&
		filetype != "image/tif" && filetype != "image/tiff" {
		writeError(w, r, "INVALID_FILE_TYPE "+filetype, http.StatusBadRequest, err)
		return
	}
	save_path, save_name = r.FormValue("save_path"), r.PostFormValue("save_name")

	// 如果上传没有指定路径，就在当前目录下面
	uuid = public_func.GetGuid()
	if save_path == "" {
		save_path = g_strUpPath + "/" + uuid[:4] + "/" + uuid[4:8] + "/" + uuid[8:12]
		fl.CreateDirTree(save_path)
	}
	// 如果文件名称没有指定，就使用uuid创建
	if save_name == "" {
		save_name = uuid
	}

	save_full = save_path + "/" + save_name
	if newFile, err = os.Create(save_full); err != nil {
		writeError(w, r, "CANT_WRITE_FILE1 "+save_full, http.StatusInternalServerError, err)
		return
	}
	defer newFile.Close()
	if _, err = newFile.Write(fileBytes); err != nil {
		writeError(w, r, "CANT_WRITE_FILE2 "+save_full, http.StatusInternalServerError, err)
		return
	}
	//url = g_strDwRoute + save_full
	w.Write([]byte("SUCCESS"))
	log.WriteLog(ulog.INFO, "客户端[%s]请求上传文件[%s],返回字符串为[]", r.RemoteAddr, local_name)
	return

}

// 下载路由
func DwFileHandle(w http.ResponseWriter, r *http.Request) {

	return

}
