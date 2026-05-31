package auth

import "golang.org/x/crypto/bcrypt"

// bcryptCost define o custo computacional do hash.
// Cada incremento duplica o tempo de processamento.
// Cost 12 leva ~250ms numa máquina moderna — suficiente para tornar
// brute force impraticável sem impactar a experiência do utilizador.
const bcryptCost = 12

// HashPassword gera um hash bcrypt da password.
// O hash inclui automaticamente um salt aleatório — dois hashes da
// mesma password são sempre diferentes, o que impede rainbow table attacks.
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	return string(bytes), err
}

// CheckPassword compara uma password em plain text com um hash bcrypt.
// Usa constant-time comparison internamente — imune a timing attacks.
func CheckPassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
