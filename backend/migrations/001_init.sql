-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ═══════════════════════════════════════════════════════════════
-- RBAC: Users, Roles, Permissions
-- ═══════════════════════════════════════════════════════════════

-- Users table (no role column)
CREATE TABLE IF NOT EXISTS users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL UNIQUE,
    password     VARCHAR(255) NOT NULL,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active = TRUE;

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_active ON roles(is_active) WHERE is_active = TRUE;

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource        VARCHAR(100) NOT NULL,
    action          VARCHAR(50) NOT NULL,
    endpoint_path   TEXT NOT NULL,
    endpoint_method VARCHAR(10) NOT NULL,
    description     TEXT,
    hide            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT permissions_resource_action_unique UNIQUE(resource, action)
);

CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);
CREATE INDEX IF NOT EXISTS idx_permissions_endpoint ON permissions(endpoint_path, endpoint_method);

-- Role-Permission junction (many-to-many)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_perm ON role_permissions(permission_id);

-- User-Role junction (many-to-many)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role_id);

-- ═══════════════════════════════════════════════════════════════
-- Suppliers
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS suppliers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code         VARCHAR(20) NOT NULL UNIQUE,
    supplier_no  VARCHAR(20) NOT NULL UNIQUE,
    name         VARCHAR(255) NOT NULL,
    alias        VARCHAR(100),
    logo_url     TEXT,
    address      TEXT,
    city         VARCHAR(100),
    province     VARCHAR(100),
    country      VARCHAR(100) DEFAULT 'Indonesia',
    postal_code  VARCHAR(20),
    phone        VARCHAR(50),
    email        VARCHAR(255),
    website      VARCHAR(255),
    status       VARCHAR(50) NOT NULL DEFAULT 'draft',
    stage        VARCHAR(50) NOT NULL DEFAULT 'draft',
    sla_hours    INT NOT NULL DEFAULT 72,
    is_blocked   BOOLEAN NOT NULL DEFAULT FALSE,
    block_reason TEXT,
    notes        TEXT,
    created_by   UUID,
    updated_by   UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_suppliers_status ON suppliers(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers(deleted_at);

-- Supplier contacts
CREATE TABLE IF NOT EXISTS supplier_contacts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id  UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    position     VARCHAR(100),
    phone        VARCHAR(50),
    mobile       VARCHAR(50),
    email        VARCHAR(255),
    is_primary   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_supplier_contacts_supplier_id ON supplier_contacts(supplier_id);

-- Supplier addresses (multiple per supplier: Head Office, Branch, etc.)
CREATE TABLE IF NOT EXISTS supplier_addresses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id  UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,  -- e.g. "Head Office", "Branch Office"
    address      TEXT NOT NULL,
    city         VARCHAR(100),
    province     VARCHAR(100),
    country      VARCHAR(100) DEFAULT 'Indonesia',
    postal_code  VARCHAR(20),
    is_main      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_supplier_addresses_supplier_id ON supplier_addresses(supplier_id);

-- Supplier groups (classification tags: Industry, Telkom Group, etc.)
CREATE TABLE IF NOT EXISTS supplier_groups (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id  UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    group_name   VARCHAR(100) NOT NULL,  -- e.g. "Industry", "Telkom Group"
    value        VARCHAR(100) NOT NULL,  -- e.g. "Manufacture", "Non Telkom Group"
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_supplier_groups_supplier_id ON supplier_groups(supplier_id);

-- Supplier materials
CREATE TABLE IF NOT EXISTS supplier_materials (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id    UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    material_group VARCHAR(100) NOT NULL,
    material_id    VARCHAR(100) NOT NULL,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_supplier_materials_supplier_id ON supplier_materials(supplier_id);

-- Supplier performance ratings
CREATE TABLE IF NOT EXISTS supplier_performance_ratings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id     UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    price_rating    INT NOT NULL CHECK (price_rating BETWEEN 1 AND 5),
    delivery_rating INT NOT NULL CHECK (delivery_rating BETWEEN 1 AND 5),
    notes           TEXT,
    reviewed_by     VARCHAR(255),
    reviewed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ratings_supplier_id ON supplier_performance_ratings(supplier_id);

-- Supplier stage histories
CREATE TABLE IF NOT EXISTS supplier_stage_histories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id  UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    from_stage   VARCHAR(50),
    to_stage     VARCHAR(50) NOT NULL,
    notes        TEXT,
    changed_by   VARCHAR(255),
    elapsed_ms   BIGINT DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stage_histories_supplier_id ON supplier_stage_histories(supplier_id);

-- Supplier invoices (outstandings)
CREATE TABLE IF NOT EXISTS supplier_invoices (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id    UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    invoice_number VARCHAR(50)   NOT NULL,
    project_name   VARCHAR(255),
    amount         NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency       VARCHAR(10)   NOT NULL DEFAULT 'IDR',
    invoice_date   DATE          NOT NULL,
    due_date       DATE          NOT NULL,
    paid_date      DATE,
    status         VARCHAR(30)   NOT NULL DEFAULT 'unpaid',
    paid_amount    NUMERIC(18,2) NOT NULL DEFAULT 0,
    notes          TEXT,
    created_by     UUID,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT supplier_invoices_number_unique UNIQUE (supplier_id, invoice_number)
);

CREATE INDEX IF NOT EXISTS idx_invoices_supplier_id ON supplier_invoices(supplier_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status      ON supplier_invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date    ON supplier_invoices(due_date);

-- ═══════════════════════════════════════════════════════════════
-- Seed Data: Roles, Permissions, Admin User
-- ═══════════════════════════════════════════════════════════════

-- Insert system roles
INSERT INTO roles (name, description, is_system, is_active)
VALUES
    ('admin',    'Full system access',                      TRUE, TRUE),
    ('manager',  'Manage suppliers and approve workflows',  TRUE, TRUE),
    ('viewer',   'Read-only access',                        TRUE, TRUE),
    ('supplier', 'Supplier portal access',                  TRUE, TRUE)
ON CONFLICT (name) DO NOTHING;

-- Insert permissions
INSERT INTO permissions (resource, action, endpoint_path, endpoint_method, description, hide)
VALUES
    -- Supplier CRUD
    ('supplier', 'read',   '/suppliers',           'GET',    'List and view suppliers',          FALSE),
    ('supplier', 'create', '/suppliers',           'POST',   'Create new supplier',              FALSE),
    ('supplier', 'update', '/suppliers/:id',       'PUT',    'Update supplier details',          FALSE),
    ('supplier', 'delete', '/suppliers/:id',       'DELETE', 'Soft delete supplier',             FALSE),
    ('supplier', 'block',  '/suppliers/:id/block', 'POST',   'Block or unblock supplier',        FALSE),
    ('supplier', 'export', '/suppliers/export',    'GET',    'Export supplier data',             FALSE),

    -- Supplier Materials
    ('material', 'read',   '/suppliers/:id/materials', 'GET', 'View supplier materials',       FALSE),
    ('material', 'update', '/suppliers/:id/materials', 'PUT', 'Update supplier material list', FALSE),

    -- Supplier Ratings
    ('rating', 'read',   '/suppliers/:id/ratings', 'GET',  'View performance ratings',      FALSE),
    ('rating', 'create', '/suppliers/:id/ratings', 'POST', 'Add performance rating',        FALSE),

    -- Workflow
    ('workflow', 'advance', '/suppliers/:id/next-stage', 'POST', 'Advance supplier stage', FALSE),

    -- Review & Approval
    ('review', 'approve', '/reviews/:id/approve', 'POST', 'Approve supplier review',       TRUE),
    ('review', 'reject',  '/reviews/:id/reject',  'POST', 'Reject supplier review',        TRUE)
ON CONFLICT (resource, action) DO NOTHING;

-- Map permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.resource IN ('supplier', 'material', 'rating', 'workflow', 'review')
    AND p.action NOT IN ('delete')
WHERE r.name = 'manager'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.action = 'read'
WHERE r.name = 'viewer'
ON CONFLICT DO NOTHING;

-- Create system admin user (password: Admin@123)
INSERT INTO users (name, email, password, is_active)
VALUES (
    'System Admin',
    'admin@erp.local',
    '$2a$10$6InfEA2HKH9PgEIepo3qfuJi0mf7T2.PNP9l/0BwzS/UAoVAIi44K',
    TRUE
) ON CONFLICT (email) DO NOTHING;

-- Assign admin role to system admin
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@erp.local'
  AND r.name = 'admin'
ON CONFLICT DO NOTHING;
