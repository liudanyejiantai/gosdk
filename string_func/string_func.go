// Copyright 2018 yejiantai Authors
//
// package string_func 字符串处理
package string_func

import (
	"fmt"
	"strings"
)

var (
	// 原始字典表
	SRC_CHARSETS = []string{"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=",
		"~", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "\"",
		"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "[", "]", "\\",
		"{", "}", "|", "A", "S", "D", "F", "G", "H", "J", "K", "L", ";",
		"'", ":", "+", "Z", "X", "C", "V", "B", "N", "M", ".", "/", "<",
		">", "?", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a",
		"s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b",
		"n", "m"}
	// 转换后的字典表
	DST_CHARSETS = []string{"v", "b", "n", "p", "a", "s", "d", "f", "g", "h", "j", "k", "c",
		"(", ")", "_", "\"", "`", "1", "2", "3", "4", "5", "6", "7", "8",
		"9", "0", "-", "=", "Q", "W", "E", "R", "T", "Y", "O", "P", "[",
		"]", "\\", "{", "}", "|", "A", "S", "D", "F", "G", "H", "J", "K",
		"L", ";", "'", ":", "+", "Z", "X", "C", "V", "B", "N", "M", ".",
		"/", "<", ">", "?", "q", "w", "e", "r", "t", "y", "u", "i", "o",
		"m", "I", "U", "@", "#", "$", "%", "^", "&", "*", "l", "z", "x",
		"!", "~"}

	// 纯粹字母表
	SRC_CHAR_ORDER = []string{"0", "1", "2", "3", "4", "5", "6", "7",
		"8", "9", "a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z"}

	// 随机字母表
	DST_CHAR_HASH_RAND = []string{"N", "X", "F", "Z", "e", "1", "D", "K",
		"k", "5", "m", "T", "o", "b", "q", "r",
		"I", "J", "3", "L", "A", "7", "O", "P",
		"V", "B", "H", "2", "E", "Y", "G", "C",
		"s", "t", "u", "v", "w", "8", "y", "z",
		"S", "l", "6", "W", "x", "9", "a", "p",
		"c", "d", "0", "f", "g", "h", "i", "j",
		"Q", "R", "4", "n", "U", "M"}
)

// 通过跑字典表实现字符串和数字类型替换,做简单加密用
func LetterEnc(str_src string) string {
	var (
		i           int
		str_dst     string
		str_replace string
		new_str     string
		ok          bool
	)
	m := make(map[string]string)
	for i = 0; i < len(SRC_CHAR_ORDER); i++ {
		m[SRC_CHAR_ORDER[i]] = DST_CHAR_HASH_RAND[i]
	}

	for i = 0; i < len(str_src); i++ {
		new_str, ok = m[str_src[i:i+1]]
		if ok {
			str_replace = new_str
		} else {
			str_replace = str_src[i : i+1]
		}
		str_dst += str_replace
	}
	return str_dst
}

// 通过跑字典表字符串和数字类型串替换,做简单解密用
func LetterDec(str_dec string) string {
	m := make(map[string]string)
	for i := 0; i < len(DST_CHAR_HASH_RAND); i++ {
		m[DST_CHAR_HASH_RAND[i]] = SRC_CHAR_ORDER[i]
	}
	str_return, str_replace := "", ""

	for i := 0; i < len(str_dec); i++ {
		new_str, ok := m[str_dec[i:i+1]]
		if ok {
			str_replace = new_str
		} else {
			str_replace = str_dec[i : i+1]
		}
		str_return += str_replace
	}
	return str_return
}

// 通过跑字典表实现字符串替换,做简单加密用
func YjtEnc(str_src string) string {
	var (
		i           int
		str_dst     string
		str_replace string
		new_str     string
		ok          bool
	)
	m := make(map[string]string)
	for i = 0; i < len(SRC_CHARSETS); i++ {
		m[SRC_CHARSETS[i]] = DST_CHARSETS[i]
	}

	for i = 0; i < len(str_src); i++ {
		new_str, ok = m[str_src[i:i+1]]
		if ok {
			str_replace = new_str
		} else {
			str_replace = str_src[i : i+1]
		}
		str_dst += str_replace
	}
	return str_dst
}

// 通过跑字典表实现字符串替换,做简单解密用
func YjtDec(str_dec string) string {
	m := make(map[string]string)
	for i := 0; i < len(DST_CHARSETS); i++ {
		m[DST_CHARSETS[i]] = SRC_CHARSETS[i]
	}
	str_return, str_replace := "", ""

	for i := 0; i < len(str_dec); i++ {
		new_str, ok := m[str_dec[i:i+1]]
		if ok {
			str_replace = new_str
		} else {
			str_replace = str_dec[i : i+1]
		}
		str_return += str_replace
	}
	return str_return
}

// 去除左边的空格等特殊字符
func TrimLeft(str string) string {
	str_return := ""
	var i = 0
	for i = 0; i < len(str); i++ {
		if str[i:i+1] != " " &&
			str[i:i+1] != "\r" &&
			str[i:i+1] != "\n" &&
			str[i:i+1] != "\t" {
			break
		}
	}
	str_return = str[i:]

	return str_return
}

// 去除右边的空格等特殊字符
func TrimRight(str string) string {
	return ReverseString(TrimLeft(ReverseString(str)))
}

// 去除左右空格
func Trim(str string) string {
	return TrimLeft(TrimRight(str))
}

// 从右边开始用sep对str分割成count部分，如果分割不足返回error
func SplitRight(str, sep string, count int) ([]string, error) {
	resStr := ReverseString(str)
	resSep := ReverseString(sep)
	array, err := SplitLeft(resStr, resSep, count)
	if err != nil {
		return []string{str}, fmt.Errorf("split by sep [%s] don't have [%d]number", sep, count)
	}
	var arrayReturn []string
	for i := len(array) - 1; i >= 0; i-- {
		arrayReturn = append(arrayReturn, ReverseString(array[i]))
	}

	return arrayReturn, nil
}

// 从左边开始用sep分隔符针对str做分割，分割成count部分，如果分割的数目不满足count返回error
func SplitLeft(str, sep string, count int) ([]string, error) {
	arrayTemp := strings.Split(str, sep)
	if len(arrayTemp) < count {
		return []string{str}, fmt.Errorf("split by sep [%s] don't have [%d]number", sep, count)
	}

	var array []string
	var i int
	var temp string
	for i = 0; i < count-1; i++ {
		array = append(array, arrayTemp[i])
	}

	for i = count - 1; i < len(arrayTemp)-1; i++ {
		temp += arrayTemp[i] + sep
	}
	temp += arrayTemp[i]
	array = append(array, temp)

	return array, nil
}

func SplitAll(str, sep string) []string {
	arrayInt := IndexAll(str, sep)
	if len(arrayInt) == 1 {
		if arrayInt[0] == -1 {
			return []string{str}
		}

	}

	var arrayStr []string

	return arrayStr
}

// 将array中相同的数据清空，确保数据不存在重复的
func ArrayRemoveRepeat(array []string) []string {
	var arrayTemp []string
	var exists bool = false
	for i := 0; i < len(array); i++ {
		if i == 0 {
			arrayTemp = append(arrayTemp, array[i])
			continue
		}

		exists = false
		for j := 0; j < len(arrayTemp); j++ {
			if arrayTemp[j] == array[i] {
				exists = true
			}
		}
		if !exists {
			arrayTemp = append(arrayTemp, array[i])
		}
	}

	return arrayTemp
}

// 获得字符串sep在原始字符串str中出现的全部下标
func IndexAll(str, sep string) []int {
	var array []int
	if !strings.Contains(str, sep) || sep == "" {
		return []int{-1}
	}

	var temp string = str
	var pos, tempPos, sepCount int = 0, -1, len(sep)
	for {
		if tempPos = strings.Index(temp, sep); tempPos >= 0 {
			temp = temp[tempPos+sepCount:]
			fmt.Println("temp=[" + temp + "]")
			if pos == 0 {
				pos += tempPos
			} else {
				pos += tempPos + sepCount
			}
			array = append(array, pos)
		}
		if !strings.Contains(temp, sep) {
			break
		}
	}

	return array
}

// 将多空格转换为单空格
func ConvertToSingleTrim(str string) string {
	temp := str
	temp = strings.Replace(temp, "\t", " ", -1)
	temp = strings.Replace(temp, "\r", " ", -1)
	temp = strings.Replace(temp, "\n", " ", -1)
	temp = strings.Replace(temp, "\v", " ", -1)
	temp = strings.Replace(temp, "\a", " ", -1)
	temp = strings.Replace(temp, "\f", " ", -1)
	for {
		if !strings.Contains(temp, "  ") {
			break
		}
		temp = strings.Replace(temp, "  ", " ", -1)
	}
	return temp
}

// 反转字符串函数
func ReverseString(str string) string {
	runes := []rune(str)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// 将十进制转换为62进制  0-9a-zA-Z 六十二进制
func TransTo62(id int64) string {
	// 1 -- > 1
	// 10-- > a
	// 61-- > Z
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var shortUrl []byte
	for {
		var result byte
		number := id % 62
		result = charset[number]
		var tmp []byte
		tmp = append(tmp, result)
		shortUrl = append(tmp, shortUrl...)
		id = id / 62
		if id == 0 {
			break
		}
	}
	return string(shortUrl)
}

// 获得单纯的路径
func GetOnlyDir(str string) string {
	temp := strings.Replace(str, "\\", "/", -1)
	array := strings.Split(temp, "/")
	if len(array) < 2 {
		return str
	}
	var ret string
	for i := 0; i < len(array)-1; i++ {
		ret += array[i] + "/"
	}
	return ret
}
