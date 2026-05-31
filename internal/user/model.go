package user

import (
	"time"

	"github.com/google/uuid"
)

// Role define o nível de acesso de um utilizador no sistema.
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleSeller  Role = "seller"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleSeller:
		return true
	}
	return false
}

// User representa um utilizador do GoRM CRM.
// Em Go, structs são a unidade central de dados — não há classes.
type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" exclui do JSON — nunca exposto na API
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Repository define o contrato de acesso a dados para User.
// Em Go, interfaces são implícitas — qualquer tipo que implemente
// estes métodos satisfaz esta interface, sem declaração explícita.
type Repository interface {
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) (*User, error)
	Update(user *User) (*User, error)
}
