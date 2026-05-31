CREATE TYPE lead_status AS ENUM ('new', 'contacted', 'qualified', 'lost');

CREATE TABLE leads (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title      VARCHAR(200) NOT NULL,
    value      DECIMAL(12,2) NOT NULL DEFAULT 0,
    status     lead_status NOT NULL DEFAULT 'new',
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_owner_id   ON leads(owner_id);
CREATE INDEX idx_leads_contact_id ON leads(contact_id);
CREATE INDEX idx_leads_status     ON leads(status);
