// Package security implements authentication and authorization
package security

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*UserClaims, error)
	ValidateServiceToken(ctx context.Context, token string) (*ServiceClaims, error)
}

type UserClaims struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Roles     []string  `json:"roles"`
	ExpiresAt time.Time `json:"exp"`
	jwt.RegisteredClaims
}

type ServiceClaims struct {
	ServiceID string    `json:"service_id"`
	Scopes    []string  `json:"scopes"`
	ExpiresAt time.Time `json:"exp"`
	jwt.RegisteredClaims
}

type JWTAuthService struct {
	secretKey []byte
}

func NewJWTAuthService(secretKey string) AuthService {
	return &JWTAuthService{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTAuthService) ValidateToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if time.Now().After(claims.ExpiresAt) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (j *JWTAuthService) ValidateServiceToken(ctx context.Context, tokenString string) (*ServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid service token: %w", err)
	}

	if claims, ok := token.Claims.(*ServiceClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid service token claims")
}

// Middleware for HTTP authentication
type AuthMiddleware struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Context keys for storing auth information
type contextKey string

const (
	UserClaimsKey    contextKey = "user_claims"
	ServiceClaimsKey contextKey = "service_claims"
)

func WithUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, UserClaimsKey, claims)
}

func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims, ok
}

func GetTenantIDFromContext(ctx context.Context) string {
	if claims, ok := GetUserClaims(ctx); ok {
		return claims.TenantID
	}
	return ""
}

func GetUserIDFromContext(ctx context.Context) string {
	if claims, ok := GetUserClaims(ctx); ok {
		return claims.UserID
	}
	return ""
}