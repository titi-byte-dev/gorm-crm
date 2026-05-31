package contact

import (
	"fmt"
	"strings"

	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
)

// Rule é a interface Chain of Responsibility para validação de contactos.
//
// Chain of Responsibility pattern: cada regra decide se processa o pedido
// ou o passa para a próxima. Aqui, uma regra que falha interrompe a cadeia.
// Novas regras adicionam-se sem tocar no Service.
type Rule interface {
	Validate(repo Reader, dto CreateContactDTO) error
}

// Reader é o sub-conjunto de Repository necessário para validação.
// Interface segregation: as regras não precisam de Save/Update/Delete.
type Reader interface {
	FindByEmail(email string) (*Contact, error)
}

// Chain executa uma sequência de regras em ordem.
// A primeira regra que falhar interrompe a validação.
type Chain []Rule

func (c Chain) Validate(repo Reader, dto CreateContactDTO) error {
	for _, rule := range c {
		if err := rule.Validate(repo, dto); err != nil {
			return err
		}
	}
	return nil
}

// DefaultChain devolve as regras de validação padrão para contactos.
func DefaultChain() Chain {
	return Chain{
		UniqueEmailRule{},
		EmailDomainRule{},
	}
}

// ---

// UniqueEmailRule garante que o email não está registado por outro contacto.
type UniqueEmailRule struct{}

func (UniqueEmailRule) Validate(repo Reader, dto CreateContactDTO) error {
	existing, err := repo.FindByEmail(dto.Email)
	if err == nil && existing != nil {
		return fmt.Errorf("email already exists: %w", sharederrors.ErrConflict)
	}
	return nil
}

// EmailDomainRule rejeita domínios descartáveis conhecidos.
// Demonstra que regras podem ser adicionadas sem tocar no Service.
type EmailDomainRule struct{}

var blockedDomains = []string{"mailinator.com", "guerrillamail.com", "temp-mail.org"}

func (EmailDomainRule) Validate(_ Reader, dto CreateContactDTO) error {
	parts := strings.SplitN(dto.Email, "@", 2)
	if len(parts) != 2 {
		return nil // validação de formato fica no validator struct tag
	}
	domain := strings.ToLower(parts[1])
	for _, blocked := range blockedDomains {
		if domain == blocked {
			return fmt.Errorf("disposable email domain not allowed: %w", sharederrors.ErrValidation)
		}
	}
	return nil
}
