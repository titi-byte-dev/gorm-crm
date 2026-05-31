CREATE TYPE deal_stage AS ENUM ('proposal', 'negotiation', 'won', 'lost');

CREATE TABLE deals (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title      VARCHAR(200) NOT NULL,
    value      DECIMAL(12,2) NOT NULL DEFAULT 0,
    stage      deal_stage NOT NULL DEFAULT 'proposal',
    lead_id    UUID REFERENCES leads(id) ON DELETE SET NULL,
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    closed_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deals_owner_id   ON deals(owner_id);
CREATE INDEX idx_deals_contact_id ON deals(contact_id);
CREATE INDEX idx_deals_stage      ON deals(stage);
