package com

import (
	"encoding/base64"
	"net/url"
)

//URLEncode url encode string, is + not %20
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

//URLDecode url decode string
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

//Base64Encode base64 encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//Base64Decode base64 decode
func Base64Decode(str string) (string, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return string(s), e
}
