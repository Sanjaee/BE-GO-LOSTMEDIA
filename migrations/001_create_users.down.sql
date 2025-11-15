-- Drop trigger and function
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_email_verified;
DROP INDEX IF EXISTS idx_password_reset_token;
DROP INDEX IF EXISTS idx_email_verify_token;
DROP INDEX IF EXISTS idx_banned_users;
DROP INDEX IF EXISTS idx_google_id;
DROP INDEX IF EXISTS idx_user_lookup;

-- Drop table
DROP TABLE IF EXISTS users;

