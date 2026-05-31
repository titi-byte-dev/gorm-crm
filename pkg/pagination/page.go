package pagination

const (
	defaultPage    = 1
	defaultLimit   = 20
	maxLimit       = 100
	defaultSortDir = "desc"
)

// Base contém os parâmetros de paginação partilhados por todas as listagens do CRM.
//
// Em Go não existe herança — usa-se composição via embedding.
// Cada Filters de domínio embebe Base e herda os seus métodos directamente.
//
// Diferença chave vs herança:
//   - Herança (OOP clássico): "ContactFilters É UM Filters"
//   - Embedding (Go):         "ContactFilters TEM UM Base" — promoção de métodos
type Base struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

// Normalize aplica valores por omissão.
// Recebe o campo de ordenação por omissão porque cada domínio tem o seu.
func (b *Base) Normalize(defaultSort string) {
	if b.Page <= 0 {
		b.Page = defaultPage
	}
	if b.Limit <= 0 || b.Limit > maxLimit {
		b.Limit = defaultLimit
	}
	if b.SortBy == "" {
		b.SortBy = defaultSort
	}
	if b.SortDir == "" {
		b.SortDir = defaultSortDir
	}
}

// Offset calcula quantos registos saltar para a página actual.
// (Page - 1) * Limit — página 1 começa em 0, página 2 em Limit, etc.
func (b Base) Offset() int {
	return (b.Page - 1) * b.Limit
}
