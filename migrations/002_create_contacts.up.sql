CREATE TABLE contacts (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100)  NOT NULL,
    email      VARCHAR(255)  NOT NULL,
    phone      VARCHAR(20),
    company    VARCHAR(100),
    notes      TEXT,
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contacts_owner_id ON contacts(owner_id);
CREATE INDEX idx_contacts_email    ON contacts(email);
-- Index para pesquisa case-insensitive (ILIKE)
CREATE INDEX idx_contacts_name_search    ON contacts USING gin(name gin_trgm_ops);
CREATE INDEX idx_contacts_company_search ON contacts USING gin(company gin_trgm_ops);
