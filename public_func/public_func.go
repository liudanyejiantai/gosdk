// Copyright 2018 yejiantai Authors
//
// package public_func 公共方法
package public_func

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	rd "math/rand"
	"net"
	"os"
	"time"
)

func init() {
	rd.Seed(time.Now().Unix())
}

// 获得nstart到nEnd之间的随机数
func GetRand(nstart, nEnd int) int {
	// int随机值，返回值为int
	return nstart + rd.Intn(nEnd-nstart)
}

// 获得本机IP地址
func GetLocalIP() string {
	var (
		addrs   []net.Addr
		err     error
		address net.Addr
		ipnet   *net.IPNet
		ok      bool
	)

	if addrs, err = net.InterfaceAddrs(); err != nil {
		return ""
	}
	for _, address = range addrs {
		if ipnet, ok = address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 对一个字符串进行MD5加密,不可解密
func GetMd5String(s string) string {
	var (
		h hash.Hash
	)
	h = md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

// 判断文件是否存在  存在返回 true 不存在返回false
func CheckFileIsExist(filename string) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}
	// 如果是文件夹，返回失败
	if stat.IsDir() {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// 进行zlib压缩
func ZlibCompress(src []byte) []byte {
	var (
		in bytes.Buffer
		w  *zlib.Writer
	)
	w = zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// 进行zlib解压缩
func ZlibUnCompress(compress []byte) ([]byte, error) {
	var (
		b   *bytes.Reader
		out bytes.Buffer
		err error
		r   io.ReadCloser
	)
	b = bytes.NewReader(compress)
	if r, err = zlib.NewReader(b); err != nil {
		return []byte(""), err
	}
	io.Copy(&out, r)
	return out.Bytes(), nil
}

// Base64解码
func Base64_Decode(base64_str string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64_str)
	return string(decoded), err
}

// Base64编码
func Base64_Encode(byt_src []byte) string {
	return base64.StdEncoding.EncodeToString(byt_src)
}
