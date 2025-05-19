-- Drop triggers first
DROP TRIGGER IF EXISTS check_savings_amount_trigger ON savings_plans;
DROP TRIGGER IF EXISTS update_plan_status_trigger ON savings_plans;

-- Drop indexes
DROP INDEX IF EXISTS idx_savings_plans_user_id;
DROP INDEX IF EXISTS idx_savings_plans_status;
DROP INDEX IF EXISTS idx_savings_plans_maturity_date;

-- Drop the table
DROP TABLE IF EXISTS savings_plans;

-- Drop functions
DROP FUNCTION IF EXISTS check_savings_amount;
DROP FUNCTION IF EXISTS update_plan_status;

-- Note: We don't drop uuid-ossp extension as it might be used by other tables