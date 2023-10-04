package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
