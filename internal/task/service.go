package task

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

type Service struct {
	repo Repository
	bus  *events.Bus
}

func NewService(repo Repository, bus *events.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

type CreateTaskDTO struct {
	Title       string     `json:"title"       validate:"required,min=2,max=200"`
	Description string     `json:"description" validate:"omitempty,max=2000"`
	Priority    Priority   `json:"priority"    validate:"required,oneof=low medium high urgent"`
	AssignedTo  uuid.UUID  `json:"assigned_to" validate:"required"`
	ContactID   *uuid.UUID `json:"contact_id"`
	DealID      *uuid.UUID `json:"deal_id"`
	DueDate     *string    `json:"due_date"    validate:"omitempty"`
}

type UpdateTaskDTO struct {
	Title       *string   `json:"title"       validate:"omitempty,min=2,max=200"`
	Description *string   `json:"description" validate:"omitempty,max=2000"`
	Priority    *Priority `json:"priority"    validate:"omitempty,oneof=low medium high urgent"`
	Status      *Status   `json:"status"      validate:"omitempty,oneof=todo in_progress done cancelled"`
}

func (s *Service) Create(rctx ctxutil.RequestCtx, dto CreateTaskDTO) (*Task, error) {
	task := &Task{
		Title:       dto.Title,
		Description: dto.Description,
		Priority:    dto.Priority,
		Status:      StatusTodo,
		AssignedTo:  dto.AssignedTo,
		TenantID:    rctx.TenantID,
		ContactID:   dto.ContactID,
		DealID:      dto.DealID,
	}
	saved, err := s.repo.Save(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID, rctx ctxutil.RequestCtx) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.checkAccess(task, rctx); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) List(rctx ctxutil.RequestCtx, filters Filters) ([]*Task, int64, error) {
	return s.repo.FindAll(rctx.TenantID, rctx.UserID, rctx.IsManager(), filters)
}

func (s *Service) UpdateStatus(id uuid.UUID, rctx ctxutil.RequestCtx, newStatus Status) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update task status: %w", err)
	}
	if err := s.checkAccess(task, rctx); err != nil {
		return nil, err
	}
	if task.Status.IsFinal() {
		return nil, fmt.Errorf("task %s is %s and cannot be updated: %w",
			id, task.Status, sharederrors.ErrValidation)
	}
	task.Status = newStatus
	updated, err := s.repo.Update(task)
	if err != nil {
		return nil, fmt.Errorf("save task status: %w", err)
	}
	if newStatus == StatusDone {
		s.bus.Publish(events.Event{
			Type:    events.TaskCompleted,
			Payload: map[string]string{"id": id.String()},
			UserID:  rctx.UserID.String(),
		})
	}
	return updated, nil
}

func (s *Service) Update(id uuid.UUID, rctx ctxutil.RequestCtx, dto UpdateTaskDTO) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	if err := s.checkAccess(task, rctx); err != nil {
		return nil, err
	}
	if task.Status.IsFinal() {
		return nil, fmt.Errorf("task is %s, cannot update: %w",
			task.Status, sharederrors.ErrValidation)
	}
	applyUpdates(task, dto)
	return s.repo.Update(task)
}

func (s *Service) Delete(id uuid.UUID, rctx ctxutil.RequestCtx) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if err := s.checkAccess(task, rctx); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) GetOverdue() ([]*Task, error) {
	return s.repo.FindOverdue()
}

// checkAccess: manager vê toda a org, seller só as suas tarefas atribuídas.
func (s *Service) checkAccess(t *Task, rctx ctxutil.RequestCtx) error {
	if t.TenantID != rctx.TenantID {
		return fmt.Errorf("task %s: %w", t.ID, sharederrors.ErrNotFound)
	}
	if !rctx.IsManager() && t.AssignedTo != rctx.UserID {
		return fmt.Errorf("task %s: %w", t.ID, sharederrors.ErrNotFound)
	}
	return nil
}

func applyUpdates(t *Task, dto UpdateTaskDTO) {
	if dto.Title != nil {
		t.Title = *dto.Title
	}
	if dto.Description != nil {
		t.Description = *dto.Description
	}
	if dto.Priority != nil {
		t.Priority = *dto.Priority
	}
	if dto.Status != nil {
		t.Status = *dto.Status
	}
}
