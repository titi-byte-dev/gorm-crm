// Package valueobject define primitivos de dominio com comportamento proprio.
// Object Calisthenics — Regra 3: Envolve todos os primitivos e strings.
//
// Um float64 nao sabe se e negativo, nao se formata como moeda,
// nao impede que 999999999.99 seja passado como preco de uma caneta.
// Money encapsula o valor e garante as invariantes do dominio.
package valueobject

import (
	"errors"
	"fmt"
)

// Money representa um valor monetario nao-negativo.
// Tipo base float64 — GORM e JSON funcionam sem configuracao adicional.
type Money float64

var ErrNegativeMoney = errors.New("money cannot be negative")

// ParseMoney cria um Money validado. Retorna erro se negativo.
func ParseMoney(v float64) (Money, error) {
	if v < 0 {
		return 0, fmt.Errorf("%.2f: %w", v, ErrNegativeMoney)
	}
	return Money(v), nil
}

func (m Money) Float64() float64 { return float64(m) }

func (m Money) String() string { return fmt.Sprintf("%.2f", float64(m)) }

func (m Money) IsZero() bool { return m == 0 }

func (m Money) Add(other Money) Money { return m + other }

func (m Money) GreaterThan(other Money) bool { return m > other }
