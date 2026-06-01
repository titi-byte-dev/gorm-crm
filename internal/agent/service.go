package agent

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	"github.com/titi-byte-dev/gorm-crm/internal/deal"
	"github.com/titi-byte-dev/gorm-crm/internal/lead"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
)

// Service orquestra a execução de um agente sobre uma entidade CRM.
type Service struct {
	repo     Repository
	contacts contact.Repository
	leads    lead.Repository
	deals    deal.Repository
	tasks    task.Repository
	executor *Executor
	llm      *LLMClient // nil se ANTHROPIC_API_KEY não estiver configurada
	bus      *events.Bus
}

func NewService(
	repo Repository,
	contacts contact.Repository,
	leads lead.Repository,
	deals deal.Repository,
	tasks task.Repository,
	taskSvc *task.Service,
	bus *events.Bus,
) *Service {
	return &Service{
		repo:     repo,
		contacts: contacts,
		leads:    leads,
		deals:    deals,
		tasks:    tasks,
		executor: newExecutor(taskSvc),
		llm:      NewLLMClient(),
		bus:      bus,
	}
}

// Run executa ou simula um agente, dependendo do modo e da disponibilidade do LLM.
func (s *Service) Run(rctx ctxutil.RequestCtx, dto RunDTO) (*AgentRun, error) {
	// ModeAuto requer manager ou admin
	if dto.Mode == ModeAuto && !rctx.IsManager() {
		return nil, fmt.Errorf("auto mode requires manager role: %w", sharederrors.ErrForbidden)
	}

	entityID, err := uuid.Parse(dto.EntityID)
	if err != nil {
		return nil, fmt.Errorf("invalid entity id: %w", sharederrors.ErrValidation)
	}

	run := &AgentRun{
		AgentType:  dto.AgentType,
		EntityType: dto.EntityType,
		EntityID:   entityID,
		TenantID:   rctx.TenantID,
		RunnerID:   rctx.UserID,
		Mode:       dto.Mode,
		Status:     RunStatusRunning,
	}

	saved, err := s.repo.Save(run)
	if err != nil {
		return nil, fmt.Errorf("save agent run: %w", err)
	}

	ctx, err := s.loadEntityContext(rctx, dto.EntityType, entityID)
	if err != nil {
		return s.failRun(saved, err)
	}

	if s.llm == nil {
		// sem API key: modo regra simples baseado em heurísticas
		return s.runRuleBased(saved, rctx, ctx)
	}

	prompt := BuildPrompt(dto.AgentType, *ctx)
	llmResult, err := s.llm.Run(prompt)
	if err != nil {
		return s.failRun(saved, fmt.Errorf("llm error: %w", err))
	}

	saved.Summary = llmResult.Summary
	saved.TokensUsed = llmResult.TokensUsed

	for _, call := range llmResult.ToolCalls {
		action := AgentAction{
			Tool:   call.Name,
			Input:  call.Input,
			Status: ActionPendingApproval,
		}

		if dto.Mode == ModeAuto {
			output, execErr := s.executor.Execute(call, dto.EntityType, entityID, rctx)
			if execErr != nil {
				action.Status = ActionFailed
				action.Error = execErr.Error()
			} else {
				action.Status = ActionExecuted
				action.Output = output
			}
		}

		saved.Actions = append(saved.Actions, action)
	}

	now := time.Now()
	saved.Status = RunStatusCompleted
	saved.CompletedAt = &now

	updated, err := s.repo.Update(saved)
	if err != nil {
		return nil, fmt.Errorf("update agent run: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.AgentRunCompleted,
		Payload: map[string]string{"run_id": saved.ID.String(), "entity_id": entityID.String()},
		UserID:  rctx.UserID.String(),
	})

	return updated, nil
}

// ApproveActions executa as ações pendentes de um run em modo suggest.
func (s *Service) ApproveActions(runID uuid.UUID, rctx ctxutil.RequestCtx, indices []int) (*AgentRun, error) {
	run, err := s.repo.FindByID(runID)
	if err != nil {
		return nil, err
	}
	if run.TenantID != rctx.TenantID {
		return nil, fmt.Errorf("run %s: %w", runID, sharederrors.ErrNotFound)
	}

	for _, i := range indices {
		if i < 0 || i >= len(run.Actions) {
			continue
		}
		action := &run.Actions[i]
		if action.Status != ActionPendingApproval {
			continue
		}
		call := ToolCall{Name: action.Tool, Input: action.Input}
		output, execErr := s.executor.Execute(call, run.EntityType, run.EntityID, rctx)
		if execErr != nil {
			action.Status = ActionFailed
			action.Error = execErr.Error()
		} else {
			action.Status = ActionExecuted
			action.Output = output
		}
	}

	return s.repo.Update(run)
}

// GetRunsByEntity devolve o histórico de runs para uma entidade.
func (s *Service) GetRunsByEntity(rctx ctxutil.RequestCtx, entityType string, entityID uuid.UUID) ([]*AgentRun, error) {
	return s.repo.FindByEntity(rctx.TenantID, entityType, entityID)
}

// loadEntityContext carrega os dados da entidade e tasks associadas.
func (s *Service) loadEntityContext(rctx ctxutil.RequestCtx, entityType string, entityID uuid.UUID) (*EntityContext, error) {
	ctx := &EntityContext{}
	var err error

	switch entityType {
	case "contact":
		ctx.Contact, err = s.contacts.FindByID(entityID)
		if err != nil {
			return nil, err
		}
		ctx.Tasks, _ = s.tasks.FindByContact(entityID)
	case "deal":
		ctx.Deal, err = s.deals.FindByID(entityID)
		if err != nil {
			return nil, err
		}
		ctx.Tasks, _ = s.tasks.FindByDeal(entityID)
	case "lead":
		ctx.Lead, err = s.leads.FindByID(entityID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown entity type: %s", entityType)
	}

	return ctx, nil
}

// runRuleBased aplica heurísticas simples quando não há LLM disponível.
func (s *Service) runRuleBased(run *AgentRun, rctx ctxutil.RequestCtx, ctx *EntityContext) (*AgentRun, error) {
	var actions []AgentAction

	// regra: contacto sem atualização há mais de 7 dias → follow-up task
	if ctx.Contact != nil {
		daysSince := int(time.Since(ctx.Contact.UpdatedAt).Hours() / 24)
		if daysSince >= 7 {
			actions = append(actions, AgentAction{
				Tool: "create_task",
				Input: map[string]any{
					"title":    fmt.Sprintf("Follow-up: %s", ctx.Contact.Name),
					"priority": "high",
					"due_days": float64(2),
				},
				Status: ActionPendingApproval,
			})
		}
	}

	// regra: tasks em atraso → alertar
	for _, t := range ctx.Tasks {
		if t.IsOverdue() {
			actions = append(actions, AgentAction{
				Tool: "escalate_to_manager",
				Input: map[string]any{
					"reason": fmt.Sprintf("Task '%s' está em atraso", t.Title),
				},
				Status: ActionPendingApproval,
			})
			break // um alerta chega
		}
	}

	if len(actions) == 0 {
		actions = append(actions, AgentAction{
			Tool:   "summarize_only",
			Input:  map[string]any{"summary": "Nenhuma ação necessária neste momento."},
			Status: ActionExecuted,
			Output: map[string]any{"summary": "Nenhuma ação necessária neste momento."},
		})
	}

	if run.Mode == ModeAuto {
		for i, a := range actions {
			if a.Status == ActionPendingApproval {
				call := ToolCall{Name: a.Tool, Input: a.Input}
				output, execErr := s.executor.Execute(call, run.EntityType, run.EntityID, rctx)
				if execErr != nil {
					actions[i].Status = ActionFailed
					actions[i].Error = execErr.Error()
				} else {
					actions[i].Status = ActionExecuted
					actions[i].Output = output
				}
			}
		}
	}

	now := time.Now()
	run.Actions = actions
	run.Summary = "Análise baseada em regras (LLM não configurado)."
	run.Status = RunStatusCompleted
	run.CompletedAt = &now

	return s.repo.Update(run)
}

func (s *Service) failRun(run *AgentRun, cause error) (*AgentRun, error) {
	run.Status = RunStatusFailed
	run.Summary = cause.Error()
	_, _ = s.repo.Update(run)
	return nil, cause
}
