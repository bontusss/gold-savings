CREATE TABLE
  admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
  );

CREATE TABLE
  investment_plans (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    interest_rate NUMERIC(5, 2) NOT NULL, -- e.g. 12.50 (%)
    min_amount NUMERIC(12, 2) NOT NULL CHECK (min_amount >= 0),
    max_amount NUMERIC(12, 2) NOT NULL CHECK (max_amount > min_amount),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE INDEX idx_investment_plans_name ON investment_plans (name);