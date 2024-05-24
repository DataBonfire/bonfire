package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	issuer = "kolplanet.com"
)

type UserSession struct {
	UserId        uint  `json:"user_id"`
	TokenIssuedAt int64 `json:"token_issued_at"`
}

func GenToken(us *UserSession, jwtKey string) (string, error) {
	j := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strconv.FormatUint(uint64(us.UserId), 10),
		Audience:  jwt.ClaimStrings{},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenString, jwtKey string) (*UserSession, error) {
	if len(tokenString) == 0 || len(jwtKey) == 0 {
		return nil, errors.New("invalid token")
	}
	t := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	mc, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("claims")
	}
	userId, err := strconv.Atoi(mc.Subject)
	if err != nil {
		return nil, err
	}
	if mc.IssuedAt == nil {
		return nil, errors.New("invalid token")
	}
	issueAt := mc.IssuedAt.Unix()

	return &UserSession{
		UserId:        uint(userId),
		TokenIssuedAt: issueAt,
	}, nil
}
