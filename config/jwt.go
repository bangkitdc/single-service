package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY4ODcxMzcyNiwiaWF0IjoxNjg4NzEzNzI2fQ.-ySqrsyecVvpuixK3GPyoDHsfmJJnYTUiKBGYxox7Mk")

type JWTClaim struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}
