CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create trigger function to validate savings amounts
CREATE OR REPLACE FUNCTION check_savings_amount()
RETURNS TRIGGER AS $$
BEGIN
    -- Prevent current_amount from exceeding target_amount
    IF NEW.current_amount > NEW.target_amount THEN
        RAISE EXCEPTION 'Current amount cannot exceed target amount';
END IF;

    -- Ensure savings amount is positive
    IF NEW.savings_amount <= 0 THEN
        RAISE EXCEPTION 'Savings amount must be positive';
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger function to auto-update plan status
CREATE OR REPLACE FUNCTION update_plan_status()
RETURNS TRIGGER AS $$
BEGIN
    -- Auto-complete plan if target reached
    IF NEW.current_amount >= NEW.target_amount AND NEW.status != 'completed' THEN
        NEW.status := 'completed';
        NEW.updated_at := NOW();
END IF;

    -- Re-activate plan if amount drops below target
    IF NEW.current_amount < NEW.target_amount AND NEW.status = 'completed' THEN
        NEW.status := 'active';
        NEW.updated_at := NOW();
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the savings_plans table
CREATE TABLE savings_plans (
                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

                                -- Reference to the user who owns this savings plan
                                -- Cascading delete: if user is deleted, their plans are automatically deleted
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

                                -- Public-facing unique reference code for the plan (shown to users)
                                -- Used for deposits/withdrawals and customer support references
                               plan_ref TEXT UNIQUE NOT NULL,

                                -- The total target amount the user wants to save (in local currency)
                                -- Example: 500000.00 for ₦500,000 target
                               target_amount DECIMAL(12, 2) NOT NULL,

                                -- Current accumulated amount in the plan (in local currency)
                                -- Starts at 0 and increments with each approved deposit
                               current_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,

                                -- Duration of the savings plan in days
                                -- Example: 90 for a 3-month plan
                               duration_days INTEGER NOT NULL CHECK (duration_days > 0),

                                -- How often the user intends to save (daily/weekly/monthly)
                                -- Determines the expected savings cadence
                               savings_frequency TEXT NOT NULL CHECK (savings_frequency IN ('daily', 'weekly', 'monthly')),

                                -- Fixed amount to be saved per interval (in local currency)
                                -- Example: If weekly, this would be the amount saved each week
                               savings_amount DECIMAL(12, 2) NOT NULL,

                                -- Current status of the savings plan
                                -- active: Plan is ongoing
                                -- completed: Reached target/maturity date
                                -- cancelled: Terminated before completion
                               status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
                               created_at TIMESTAMP DEFAULT NOW(),
                               updated_at TIMESTAMP DEFAULT NOW(),

                                -- Auto-calculated maturity date based on created_at + duration_days
                                -- This is a generated column that's stored physically
                               maturity_date TIMESTAMP GENERATED ALWAYS AS (created_at + (duration_days * INTERVAL '1 day')) STORED
);

-- Add index for performance
CREATE INDEX idx_savings_plans_user_id ON savings_plans(user_id);
CREATE INDEX idx_savings_plans_status ON savings_plans(status);
CREATE INDEX idx_savings_plans_maturity_date ON savings_plans(maturity_date);

-- Attach triggers
CREATE TRIGGER check_savings_amount_trigger
    BEFORE INSERT OR UPDATE ON savings_plans
                         FOR EACH ROW EXECUTE FUNCTION check_savings_amount();

CREATE TRIGGER update_plan_status_trigger
    BEFORE UPDATE OF current_amount ON savings_plans
    FOR EACH ROW EXECUTE FUNCTION update_plan_status();