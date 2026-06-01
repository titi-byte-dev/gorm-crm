package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	"github.com/titi-byte-dev/gorm-crm/internal/deal"
	"github.com/titi-byte-dev/gorm-crm/internal/lead"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
)

// EntityContext agrega os dados de uma entidade para construção do prompt.
type EntityContext struct {
	Contact *contact.Contact
	Lead    *lead.Lead
	Deal    *deal.Deal
	Tasks   []*task.Task
}

// BuildPrompt constrói o prompt para o LLM a partir do contexto da entidade.
func BuildPrompt(agentType AgentType, ctx EntityContext) string {
	var sb strings.Builder

	switch agentType {
	case AgentFollowUp:
		sb.WriteString("=== TAREFA: Analisar contacto e propor acompanhamento ===\n\n")
		writeContactSection(&sb, ctx.Contact)
		writeTasksSection(&sb, ctx.Tasks)
	case AgentDealCloser:
		sb.WriteString("=== TAREFA: Analisar deal e propor próximos passos para fechar ===\n\n")
		writeDealSection(&sb, ctx.Deal)
		writeTasksSection(&sb, ctx.Tasks)
	case AgentTaskRouter:
		sb.WriteString("=== TAREFA: Priorizar e distribuir tasks pendentes ===\n\n")
		writeTasksSection(&sb, ctx.Tasks)
	case AgentSummarize:
		sb.WriteString("=== TAREFA: Resumir o histórico e estado atual ===\n\n")
		if ctx.Contact != nil {
			writeContactSection(&sb, ctx.Contact)
		}
		if ctx.Deal != nil {
			writeDealSection(&sb, ctx.Deal)
		}
		if ctx.Lead != nil {
			writeLeadSection(&sb, ctx.Lead)
		}
		writeTasksSection(&sb, ctx.Tasks)
	}

	sb.WriteString(fmt.Sprintf("\nData atual: %s\n", time.Now().Format("2006-01-02")))
	sb.WriteString("\nAnalisa o contexto acima e decide quais as ações mais adequadas.")

	return sb.String()
}

func writeContactSection(sb *strings.Builder, c *contact.Contact) {
	if c == nil {
		return
	}
	sb.WriteString("CONTACTO:\n")
	sb.WriteString(fmt.Sprintf("  Nome: %s\n", c.Name))
	sb.WriteString(fmt.Sprintf("  Email: %s\n", c.Email))
	if c.Phone != "" {
		sb.WriteString(fmt.Sprintf("  Telefone: %s\n", c.Phone))
	}
	if c.Company != "" {
		sb.WriteString(fmt.Sprintf("  Empresa: %s\n", c.Company))
	}
	daysSince := int(time.Since(c.UpdatedAt).Hours() / 24)
	sb.WriteString(fmt.Sprintf("  Última atualização: há %d dias\n", daysSince))
	if c.Notes != "" {
		sb.WriteString(fmt.Sprintf("  Notas: %s\n", c.Notes))
	}
	sb.WriteString("\n")
}

func writeDealSection(sb *strings.Builder, d *deal.Deal) {
	if d == nil {
		return
	}
	sb.WriteString("DEAL:\n")
	sb.WriteString(fmt.Sprintf("  Título: %s\n", d.Title))
	sb.WriteString(fmt.Sprintf("  Valor: %.2f €\n", d.Value))
	sb.WriteString(fmt.Sprintf("  Stage: %s\n", d.Stage))
	daysSince := int(time.Since(d.UpdatedAt).Hours() / 24)
	sb.WriteString(fmt.Sprintf("  Sem atualização há: %d dias\n", daysSince))
	sb.WriteString("\n")
}

func writeLeadSection(sb *strings.Builder, l *lead.Lead) {
	if l == nil {
		return
	}
	sb.WriteString("LEAD:\n")
	sb.WriteString(fmt.Sprintf("  Título: %s\n", l.Title))
	sb.WriteString(fmt.Sprintf("  Valor: %.2f €\n", l.Value))
	sb.WriteString(fmt.Sprintf("  Status: %s\n", l.Status))
	daysSince := int(time.Since(l.UpdatedAt).Hours() / 24)
	sb.WriteString(fmt.Sprintf("  Sem atualização há: %d dias\n", daysSince))
	sb.WriteString("\n")
}

func writeTasksSection(sb *strings.Builder, tasks []*task.Task) {
	if len(tasks) == 0 {
		sb.WriteString("TASKS: nenhuma task associada\n\n")
		return
	}
	sb.WriteString(fmt.Sprintf("TASKS (%d):\n", len(tasks)))
	for _, t := range tasks {
		due := "sem prazo"
		overdue := ""
		if t.DueDate != nil {
			due = t.DueDate.Format("2006-01-02")
			if t.IsOverdue() {
				overdue = " ⚠ EM ATRASO"
			}
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s | prioridade: %s | prazo: %s%s\n",
			t.Status, t.Title, t.Priority, due, overdue))
	}
	sb.WriteString("\n")
}
