package task

import (
	"time"

	"github.com/google/uuid"
)

// Priority define a urgência de uma tarefa.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Status define o estado de uma tarefa.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

func (s Status) IsFinal() bool {
	return s == StatusDone || s == StatusCancelled
}

// Task representa uma tarefa associada a um contacto ou negócio.
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	AssignedTo  uuid.UUID  `json:"assigned_to"`
	ContactID   *uuid.UUID `json:"contact_id,omitempty"` // pointer — task pode não ter contacto
	DealID      *uuid.UUID `json:"deal_id,omitempty"`    // pointer — task pode não ter deal
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// IsOverdue verifica se a tarefa está em atraso.
// Método com receiver de valor — não modifica a struct.
func (t Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status.IsFinal() {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// Repository define o contrato de acesso a dados para Task.
type Repository interface {
	FindByID(id uuid.UUID) (*Task, error)
	FindAll(assignedTo uuid.UUID, filters Filters) ([]*Task, int64, error)
	FindByContact(contactID uuid.UUID) ([]*Task, error)
	FindByDeal(dealID uuid.UUID) ([]*Task, error)
	FindOverdue() ([]*Task, error)
	Save(task *Task) (*Task, error)
	Update(task *Task) (*Task, error)
	Delete(id uuid.UUID) error
}

// Filters encapsula os parâmetros de pesquisa para tasks.
type Filters struct {
	Status   Status
	Priority Priority
	Page     int
	Limit    int
}

func (f *Filters) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
}
