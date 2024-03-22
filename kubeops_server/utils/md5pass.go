package utils

import (
	"crypto/md5"
	"fmt"
)

func PasswordToMd5(password string) string {
	md5password := md5.New()
	md5password.Write([]byte(password))
	return fmt.Sprintf("%x", md5password.Sum(nil))
}
