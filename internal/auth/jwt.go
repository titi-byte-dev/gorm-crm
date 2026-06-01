package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

// Claims são os dados que vivem dentro do JWT.
// Quanto menos dados aqui, menor o token e menor a superfície de exposição.
// NUNCA colocar dados sensíveis (password, dados de pagamento) no JWT —
// o payload é apenas Base64, não encriptado (qualquer um pode ler).
type Claims struct {
	UserID string    `json:"uid"`
	OrgID  string    `json:"org_id"`
	Role   user.Role `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair é o par access + refresh devolvido no login.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // segundos
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// Só em desenvolvimento — em produção JWT_SECRET é obrigatório
		s = "dev-secret-change-in-production"
	}
	return []byte(s)
}

// GenerateTokenPair cria um access token (curto) e um refresh token (longo).
//
// Access token  — 24h — enviado em cada request no header Authorization.
// Refresh token — 7d  — só usado para pedir um novo access token.
//
// Porquê dois tokens?
// Se o access token for roubado, expira em 24h sem ação do utilizador.
// O refresh token tem vida mais longa mas é usado raramente, reduzindo
// a janela de exposição.
func GenerateTokenPair(userID, orgID uuid.UUID, role user.Role) (*TokenPair, error) {
	accessExp := time.Now().Add(24 * time.Hour)
	access, err := generateToken(userID, orgID, role, accessExp)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refresh, err := generateToken(userID, orgID, role, refreshExp)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(time.Until(accessExp).Seconds()),
	}, nil
}

func generateToken(userID, orgID uuid.UUID, role user.Role, exp time.Time) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		OrgID:  orgID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gorm-crm",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// ValidateToken valida a assinatura e a expiração do token.
// Devolve os Claims para o handler/middleware extrair userID e role.
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		// Garante que o algoritmo é o esperado — previne algorithm confusion attacks
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
