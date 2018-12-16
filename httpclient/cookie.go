package httpclient

import (
	"net/http"
	"time"
)

func NewCookie(name, value string, age time.Duration, path string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(age)
	cookie.Path = path
	return cookie
}
