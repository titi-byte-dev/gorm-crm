package agent

import (
	"time"

	"github.com/google/uuid"
)

// AgentType identifica que tipo de raciocínio o agente deve aplicar.
type AgentType string

const (
	AgentFollowUp   AgentType = "follow_up"   // analisa contacto e propõe acompanhamento
	AgentDealCloser AgentType = "deal_closer"  // analisa deal e propõe próximos passos
	AgentTaskRouter AgentType = "task_router"  // distribui/prioriza tasks da equipa
	AgentSummarize  AgentType = "summarize"    // resume histórico de uma entidade
)

// RunMode controla se o agente executa ou apenas sugere.
type RunMode string

const (
	// ModeSuggest devolve ações propostas — o utilizador aprova antes de executar.
	ModeSuggest RunMode = "suggest"
	// ModeAuto executa imediatamente (requer role manager ou admin).
	ModeAuto RunMode = "auto"
)

// RunStatus representa o ciclo de vida de uma execução.
type RunStatus string

const (
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
)

// ActionStatus representa o estado de uma ação individual.
type ActionStatus string

const (
	ActionPendingApproval ActionStatus = "pending_approval"
	ActionExecuted        ActionStatus = "executed"
	ActionSkipped         ActionStatus = "skipped"
	ActionFailed          ActionStatus = "failed"
)

// AgentAction é uma ação proposta ou executada pelo agente.
type AgentAction struct {
	Tool   string         `json:"tool"`   // nome da ferramenta (create_task, update_lead_status…)
	Input  map[string]any `json:"input"`  // parâmetros enviados ao executor
	Output map[string]any `json:"output"` // resultado após execução
	Status ActionStatus   `json:"status"`
	Error  string         `json:"error,omitempty"`
}

// AgentRun representa uma execução completa de um agente sobre uma entidade.
type AgentRun struct {
	ID          uuid.UUID    `json:"id"`
	AgentType   AgentType    `json:"agent_type"`
	EntityType  string       `json:"entity_type"` // "contact", "deal", "lead"
	EntityID    uuid.UUID    `json:"entity_id"`
	TenantID    uuid.UUID    `json:"tenant_id"`
	RunnerID    uuid.UUID    `json:"runner_id"` // utilizador que ativou
	Mode        RunMode      `json:"mode"`
	Status      RunStatus    `json:"status"`
	Actions     []AgentAction `json:"actions"`
	Summary     string       `json:"summary"` // explicação em linguagem natural
	TokensUsed  int          `json:"tokens_used"`
	CreatedAt   time.Time    `json:"created_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

// Repository define o contrato de persistência de AgentRun.
type Repository interface {
	Save(run *AgentRun) (*AgentRun, error)
	Update(run *AgentRun) (*AgentRun, error)
	FindByID(id uuid.UUID) (*AgentRun, error)
	FindByEntity(tenantID uuid.UUID, entityType string, entityID uuid.UUID) ([]*AgentRun, error)
}

// RunDTO é o payload do pedido HTTP para iniciar uma execução.
type RunDTO struct {
	AgentType  AgentType `json:"agent_type"  validate:"required,oneof=follow_up deal_closer task_router summarize"`
	EntityType string    `json:"entity_type" validate:"required,oneof=contact deal lead"`
	EntityID   string    `json:"entity_id"   validate:"required,uuid"`
	Mode       RunMode   `json:"mode"        validate:"required,oneof=suggest auto"`
}

// ApproveDTO aprova uma ação pendente de uma run em modo suggest.
type ApproveDTO struct {
	ActionIndices []int `json:"action_indices" validate:"required,min=1"`
}
