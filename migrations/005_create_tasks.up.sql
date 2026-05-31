CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE task_status   AS ENUM ('todo', 'in_progress', 'done', 'cancelled');

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(200) NOT NULL,
    description TEXT,
    priority    task_priority NOT NULL DEFAULT 'medium',
    status      task_status   NOT NULL DEFAULT 'todo',
    assigned_to UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_id  UUID REFERENCES contacts(id) ON DELETE SET NULL,
    deal_id     UUID REFERENCES deals(id) ON DELETE SET NULL,
    due_date    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_contact_id  ON tasks(contact_id);
CREATE INDEX idx_tasks_deal_id     ON tasks(deal_id);
CREATE INDEX idx_tasks_due_date    ON tasks(due_date) WHERE status NOT IN ('done', 'cancelled');
