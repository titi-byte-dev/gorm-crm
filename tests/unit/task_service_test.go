package unit_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

var _ task.Repository = (*mockTaskRepository)(nil)

type mockTaskRepository struct {
	tasks map[uuid.UUID]*task.Task
}

func newMockRepo() *mockTaskRepository {
	return &mockTaskRepository{tasks: make(map[uuid.UUID]*task.Task)}
}

func (m *mockTaskRepository) FindByID(id uuid.UUID) (*task.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %s: %w", id, sharederrors.ErrNotFound)
	}
	return t, nil
}

func (m *mockTaskRepository) Save(t *task.Task) (*task.Task, error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	m.tasks[t.ID] = t
	return t, nil
}

func (m *mockTaskRepository) Update(t *task.Task) (*task.Task, error) {
	m.tasks[t.ID] = t
	return t, nil
}

func (m *mockTaskRepository) Delete(id uuid.UUID) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepository) FindAll(_ uuid.UUID, _ uuid.UUID, _ bool, _ task.Filters) ([]*task.Task, int64, error) {
	return nil, 0, nil
}

func (m *mockTaskRepository) FindByContact(_ uuid.UUID) ([]*task.Task, error) { return nil, nil }
func (m *mockTaskRepository) FindByDeal(_ uuid.UUID) ([]*task.Task, error)    { return nil, nil }
func (m *mockTaskRepository) FindOverdue() ([]*task.Task, error)               { return nil, nil }

func isValidationError(err error) bool {
	return err != nil && fmt.Sprintf("%s", err) != "" &&
		containsErr(err, sharederrors.ErrValidation)
}

func containsErr(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
		} else {
			return false
		}
	}
	return false
}

// ---

func makeRctx(userID uuid.UUID) ctxutil.RequestCtx {
	return ctxutil.RequestCtx{
		UserID:   userID,
		TenantID: uuid.New(),
		Role:     user.RoleSeller,
	}
}

func TestTaskService_UpdateStatus_BlocksReopeningFinalTask(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	bus := events.New(10, log)
	svc := task.NewService(newMockRepo(), bus)

	assignedTo := uuid.New()
	rctx := makeRctx(assignedTo)

	created, err := svc.Create(rctx, task.CreateTaskDTO{
		Title:      "Ligar ao cliente",
		Priority:   task.PriorityHigh,
		AssignedTo: assignedTo,
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	_, err = svc.UpdateStatus(created.ID, rctx, task.StatusDone)
	if err != nil {
		t.Fatalf("mark done: %v", err)
	}

	_, err = svc.UpdateStatus(created.ID, rctx, task.StatusTodo)
	if err == nil {
		t.Fatal("expected error reopening a done task, got nil")
	}
	if !isValidationError(err) {
		t.Errorf("expected ErrValidation, got: %v", err)
	}
}

func TestTaskService_UpdateStatus_AllowsNormalTransitions(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	bus := events.New(10, log)
	svc := task.NewService(newMockRepo(), bus)

	assignedTo := uuid.New()
	rctx := makeRctx(assignedTo)

	created, _ := svc.Create(rctx, task.CreateTaskDTO{
		Title:      "Enviar proposta",
		Priority:   task.PriorityMedium,
		AssignedTo: assignedTo,
	})

	transitions := []task.Status{
		task.StatusInProgress,
		task.StatusDone,
	}

	for _, status := range transitions {
		updated, err := svc.UpdateStatus(created.ID, rctx, status)
		if err != nil {
			t.Errorf("transition to %s failed: %v", status, err)
		}
		if updated.Status != status {
			t.Errorf("expected status %s, got %s", status, updated.Status)
		}
	}
}
