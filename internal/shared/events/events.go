package events

import (
	"context"
	"log/slog"
)

// EventType identifica o tipo de evento emitido no sistema.
type EventType string

const (
	ContactCreated EventType = "contact.created"
	ContactUpdated EventType = "contact.updated"
	ContactDeleted EventType = "contact.deleted"
	LeadCreated    EventType = "lead.created"
	LeadConverted  EventType = "lead.converted"
	LeadLost       EventType = "lead.lost"
	DealWon        EventType = "deal.won"
	DealLost       EventType = "deal.lost"
	TaskOverdue    EventType = "task.overdue"
)

// Event é a mensagem que circula no bus.
type Event struct {
	Type    EventType
	Payload any // conteúdo específico de cada evento
	UserID  string
}

// Handler é a assinatura de qualquer função que processa eventos.
// Em Go, funções são valores de primeira classe — podem ser passadas como argumentos.
type Handler func(ctx context.Context, event Event)

// Bus é o event bus em memória baseado em channels Go.
//
// Como funciona:
//   - Publish envia um evento para o channel (não bloqueia — channel tem buffer)
//   - Workers (goroutines) escutam o channel e processam cada evento
//   - Cada tipo de evento pode ter vários handlers registados
type Bus struct {
	ch       chan Event
	handlers map[EventType][]Handler
	logger   *slog.Logger
}

// DefaultBufferSize é a capacidade do channel do bus em produção.
// Suficiente para absorver picos sem bloquear os handlers HTTP.
const DefaultBufferSize = 500

// New cria um Bus com um channel com buffer de capacidade cap.
// Um channel com buffer não bloqueia o publisher enquanto houver espaço.
func New(cap int, logger *slog.Logger) *Bus {
	return &Bus{
		ch:       make(chan Event, cap),
		handlers: make(map[EventType][]Handler),
		logger:   logger,
	}
}

// Subscribe regista um handler para um tipo de evento.
func (b *Bus) Subscribe(eventType EventType, handler Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish envia um evento para o channel.
// select com default evita bloqueio se o channel estiver cheio — o evento é descartado
// e registado no log. Em produção, usaríamos uma fila persistente (Módulo 17).
func (b *Bus) Publish(event Event) {
	select {
	case b.ch <- event:
	default:
		b.logger.Warn("event bus full, dropping event", "type", event.Type)
	}
}

// Start inicia o worker loop numa goroutine separada.
// A goroutine corre em background até o contexto ser cancelado.
//
// Padrão Go: lançar goroutines com `go` e coordenar com context para shutdown limpo.
func (b *Bus) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case event := <-b.ch:
				b.dispatch(ctx, event)
			case <-ctx.Done():
				b.logger.Info("event bus shutting down")
				return
			}
		}
	}()
}

func (b *Bus) dispatch(ctx context.Context, event Event) {
	handlers, ok := b.handlers[event.Type]
	if !ok {
		return
	}
	for _, h := range handlers {
		h(ctx, event)
	}
}
