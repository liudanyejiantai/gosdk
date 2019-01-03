// Copyright 2018 yejiantai Authors
//
// package datatype 类型互相转换
package datatype

import (
	"encoding/binary"
	"math"
	"strconv"
)

func IntToString(n int) string {
	return strconv.Itoa(n)
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func Int64ToString(n64 int64) string {
	return strconv.FormatInt(n64, 10)
}

func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	if buf == nil {
		return -1
	}
	return int64(binary.BigEndian.Uint64(buf))
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func Float64ToString(f64 float64) string {
	return strconv.FormatFloat(f64, 'E', -1, 64)
}

func StringToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
