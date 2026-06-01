package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/organization"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

type Service struct {
	users user.Repository
	orgs  organization.Repository
}

func NewService(users user.Repository, orgs organization.Repository) *Service {
	return &Service{users: users, orgs: orgs}
}

type RegisterDTO struct {
	Name     string    `json:"name"     validate:"required,min=2,max=100"`
	Email    string    `json:"email"    validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=72"`
	Role     user.Role `json:"role"     validate:"required,oneof=admin manager seller"`
	OrgName  string    `json:"org_name" validate:"required,min=2,max=100"`
}

type LoginDTO struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register cria uma nova organização e o seu primeiro utilizador (admin por defeito).
// Cada registo cria um tenant isolado — utilizadores adicionais são convidados
// pelo admin através de um endpoint separado (futuro).
func (s *Service) Register(dto RegisterDTO) (*user.User, error) {
	existing, err := s.users.FindByEmail(dto.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already registered: %w", sharederrors.ErrConflict)
	}

	org, err := s.orgs.Save(&organization.Organization{Name: dto.OrgName})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	hash, err := HashPassword(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &user.User{
		ID:             uuid.New(),
		Name:           dto.Name,
		Email:          dto.Email,
		PasswordHash:   hash,
		Role:           dto.Role,
		OrganizationID: org.ID,
	}

	return s.users.Save(u)
}

func (s *Service) Login(dto LoginDTO) (*TokenPair, error) {
	u, err := s.users.FindByEmail(dto.Email)
	if err != nil || !CheckPassword(dto.Password, u.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials: %w", sharederrors.ErrUnauthorized)
	}

	tokens, err := GenerateTokenPair(u.ID, u.OrganizationID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}
	return tokens, nil
}

func (s *Service) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", sharederrors.ErrUnauthorized)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, sharederrors.ErrUnauthorized
	}

	u, err := s.users.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", sharederrors.ErrUnauthorized)
	}

	return GenerateTokenPair(u.ID, u.OrganizationID, u.Role)
}
