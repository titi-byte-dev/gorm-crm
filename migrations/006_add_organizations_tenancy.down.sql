DROP INDEX IF EXISTS idx_contacts_email_owner;
DROP INDEX IF EXISTS idx_tasks_tenant_id;
DROP INDEX IF EXISTS idx_deals_tenant_id;
DROP INDEX IF EXISTS idx_leads_tenant_id;
DROP INDEX IF EXISTS idx_contacts_tenant_id;

ALTER TABLE tasks    DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE deals    DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE leads    DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE contacts DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_users_organization_id;
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;

DROP TABLE IF EXISTS organizations;
