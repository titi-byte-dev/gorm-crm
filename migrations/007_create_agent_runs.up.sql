CREATE TYPE agent_type   AS ENUM ('follow_up', 'deal_closer', 'task_router', 'summarize');
CREATE TYPE agent_mode   AS ENUM ('suggest', 'auto');
CREATE TYPE agent_status AS ENUM ('running', 'completed', 'failed');

CREATE TABLE agent_runs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_type   agent_type   NOT NULL,
    entity_type  VARCHAR(50)  NOT NULL,
    entity_id    UUID         NOT NULL,
    tenant_id    UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    runner_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode         agent_mode   NOT NULL DEFAULT 'suggest',
    status       agent_status NOT NULL DEFAULT 'running',
    actions      JSONB        NOT NULL DEFAULT '[]',
    summary      TEXT,
    tokens_used  INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_agent_runs_tenant_entity ON agent_runs(tenant_id, entity_type, entity_id);
CREATE INDEX idx_agent_runs_runner        ON agent_runs(runner_id);
CREATE INDEX idx_agent_runs_created_at    ON agent_runs(created_at DESC);
