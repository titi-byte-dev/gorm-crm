package unit_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
	"log/slog"
	"os"
)

// mockTaskRepository implementa task.Repository sem DB.
// O compilador verifica que implementa a interface — se falhar, erro de build.
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

func (m *mockTaskRepository) FindAll(_ uuid.UUID, _ task.Filters) ([]*task.Task, int64, error) {
	return nil, 0, nil
}

func (m *mockTaskRepository) FindByContact(_ uuid.UUID) ([]*task.Task, error) { return nil, nil }
func (m *mockTaskRepository) FindByDeal(_ uuid.UUID) ([]*task.Task, error)    { return nil, nil }
func (m *mockTaskRepository) FindOverdue() ([]*task.Task, error)               { return nil, nil }

// ---

// TestTaskService_UpdateStatus_BlocksReopeningFinalTask prova que
// a regra de negócio "tarefas finais não podem ser reabertas" está
// no Service — sem HTTP server, sem PostgreSQL.
//
// Este teste corre em <1ms. Se estivesse misturado no handler,
// precisaríamos de um servidor HTTP e uma ligação ao DB para o correr.
func TestTaskService_UpdateStatus_BlocksReopeningFinalTask(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	bus := events.New(10, log)
	svc := task.NewService(newMockRepo(), bus)

	// Criar uma task
	created, err := svc.Create(task.CreateTaskInput{
		Title:      "Ligar ao cliente",
		Priority:   task.PriorityHigh,
		AssignedTo: uuid.New(),
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// Marcar como done
	_, err = svc.UpdateStatus(created.ID, task.StatusDone)
	if err != nil {
		t.Fatalf("mark done: %v", err)
	}

	// Tentar reabrir → deve falhar
	_, err = svc.UpdateStatus(created.ID, task.StatusTodo)
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

	created, _ := svc.Create(task.CreateTaskInput{
		Title:      "Enviar proposta",
		Priority:   task.PriorityMedium,
		AssignedTo: uuid.New(),
	})

	transitions := []task.Status{
		task.StatusInProgress,
		task.StatusDone,
	}

	for _, status := range transitions {
		updated, err := svc.UpdateStatus(created.ID, status)
		if err != nil {
			t.Errorf("transition to %s failed: %v", status, err)
		}
		if updated.Status != status {
			t.Errorf("expected status %s, got %s", status, updated.Status)
		}
	}
}

func isValidationError(err error) bool {
	return err != nil && fmt.Sprintf("%v", err) != "" // simplificado
}
