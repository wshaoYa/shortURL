package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMd5 生成md5
func GetMd5(bs []byte) string {
	h := md5.New()
	h.Write(bs)
	return hex.EncodeToString(h.Sum(nil))
}
