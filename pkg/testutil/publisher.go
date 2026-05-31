// Package testutil fornece implementacoes de teste para interfaces do sistema.
// Usa-se em testes unitarios em vez das implementacoes reais.
package testutil

import "github.com/titi-byte-dev/gorm-crm/internal/shared/events"

// NullPublisher implementa events.Publisher sem fazer nada.
// LSP — Liskov Substitution Principle: qualquer implementacao de Publisher
// deve ser substituivel onde Publisher e esperado, sem quebrar o chamador.
//
// NullPublisher e o caso minimo: satisfaz o contrato (nao panics, nao erros)
// mas descarta os eventos — util em testes que nao verificam eventos.
type NullPublisher struct{}

// Verify compile-time que NullPublisher satisfaz events.Publisher.
// Se a interface mudar e NullPublisher nao for actualizado, o compilador avisa aqui.
var _ events.Publisher = (*NullPublisher)(nil)

func (NullPublisher) Publish(_ events.Event) {}

// SpyPublisher implementa events.Publisher e regista todos os eventos publicados.
// Util em testes que verificam SE e QUAIS eventos foram emitidos.
type SpyPublisher struct {
	Events []events.Event
}

var _ events.Publisher = (*SpyPublisher)(nil)

func (s *SpyPublisher) Publish(event events.Event) {
	s.Events = append(s.Events, event)
}

// Published devolve true se pelo menos um evento do tipo dado foi publicado.
func (s *SpyPublisher) Published(t events.EventType) bool {
	for _, e := range s.Events {
		if e.Type == t {
			return true
		}
	}
	return false
}
