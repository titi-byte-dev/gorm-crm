package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

// taskRecord é o modelo GORM — separado do domain model Task.
// Esta separação é intencional: o domain model não conhece GORM,
// e o GORM record não é exposto fora deste ficheiro.
type taskRecord struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Title       string     `gorm:"not null"`
	Description string
	Priority    string     `gorm:"not null;default:'medium'"`
	Status      string     `gorm:"not null;default:'todo'"`
	AssignedTo  uuid.UUID  `gorm:"type:uuid;not null;index"`
	TenantID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ContactID   *uuid.UUID `gorm:"type:uuid;index"`
	DealID      *uuid.UUID `gorm:"type:uuid;index"`
	DueDate     *time.Time
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

func (taskRecord) TableName() string { return "tasks" }

// Verificação de interface em compile-time.
// Se postgresRepository não implementar Repository, o build falha aqui.
var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*Task, error) {
	var rec taskRecord
	if err := r.db.First(&rec, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find task: %w", err)
	}
	return recordToTask(rec), nil
}

func (r *postgresRepository) FindAll(tenantID, assignedTo uuid.UUID, isManager bool, filters Filters) ([]*Task, int64, error) {
	filters.SetDefaults()
	query := r.db.Model(&taskRecord{}).Where("tenant_id = ?", tenantID)
	if !isManager {
		query = query.Where("assigned_to = ?", assignedTo)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
	}
	var total int64
	query.Count(&total)
	var recs []taskRecord
	if err := query.Limit(filters.Limit).Offset((filters.Page-1)*filters.Limit).Find(&recs).Error; err != nil {
		return nil, 0, fmt.Errorf("list tasks: %w", err)
	}
	tasks := make([]*Task, len(recs))
	for i, rec := range recs {
		tasks[i] = recordToTask(rec)
	}
	return tasks, total, nil
}

func (r *postgresRepository) FindByContact(contactID uuid.UUID) ([]*Task, error) {
	var recs []taskRecord
	if err := r.db.Where("contact_id = ?", contactID).Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("find tasks by contact: %w", err)
	}
	return recsToTasks(recs), nil
}

func (r *postgresRepository) FindByDeal(dealID uuid.UUID) ([]*Task, error) {
	var recs []taskRecord
	if err := r.db.Where("deal_id = ?", dealID).Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("find tasks by deal: %w", err)
	}
	return recsToTasks(recs), nil
}

// FindOverdue usa o partial index criado na migration 005.
// A query WHERE due_date < NOW() AND status NOT IN (...)
// é eficiente porque o index só contém tasks ativas.
func (r *postgresRepository) FindOverdue() ([]*Task, error) {
	var recs []taskRecord
	err := r.db.Where("due_date < ? AND status NOT IN ('done','cancelled')", time.Now()).Find(&recs).Error
	if err != nil {
		return nil, fmt.Errorf("find overdue tasks: %w", err)
	}
	return recsToTasks(recs), nil
}

func (r *postgresRepository) Save(t *Task) (*Task, error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	rec := taskToRecord(t)
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save task: %w", err)
	}
	return recordToTask(rec), nil
}

func (r *postgresRepository) Update(t *Task) (*Task, error) {
	rec := taskToRecord(t)
	result := r.db.Model(&rec).Where("id = ?", t.ID).Updates(rec)
	if result.Error != nil {
		return nil, fmt.Errorf("update task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("task %s: %w", t.ID, sharederrors.ErrNotFound)
	}
	return t, nil
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&taskRecord{})
	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
	}
	return nil
}

func recordToTask(r taskRecord) *Task {
	return &Task{
		ID: r.ID, Title: r.Title, Description: r.Description,
		Priority: Priority(r.Priority), Status: Status(r.Status),
		AssignedTo: r.AssignedTo, TenantID: r.TenantID,
		ContactID: r.ContactID, DealID: r.DealID,
		DueDate: r.DueDate, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func taskToRecord(t *Task) taskRecord {
	return taskRecord{
		ID: t.ID, Title: t.Title, Description: t.Description,
		Priority: string(t.Priority), Status: string(t.Status),
		AssignedTo: t.AssignedTo, TenantID: t.TenantID,
		ContactID: t.ContactID, DealID: t.DealID, DueDate: t.DueDate,
	}
}

func recsToTasks(recs []taskRecord) []*Task {
	tasks := make([]*Task, len(recs))
	for i, r := range recs {
		tasks[i] = recordToTask(r)
	}
	return tasks
}
