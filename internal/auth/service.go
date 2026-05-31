package auth

import (
	"fmt"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

type Service struct {
	users user.Repository
}

func NewService(users user.Repository) *Service {
	return &Service{users: users}
}

type RegisterDTO struct {
	Name     string    `json:"name"     validate:"required,min=2,max=100"`
	Email    string    `json:"email"    validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=72"`
	Role     user.Role `json:"role"     validate:"required,oneof=admin manager seller"`
}

type LoginDTO struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *Service) Register(dto RegisterDTO) (*user.User, error) {
	// Verificar se o email já existe antes de fazer hash (evita trabalho desnecessário)
	existing, err := s.users.FindByEmail(dto.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already registered: %w", sharederrors.ErrConflict)
	}

	hash, err := HashPassword(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &user.User{
		ID:           uuid.New(),
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: hash,
		Role:         dto.Role,
	}

	return s.users.Save(u)
}

// Login valida as credenciais e devolve um par de tokens JWT.
//
// Segurança importante: devolvemos sempre o mesmo erro genérico para
// email não encontrado E password errada. Mensagens diferentes permitiriam
// a um atacante descobrir quais emails existem no sistema (user enumeration).
func (s *Service) Login(dto LoginDTO) (*TokenPair, error) {
	u, err := s.users.FindByEmail(dto.Email)
	if err != nil || !CheckPassword(dto.Password, u.PasswordHash) {
		// Mensagem genérica — não revela se o email existe ou não
		return nil, fmt.Errorf("invalid credentials: %w", sharederrors.ErrUnauthorized)
	}

	tokens, err := GenerateTokenPair(u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}
	return tokens, nil
}

// Refresh valida um refresh token e emite um novo par de tokens.
func (s *Service) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", sharederrors.ErrUnauthorized)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, sharederrors.ErrUnauthorized
	}

	// Confirma que o utilizador ainda existe e não foi desativado
	u, err := s.users.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", sharederrors.ErrUnauthorized)
	}

	return GenerateTokenPair(u.ID, u.Role)
}
