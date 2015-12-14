package session

import (
	"crypto/sha256"
	"crypto/hmac"
	"crypto/subtle"
	"encoding/base64"
	"strings"
	"errors"
)

var ErrInvalidSignature = errors.New("session: cookie signature is invalid")

func Sign(unsigned, key string) (signed string) {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(unsigned))
	h := mac.Sum(nil)
	b := base64.StdEncoding.EncodeToString(h)
	r := strings.NewReplacer("+", "", "/", "", "=", "")
	return unsigned + "." + r.Replace(b)
}

func Unsign(signed, key string) (unsiged string, err error) {
	splits := strings.Split(signed, ".")
	if len(splits) != 2 {
		return "", ErrInvalidSignature
	}
	if subtle.ConstantTimeCompare([]byte(Sign(splits[0], key)), []byte(signed)) == 1 {
		return splits[0], nil
	}
	return "", ErrInvalidSignature
}
