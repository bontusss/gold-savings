CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users' table
CREATE TABLE
  users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    total_savings NUMERIC(10, 2),
    total_withdrawn NUMERIC(10, 2),
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