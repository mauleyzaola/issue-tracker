package tecutils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Encrypt(text string) (result string) {
	h := md5.New()
	io.WriteString(h, text)

	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))

	// Specify two salt: salt1 = @#$% salt2 = ^&*()
	salt1 := "@#$%"
	salt2 := "^&*()"

	// salt1 + username + salt2 + MD5 splicing
	io.WriteString(h, salt1)
	io.WriteString(h, salt2)
	io.WriteString(h, pwmd5)

	result = fmt.Sprintf("%x", h.Sum(nil))
	return
}
