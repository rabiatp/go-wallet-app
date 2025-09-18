package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

func jwtSecret() []byte {
  s := os.Getenv("JWT_SECRET")
  if s == "" { s = "dev-secret-change-me" }
  return []byte(s)
}

func GenerateToken(userID string, ttl time.Duration) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "exp": time.Now().Add(ttl).Unix(),
    "iat": time.Now().Unix(),
  }
  t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return t.SignedString(jwtSecret())
}

func ParseToken(tok string) (string, error) {
  p, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { return nil, ErrInvalidToken }
    return jwtSecret(), nil
  })
  if err != nil || !p.Valid { return "", ErrInvalidToken }
  if m, ok := p.Claims.(jwt.MapClaims); ok {
    if sub, ok := m["sub"].(string); ok { return sub, nil }
  }
  return "", ErrInvalidToken
}
