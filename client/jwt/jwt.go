package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type TokenPair struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

type Claims[T any] struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type"`
	Payload   T      `json:"payload"`
}

type Client[T any] interface {
	Generate(payload T) (string, error)
	GenerateRefresh(payload T) (string, error)
	GeneratePair(payload T) (*TokenPair, error)
	Refresh(refreshToken string) (*TokenPair, error)
	Parse(tokenString string) (*Claims[T], error)
	ParseRefresh(tokenString string) (*Claims[T], error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

var _ Client[any] = (*client[any])(nil)

type client[T any] struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New[T any](secret, issuer string, accessTTL time.Duration, refreshTTL ...time.Duration) *client[T] {
	refreshTokenTTL := accessTTL * 30
	if len(refreshTTL) > 0 {
		refreshTokenTTL = refreshTTL[0]
	}

	return &client[T]{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTokenTTL,
	}
}

func (s *client[T]) Generate(payload T) (string, error) {
	return s.generate(payload, accessTokenType, s.accessTTL)
}

func (s *client[T]) GenerateRefresh(payload T) (string, error) {
	return s.generate(payload, refreshTokenType, s.refreshTTL)
}

func (s *client[T]) GeneratePair(payload T) (*TokenPair, error) {
	var (
		accessToken  string
		refreshToken string
		err          error
	)

	if accessToken, err = s.Generate(payload); err != nil {
		return nil, err
	}

	if refreshToken, err = s.GenerateRefresh(payload); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenType:             "Bearer",
		ExpiresIn:             int64(s.accessTTL.Seconds()),
		RefreshTokenExpiresIn: int64(s.refreshTTL.Seconds()),
	}, nil
}

func (s *client[T]) Refresh(refreshToken string) (*TokenPair, error) {
	var (
		claims *Claims[T]
		err    error
	)

	if claims, err = s.ParseRefresh(refreshToken); err != nil {
		return nil, err
	}

	return s.GeneratePair(claims.Payload)
}

func (s *client[T]) Parse(tokenString string) (*Claims[T], error) {
	return s.parse(tokenString, accessTokenType)
}

func (s *client[T]) ParseRefresh(tokenString string) (*Claims[T], error) {
	return s.parse(tokenString, refreshTokenType)
}

func (s *client[T]) AccessTTL() time.Duration {
	return s.accessTTL
}

func (s *client[T]) RefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *client[T]) generate(payload T, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := Claims[T]{
		TokenType: tokenType,
		Payload:   payload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *client[T]) parse(tokenString string, tokenType string) (*Claims[T], error) {
	var (
		token  *jwt.Token
		claims *Claims[T]
		err    error
		ok     bool
	)

	if token, err = jwt.ParseWithClaims(
		tokenString,
		&Claims[T]{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return s.secret, nil
		},
		jwt.WithIssuer(s.issuer),
	); err != nil {
		return nil, err
	}

	if claims, ok = token.Claims.(*Claims[T]); !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.TokenType != tokenType {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
