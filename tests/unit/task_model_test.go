package unit_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
)

func TestTask_IsOverdue(t *testing.T) {
	t.Parallel()

	yesterday := time.Now().Add(-24 * time.Hour)
	tomorrow := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name     string
		task     task.Task
		expected bool
	}{
		{
			name: "past due date and todo → overdue",
			task: task.Task{
				ID:         uuid.New(),
				AssignedTo: uuid.New(),
				Status:     task.StatusTodo,
				DueDate:    &yesterday,
			},
			expected: true,
		},
		{
			name: "future due date → not overdue",
			task: task.Task{
				ID:         uuid.New(),
				AssignedTo: uuid.New(),
				Status:     task.StatusTodo,
				DueDate:    &tomorrow,
			},
			expected: false,
		},
		{
			name: "past due date but done → not overdue",
			task: task.Task{
				ID:         uuid.New(),
				AssignedTo: uuid.New(),
				Status:     task.StatusDone,
				DueDate:    &yesterday,
			},
			expected: false,
		},
		{
			name: "no due date → never overdue",
			task: task.Task{
				ID:         uuid.New(),
				AssignedTo: uuid.New(),
				Status:     task.StatusTodo,
				DueDate:    nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.task.IsOverdue(); got != tt.expected {
				t.Errorf("IsOverdue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTaskStatus_IsFinal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status task.Status
		final  bool
	}{
		{task.StatusTodo, false},
		{task.StatusInProgress, false},
		{task.StatusDone, true},
		{task.StatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			t.Parallel()
			if got := tt.status.IsFinal(); got != tt.final {
				t.Errorf("Status(%q).IsFinal() = %v, want %v", tt.status, got, tt.final)
			}
		})
	}
}
