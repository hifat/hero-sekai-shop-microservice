package jwtauth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

var ErrTokenMalformed = jwt.ErrTokenMalformed
var ErrTokenExpired = jwt.ErrTokenExpired
var ErrInvalidToken = "invalid token"

type (
	AuthFactory interface {
		SignToken() string
	}

	Claims struct {
		PlayerId string `json:"player_id"`
		RoleCode int32  `json:"role_code"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
	apiKey       struct{ *authConcrete }
)

func (a *authConcrete) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	ss, _ := token.SignedString(a.Secret)

	return ss
}

func jwtTimeDurationCal(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}

func jwtTimeRepeatAdapter(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(int64(t), 0))
}

func NewAccessToken(secret string, expiredAt time.Duration, claims *Claims) AuthFactory {
	return &accessToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "rust-better-than-go",
					Subject:   "access-token",
					Audience:  []string{"i-love-go"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}

}

func NewRefreshToken(secret string, expiredAt time.Duration, claims *Claims) AuthFactory {
	return &accessToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "rust-better-than-go",
					Subject:   "refresh-token",
					Audience:  []string{"i-love-go"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}

}

func ReloadToken(secret string, expiredAt time.Duration, claims *Claims) string {
	obj := &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "rust-better-than-go",
					Subject:   "refresh-token",
					Audience:  []string{"i-love-go"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}

	return obj.SignToken()
}

func NewApiKey(secret string, expiredAt time.Duration, claims *Claims) AuthFactory {
	return &apiKey{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "rust-better-than-go",
					Subject:   "api-key",
					Audience:  []string{"i-love-go"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("expected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

var apiKeyInstant string
var once sync.Once

func SetApiKey(secret string) {
	once.Do(func() {
		now := time.Now()
		d := now.AddDate(1, 0, 0).Sub(now)
		apiKeyInstant = NewApiKey(secret, d, nil).SignToken()
	})
}

func SetApiKeyInContext(ctx *context.Context) {
	*ctx = metadata.NewOutgoingContext(*ctx, metadata.Pairs("x-api-key", apiKeyInstant))
}
