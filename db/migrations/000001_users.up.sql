CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users' table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       first_name TEXT NOT NULL,
                       last_name TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       phone TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       account_number TEXT,
                       bank_name TEXT,
                       token_balance INTEGER DEFAULT 0,
                       is_active BOOLEAN DEFAULT TRUE,
                       is_admin Boolean Default FALSE,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);