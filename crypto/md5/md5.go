package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
