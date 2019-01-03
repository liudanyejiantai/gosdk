package randStr

import (
	"math/rand"
	"time"
)

const (
	cHAR1    = "1XpasdfBN37xcm890QW2ZECVRvbnTY456UIOPqwertyuioghjklASDFGHJKLzM"
	cHAR2    = "WERQrt23TY6GHUIOPqwe7yuiopasdZXCfghjklASDF45JKLzxcnmVBNM18vb90"
	cHAR3    = "QWwjNM123456ertyuiopasdfgh789klASDFERTYUIOPqGHJKLzxcvbnmZXCVB0"
	cHAR4    = "nmZ345ER789QWPqwertyuiopXCVBNM12asdfghjklASDFGHJKLzxcTYUIO6vb0"
	cHAR5    = "qwertyuiopasdfghjklASDFGHJKLzxcvbnmZXCVBNM1234567890QWERTYUIOP"
	cHAR6    = "asdfghjklASDFGHJKLzxcvbnmZXCVBNM1234567890QWERTYUIOPqwertyuiop"
	cHAR7    = "SDFGHJKLzxcvbnmZXCVBNM1234567890QWERTYUIOPqwertyuiopasdfghjklA"
	cHAR8    = "567890QWEklASDRTYUIOPqwe1234rtyuiopasdfghjFGHJKLzxcvbnmZXCVBNM"
	cHAR9    = "ertyuiopasdfghjklASDFGHJKLzxcvbnmZXCV1234567890QWBNMERTYUIOPqw"
	cHAR10   = "ERTYUIO12jklASDFGnmZXCVBNHJKL3456QW7890PqwertyuiopasdfghzxcvbM"
	cHAR11   = "FGHJ123456QWwertyuiopasdfghjklASD7890KLzxcvbnmZXCVBNMERTYUIOPq"
	cHAR12   = "5690QW1234yuiopasdfghjk78lASDFGZXCVBHJKLzxcvbnmNMERTYUIOPqwert"
	cHAR13   = "56734QWERuiopaTY89012UICVBNMOPqwertysdFGHJKLzxcvfghjklASDbnmZX"
	cHAR14   = "123QWdfghjklASERTY4bnmZXC567890UIOPqweiopasDFGHrtyuJKLzxcvVBNM"
	cHAR15   = "19ZXCV0QWEMRTYU345IOPqwertyuiophjklASDFGHJKLzxcasdfgvbnmBN2678"
	cHAR16   = "DFG2348vbn90qwertQWY567UIOPyuiopas1dfghjklASHJKLzxcmZXCERTVBNM"
	cHAR_LEN = 62
)

var (
	arrayStr = []string{cHAR1, cHAR2, cHAR3, cHAR4, cHAR5, cHAR6, cHAR7, cHAR8, cHAR9, cHAR10, cHAR11, cHAR12, cHAR13, cHAR14, cHAR15, cHAR16}
)

func init() {
	rand.Seed(time.Now().Unix())
}

// 获得32位的随机字符串
func GetRand32Str() string {
	n := cHAR_LEN
	var arrayInt []int
	var i int = 0
	for i = 0; i < 16; i++ {
		arrayInt = append(arrayInt, rand.Intn(n))
	}

	var str string
	for i = 0; i < 16; i++ {
		str += arrayStr[i][arrayInt[i] : arrayInt[i]+1]
	}
	for i = 15; i >= 0; i-- {
		str += arrayStr[i][arrayInt[i] : arrayInt[i]+1]
	}
	// 需要针对str再做顺序打散
	return str
}
