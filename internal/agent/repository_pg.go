package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

// agentRunRecord é o modelo GORM — mantido separado do domínio.
type agentRunRecord struct {
	ID          string     `gorm:"column:id;primaryKey"`
	AgentType   string     `gorm:"column:agent_type"`
	EntityType  string     `gorm:"column:entity_type"`
	EntityID    string     `gorm:"column:entity_id"`
	TenantID    string     `gorm:"column:tenant_id"`
	RunnerID    string     `gorm:"column:runner_id"`
	Mode        string     `gorm:"column:mode"`
	Status      string     `gorm:"column:status"`
	ActionsJSON []byte     `gorm:"column:actions"`
	Summary     string     `gorm:"column:summary"`
	TokensUsed  int        `gorm:"column:tokens_used"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

func (agentRunRecord) TableName() string { return "agent_runs" }

func (r *postgresRepository) Save(run *AgentRun) (*AgentRun, error) {
	rec, err := toRecord(run)
	if err != nil {
		return nil, fmt.Errorf("marshal agent run: %w", err)
	}
	rec.ID = uuid.New().String()
	if err := r.db.Create(rec).Error; err != nil {
		return nil, fmt.Errorf("save agent run: %w", err)
	}
	return fromRecord(rec)
}

func (r *postgresRepository) Update(run *AgentRun) (*AgentRun, error) {
	rec, err := toRecord(run)
	if err != nil {
		return nil, fmt.Errorf("marshal agent run: %w", err)
	}
	if err := r.db.Save(rec).Error; err != nil {
		return nil, fmt.Errorf("update agent run: %w", err)
	}
	return fromRecord(rec)
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*AgentRun, error) {
	var rec agentRunRecord
	if err := r.db.Where("id = ?", id.String()).First(&rec).Error; err != nil {
		return nil, fmt.Errorf("find agent run: %w", err)
	}
	return fromRecord(&rec)
}

func (r *postgresRepository) FindByEntity(tenantID uuid.UUID, entityType string, entityID uuid.UUID) ([]*AgentRun, error) {
	var recs []agentRunRecord
	if err := r.db.
		Where("tenant_id = ? AND entity_type = ? AND entity_id = ?",
			tenantID.String(), entityType, entityID.String()).
		Order("created_at DESC").
		Limit(50).
		Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("find agent runs by entity: %w", err)
	}
	runs := make([]*AgentRun, 0, len(recs))
	for i := range recs {
		run, err := fromRecord(&recs[i])
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func toRecord(run *AgentRun) (*agentRunRecord, error) {
	actionsJSON, err := json.Marshal(run.Actions)
	if err != nil {
		return nil, err
	}
	return &agentRunRecord{
		ID:          run.ID.String(),
		AgentType:   string(run.AgentType),
		EntityType:  run.EntityType,
		EntityID:    run.EntityID.String(),
		TenantID:    run.TenantID.String(),
		RunnerID:    run.RunnerID.String(),
		Mode:        string(run.Mode),
		Status:      string(run.Status),
		ActionsJSON: actionsJSON,
		Summary:     run.Summary,
		TokensUsed:  run.TokensUsed,
		CreatedAt:   run.CreatedAt,
		CompletedAt: run.CompletedAt,
	}, nil
}

func fromRecord(rec *agentRunRecord) (*AgentRun, error) {
	id, _ := uuid.Parse(rec.ID)
	entityID, _ := uuid.Parse(rec.EntityID)
	tenantID, _ := uuid.Parse(rec.TenantID)
	runnerID, _ := uuid.Parse(rec.RunnerID)

	var actions []AgentAction
	if len(rec.ActionsJSON) > 0 {
		if err := json.Unmarshal(rec.ActionsJSON, &actions); err != nil {
			return nil, fmt.Errorf("unmarshal actions: %w", err)
		}
	}

	return &AgentRun{
		ID:          id,
		AgentType:   AgentType(rec.AgentType),
		EntityType:  rec.EntityType,
		EntityID:    entityID,
		TenantID:    tenantID,
		RunnerID:    runnerID,
		Mode:        RunMode(rec.Mode),
		Status:      RunStatus(rec.Status),
		Actions:     actions,
		Summary:     rec.Summary,
		TokensUsed:  rec.TokensUsed,
		CreatedAt:   rec.CreatedAt,
		CompletedAt: rec.CompletedAt,
	}, nil
}
