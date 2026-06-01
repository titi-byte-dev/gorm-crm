-- Cria tabela de organizações (tenants)
CREATE TABLE organizations (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Adiciona organização ao utilizador
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);

-- Adiciona tenant_id e índices em todas as entidades
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE leads    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE deals    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE tasks    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_contacts_tenant_id ON contacts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_leads_tenant_id    ON leads(tenant_id);
CREATE INDEX IF NOT EXISTS idx_deals_tenant_id    ON deals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_id    ON tasks(tenant_id);

-- Índice único composto: email por organização (não global)
DROP INDEX IF EXISTS contacts_email_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_contacts_email_owner ON contacts(email, owner_id);
