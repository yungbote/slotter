-- Create 'company' table
CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Create 'role' table
--  - Has a JSON or text[] 'permissions' field so we can store a list of permissions
CREATE TABLE IF NOT EXISTS role (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    permissions JSONB NOT NULL DEFAULT '[]',
    company_id INT,
    CONSTRAINT fk_role_company
      FOREIGN KEY (company_id)
      REFERENCES company (id)
      ON DELETE CASCADE
);

-- Create 'warehouse' table
--  - We'll tie warehouse to a company
--  - Name is unique within the context of (company_id, name)
CREATE TABLE IF NOT EXISTS warehouse (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    company_id INT,
    CONSTRAINT fk_warehouse_company
      FOREIGN KEY (company_id)
      REFERENCES company (id)
      ON DELETE CASCADE
);

-- Create a unique index for (company_id, name) so warehouse names are only unique per company
CREATE UNIQUE INDEX IF NOT EXISTS warehouse_company_name_idx
  ON warehouse (company_id, name);

-- Create 'user' table
--  - Removes the old 'role' string
--  - Adds 'role_id' referencing 'role' table
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    company_id INT,
    role_id INT,
    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_user_company
      FOREIGN KEY (company_id)
      REFERENCES company (id)
      ON DELETE SET NULL,

    CONSTRAINT fk_user_role
      FOREIGN KEY (role_id)
      REFERENCES role (id)
      ON DELETE SET NULL
);

-- Create 'transactionrecord' table
--  - Link to 'warehouse' instead of 'company'
CREATE TABLE IF NOT EXISTS transactionrecord (
    id SERIAL PRIMARY KEY,
    warehouse_id INT NOT NULL,
    transaction_type TEXT,
    order_number TEXT,
    item_number TEXT,
    description TEXT,
    transaction_quantity INT,
    location TEXT,
    zone TEXT,
    carousel TEXT,
    row TEXT,
    shelf TEXT,
    bin TEXT,
    completed_date TIMESTAMPTZ,
    completed_by TEXT,
    completed_quantity INT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_tr_warehouse
      FOREIGN KEY (warehouse_id)
      REFERENCES warehouse (id)
      ON DELETE CASCADE
);

