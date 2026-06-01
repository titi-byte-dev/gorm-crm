package task

import (
	"fmt"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

// Service contém APENAS lógica de negócio para Tasks.
// Não importa fiber, não faz queries SQL — só regras do domínio.
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

func (s *Service) Create(dto CreateTaskDTO) (*Task, error) {
	task := &Task{
		Title:       dto.Title,
		Description: dto.Description,
		Priority:    dto.Priority,
		Status:      StatusTodo,
		AssignedTo:  dto.AssignedTo,
		ContactID:   dto.ContactID,
		DealID:      dto.DealID,
	}
	saved, err := s.repo.Save(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return saved, nil
}

func (s *Service) GetByID(id, requesterID uuid.UUID) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task.AssignedTo != requesterID {
		return nil, fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
	}
	return task, nil
}

func (s *Service) List(assignedTo uuid.UUID, filters Filters) ([]*Task, int64, error) {
	return s.repo.FindAll(assignedTo, filters)
}

// UpdateStatus aplica a regra de negócio: tarefas finais não podem ser reabertas.
// Esta regra vive no Service — o Handler não sabe nada sobre estados finais.
func (s *Service) UpdateStatus(id, requesterID uuid.UUID, newStatus Status) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update task status: %w", err)
	}
	if task.AssignedTo != requesterID {
		return nil, fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
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
			UserID:  requesterID.String(),
		})
	}
	return updated, nil
}

func (s *Service) Update(id, requesterID uuid.UUID, dto UpdateTaskDTO) (*Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	if task.AssignedTo != requesterID {
		return nil, fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
	}
	if task.Status.IsFinal() {
		return nil, fmt.Errorf("task is %s, cannot update: %w",
			task.Status, sharederrors.ErrValidation)
	}
	applyUpdates(task, dto)
	return s.repo.Update(task)
}

// applyUpdates aplica os campos opcionais do DTO à task.
// Ponteiro nil significa "não alterar" — só actualiza campos enviados.
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

func (s *Service) Delete(id, requesterID uuid.UUID) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if task.AssignedTo != requesterID {
		return fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
	}
	return s.repo.Delete(id)
}

func (s *Service) GetOverdue() ([]*Task, error) {
	return s.repo.FindOverdue()
}
