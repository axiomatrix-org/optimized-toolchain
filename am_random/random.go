package am_random

import (
	"math/rand"
	"time"
)

func CreateRandomDigits(length uint) string {
	const digits = "0123456789"                                // 划定取值范围
	result := make([]byte, length)                             // 定义结果集
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano())) // 种子
	for i := range result {
		result[i] = digits[r.Intn(len(digits))] // 种子随机产生取值范围长度之内的下标而取值
	}
	return string(result)
}

func CreateRandomString(length uint) string {
	const charas = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = charas[r.Intn(len(charas))]
	}
	return string(result)
}
