package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte("8b9a31ad-a4f5-4e64-b3d9-8e36ba4d3e1e")
)

type (
	JwtClaim struct {
		jwt.RegisteredClaims
		JwtCustomClaims
	}
	JwtCustomClaims struct {
		Name     string
		ExternId []byte
	}
)

func (jc *JwtClaim) Check() bool {
	// check flow:
	// jti -> [ttlcache|redis] get session meta
	// account,is_admin -> database get info
	return true
}

func GenToken(jcc JwtCustomClaims) (string, error) {
	claim := JwtClaim{
		jwt.RegisteredClaims{
			// should read from config
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			Issuer:    "puzzle_lib",
			ID:        "session_id[uuid]",
		},
		JwtCustomClaims{
			Name:     jcc.Name,
			ExternId: jcc.ExternId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := token.Claims.(*JwtClaim); ok && token.Valid {
		return claim, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
