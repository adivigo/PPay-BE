package lib

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5Hash(hash string) string {
	h := md5.New()
	io.WriteString(h, hash)
	md5Hash := fmt.Sprintf("%x", h.Sum(nil))
	return md5Hash
}
