package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

// ValidationError descreve um erro de validação num campo específico.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Result é devolvido quando há erros de validação.
type Result struct {
	Errors []ValidationError `json:"errors"`
}

// Check valida uma struct anotada com tags `validate:"..."`.
// Devolve nil se válida, ou um Result com todos os erros encontrados.
func Check(s any) *Result {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var errs []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidationError{
			Field:   strings.ToLower(e.Field()),
			Message: fieldMessage(e),
		})
	}
	return &Result{Errors: errs}
}

func fieldMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "campo obrigatório"
	case "email":
		return "email inválido"
	case "min":
		return fmt.Sprintf("mínimo %s caracteres", e.Param())
	case "max":
		return fmt.Sprintf("máximo %s caracteres", e.Param())
	case "oneof":
		return fmt.Sprintf("valor deve ser um de: %s", e.Param())
	default:
		return fmt.Sprintf("falhou validação '%s'", e.Tag())
	}
}
