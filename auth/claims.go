package auth

import "github.com/golang-jwt/jwt/v5"

// Claims is the payload we embed inside every JWT.
// UserID and Email identify the caller; the embedded RegisteredClaims
// carries standard fields like expiry (exp) and issued-at (iat).
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
