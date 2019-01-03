// Copyright 2018 yejiantai Authors
//
// package encrypt 加解密包
package encrypt

import (
	"github.com/liudanyejiantai/gosdk/public_func"
	"github.com/liudanyejiantai/gosdk/string_func"
)

// 基于zlib压缩后的加密方法，用做配置信息加密使用
func ZlibEncrypt(str_src string) string {
	var (
		temp string
	)

	temp = string_func.YjtEnc(str_src)
	return public_func.Base64_Encode(public_func.ZlibCompress([]byte(temp)))
}

// 基于zlib压缩后的解密方法，用做配置信息解密使用
func ZlibDecrypt(str_dec string) (string, error) {
	var (
		temp string
		err  error
		byt  []byte
	)

	if temp, err = public_func.Base64_Decode(str_dec); err != nil {
		return "", err
	}
	if byt, err = public_func.ZlibUnCompress([]byte(temp)); err != nil {
		return "", err
	}
	return string_func.YjtDec(string(byt)), nil
}

// 加密方法，用做配置信息加密使用
func Encrypt(str_src string) string {
	var (
		temp string
	)
	temp = string_func.YjtEnc(str_src)
	return public_func.Base64_Encode([]byte(temp))
}

// 解密方法，用做配置信息解密使用
func Decrypt(str_dec string) (string, error) {
	var (
		str string
		err error
	)

	if str, err = public_func.Base64_Decode(str_dec); err != nil {
		return str, err
	}
	return string_func.YjtDec(str), nil
}

// 密码加密使用，密码信息不可反向破解
func PasswordEncrypt(str_src string) string {
	var (
		temp string
	)
	temp = string_func.YjtEnc(str_src)
	temp = public_func.GetMd5String(temp)
	temp = public_func.Base64_Encode([]byte(temp))
	return string_func.ReverseString(temp)
}
