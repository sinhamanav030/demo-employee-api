package token

import (
	"errors"
	"fmt"
	"time"

	"githb.com/demo-employee-api/pkg/customErrors"
	"github.com/golang-jwt/jwt"
)

type Maker interface {
	CreateToken(userId int, role string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", 32)
	}
	return &JwtMaker{secretKey: secretKey}, nil
}

func (maker *JwtMaker) CreateToken(userId int, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, role, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New(customErrors.ErrorInvalidToken)
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errors.New(customErrors.ErrorInvalidToken)) {
			return nil, errors.New(customErrors.ErrorInvalidToken)
		}
		return nil, errors.New(customErrors.ErrorInvalidToken)

	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New(customErrors.ErrorInvalidToken)
	}
	return payload, nil
}
