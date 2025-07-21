
-- Users' table
CREATE TABLE
  users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    total_savings INTEGER NOT NULL DEFAULT 0,
    total_savings_withdrawn INTEGER NOT NULL DEFAULT 0,
    total_investment_amount INTEGER NOT NULL DEFAULT 0,
    total_investment_withdrawn INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    reference_id TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    account_number TEXT,
    bank_name TEXT,
    token_balance INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_code TEXT,
    verification_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
  );