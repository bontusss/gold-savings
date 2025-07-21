

CREATE TABLE
  savings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
  );

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the investments table
CREATE TABLE investments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id INTEGER NOT NULL REFERENCES investment_plans(id) ON DELETE CASCADE,
  reference_id TEXT NOT NULL, -- unique reference for the investment
  amount INTEGER NOT NULL CHECK (amount > 0), -- total amount invested
  interest NUMERIC(12, 2) NOT NULL DEFAULT 0.00, -- total interest earned on this investment
  interest_rate NUMERIC(5, 2) NOT NULL, -- snapshot of plan rate at time of investment
  status TEXT NOT NULL CHECK (
    status IN ('active', 'cancelled', 'pending')
  ),
  start_date TIMESTAMP, -- when the investment started
  end_date TIMESTAMP, -- when the investment ends
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create the trigger function to set end_date
CREATE OR REPLACE FUNCTION set_default_end_date()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.end_date IS NULL THEN
    NEW.end_date := NEW.created_at + INTERVAL '1 year';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the investments table
CREATE TRIGGER trg_set_end_date
BEFORE INSERT ON investments
FOR EACH ROW
EXECUTE FUNCTION set_default_end_date();


CREATE INDEX idx_investments_user_id ON investments (user_id);

CREATE INDEX idx_investments_plan_id ON investments (plan_id);

CREATE TABLE
  transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    investment_id UUID UNIQUE REFERENCES investments (id) ON DELETE SET NULL,
    type TEXT NOT NULL CHECK (type IN ('savings', 'investment')),
    status TEXT NOT NULL CHECK (
      status IN ('deposit', 'withdrawal', 'pending', 'declined')
    ),
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE
  payout_requests (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    account_name TEXT NOT NULL,
    bank_name TEXT NOT NULL,
    investment_id UUID UNIQUE REFERENCES investments (id) ON DELETE SET NULL,
    type TEXT NOT NULL CHECK (type IN ('savings', 'investment')),
    category TEXT NOT NULL CHECK (category IN ('deposit', 'withdrawal')),
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );