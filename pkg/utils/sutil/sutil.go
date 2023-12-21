package sutil

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"
)

var (
	sumString = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randS     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 效率不如 == 我: <- _ <-
//// StringCompareByHash 先使用 xxHash处理字符串来比较hash 然后再一个一个比较
//func StringCompareByHash(v1, v2 string) bool {
//	if len(v1) != len(v2) {
//		return false
//	}
//
//	v1Hash := xxhash.Sum64String(v1)
//	v2Hash := xxhash.Sum64String(v2)
//
//	if v1Hash != v2Hash {
//		return false
//	}
//
//	// cast to []byte
//	bv1 := []byte(v1)
//	bv2 := []byte(v2)
//
//	for i, v := range bv1 {
//		if v != bv2[i] {
//			return false
//		}
//	}
//
//	return true
//}

// GenerateRandomString
// 生成一个int63 并且拆分成8个uint8(可以表示256)
// 由此来实现高效的生成一个随机字符串
func GenerateRandomString(length int) string {
	if length == 0 || length > math.MaxUint16 {
		return ""
	}

	var (
		num      int64
		res      = make([]byte, length)
		uint8s   [8]uint8
		resIndex = 0
	)

	for i := 0; i < length/len(uint8s)+1; i++ {
		num = randS.Int63()

		// 使用 binary.LittleEndian 将 int64 转换为字节数组
		binary.LittleEndian.PutUint64(uint8s[:], uint64(num))

		for _, v := range uint8s {
			// 这里使用 % 是安全的处理
			res[resIndex] = sumString[int(v)%len(sumString)]
			resIndex++

			if resIndex == length {
				return string(res)
			}
		}
	}

	return ""
}
