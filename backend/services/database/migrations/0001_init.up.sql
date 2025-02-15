
-- 1) Create the "company" table
CREATE TABLE IF NOT EXISTS company (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- 2) Create the "warehouse" table
CREATE TABLE IF NOT EXISTS warehouse (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    company_id INT,
    -- Timestamps not shown in your struct, but you can add them if you want, e.g.
    -- created_at TIMESTAMPTZ DEFAULT now(),
    -- updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_warehouse_company
        FOREIGN KEY (company_id)
        REFERENCES company (id)
        ON DELETE CASCADE
);

-- 3) Create the "location" table
CREATE TABLE IF NOT EXISTS location (
    id                 SERIAL PRIMARY KEY,
    warehouse_id       INT  NOT NULL,
    parent_location_id INT  NULL,
    location_name      TEXT NOT NULL,
    location_type      TEXT NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT now(),
    updated_at         TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_location_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouse (id)
        ON DELETE CASCADE,

    -- For parent_location_id referencing itself, pick a policy (e.g. SET NULL, CASCADE, etc.)
    CONSTRAINT fk_location_parent
        FOREIGN KEY (parent_location_id)
        REFERENCES location (id)
        ON DELETE SET NULL
);

-- 4) Create the "user" table
CREATE TABLE IF NOT EXISTS "user" (
    id            SERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    company_id    INT,
    created_at    TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_user_company
        FOREIGN KEY (company_id)
        REFERENCES company (id)
        ON DELETE SET NULL
);

-- 5) Create the "transaction_record" table
CREATE TABLE IF NOT EXISTS transaction_record (
    id                    SERIAL PRIMARY KEY,
    warehouse_id          INT  NOT NULL,
    location_id           INT  NOT NULL,
    transaction_type      TEXT,
    order_number          TEXT,
    item_number           TEXT,
    description           TEXT,
    transaction_quantity  INT,
    completed_date        TIMESTAMPTZ,
    completed_quantity    INT,
    created_at            TIMESTAMPTZ DEFAULT now(),
    updated_at            TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_tr_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouse (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_tr_location
        FOREIGN KEY (location_id)
        REFERENCES location (id)
        ON DELETE CASCADE
);

-- Optional: Create any useful indexes for quick lookups
-- (Below are just examples; adjust to your usage patterns)
CREATE INDEX IF NOT EXISTS idx_warehouse_id ON warehouse (company_id);
CREATE INDEX IF NOT EXISTS idx_location_warehouse_id ON location (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_location_parent_id    ON location (parent_location_id);
CREATE INDEX IF NOT EXISTS idx_tr_warehouse_id       ON transaction_record (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_tr_location_id        ON transaction_record (location_id);

