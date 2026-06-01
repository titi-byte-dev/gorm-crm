package agent

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
)

// Executor converte ToolCalls do LLM em ações reais no CRM.
// Usa os services existentes para que todas as regras de negócio sejam respeitadas.
type Executor struct {
	tasks *task.Service
}

func newExecutor(tasks *task.Service) *Executor {
	return &Executor{tasks: tasks}
}

// Execute corre uma ToolCall contra a entidade alvo e devolve o resultado.
func (e *Executor) Execute(
	call ToolCall,
	entityType string,
	entityID uuid.UUID,
	rctx ctxutil.RequestCtx,
) (map[string]any, error) {
	switch call.Name {
	case "create_task":
		return e.createTask(call.Input, entityType, entityID, rctx)
	case "add_note":
		return e.addNote(call.Input)
	case "summarize_only":
		return map[string]any{"summary": call.Input["summary"]}, nil
	case "escalate_to_manager":
		return map[string]any{"escalated": true, "reason": call.Input["reason"]}, nil
	case "update_lead_status":
		// status update é gerido pelo lead.Service — devolvemos a intenção sem executar aqui
		// O service principal decide com base no modo (suggest vs auto)
		return map[string]any{"intended_status": call.Input["new_status"]}, nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (e *Executor) createTask(
	input map[string]any,
	entityType string,
	entityID uuid.UUID,
	rctx ctxutil.RequestCtx,
) (map[string]any, error) {
	title, _ := input["title"].(string)
	priorityStr, _ := input["priority"].(string)
	if priorityStr == "" {
		priorityStr = "medium"
	}

	dto := task.CreateTaskDTO{
		Title:      title,
		Priority:   task.Priority(priorityStr),
		AssignedTo: rctx.UserID,
	}

	// associar à entidade correta
	switch entityType {
	case "contact":
		dto.ContactID = &entityID
	case "deal":
		dto.DealID = &entityID
	}

	// calcular due_date se due_days foi fornecido
	if dueDays, ok := input["due_days"].(float64); ok && dueDays > 0 {
		due := time.Now().AddDate(0, 0, int(dueDays))
		dueStr := due.Format("2006-01-02")
		dto.DueDate = &dueStr
	}

	created, err := e.tasks.Create(rctx, dto)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return map[string]any{"task_id": created.ID.String(), "title": created.Title}, nil
}

func (e *Executor) addNote(input map[string]any) (map[string]any, error) {
	content, _ := input["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("note content is required")
	}
	// Nota: persist no activity log (MongoDB) ficará no módulo de integrações.
	// Por agora devolvemos confirmação — o handler pode publicar no events.Bus.
	return map[string]any{"noted": true, "preview": content[:min(len(content), 80)]}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
