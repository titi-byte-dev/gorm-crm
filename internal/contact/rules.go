package contact

import sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"

// Rule e uma regra de negocio aplicada antes de criar um contacto.
// OCP — Open/Closed Principle: o Service esta fechado para modificacao,
// mas aberto para extensao atraves de novas Rules.
//
// Para adicionar uma nova regra (ex: telefone unico), basta:
//   1. Criar um novo tipo que implemente Rule
//   2. Passar ao NewService — sem tocar no Service.Create
type Rule interface {
	Validate(repo Reader, dto CreateContactDTO) error
}

// UniqueEmailRule verifica que nao existe outro contacto com o mesmo email.
// E a regra de negocio que antes estava inline em Service.Create.
type UniqueEmailRule struct{}

func (r UniqueEmailRule) Validate(repo Reader, dto CreateContactDTO) error {
	existing, err := repo.FindByEmail(dto.Email)
	if err == nil && existing != nil {
		return sharederrors.ErrConflict
	}
	return nil
}
