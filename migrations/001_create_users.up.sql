-- Create users table with authentication fields
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    userId VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::text,
    googleId VARCHAR(255) UNIQUE,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    profilePic TEXT,
    bio TEXT,
    createdAt TIMESTAMP DEFAULT NOW(),
    updatedAt TIMESTAMP DEFAULT NOW(),
    followersCount INTEGER DEFAULT 0,
    followingCount INTEGER DEFAULT 0,
    role VARCHAR(50) DEFAULT 'member',
    star INTEGER DEFAULT 0,
    isBanned BOOLEAN DEFAULT false,
    banReason TEXT,
    bannedBy VARCHAR(36),
    postsCount INTEGER DEFAULT 0,
    isEmailVerified BOOLEAN DEFAULT false,
    emailVerificationToken VARCHAR(255),
    emailVerificationExpires TIMESTAMP,
    passwordResetToken VARCHAR(255),
    passwordResetExpires TIMESTAMP,
    lastLoginAt TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_lookup ON users(username, email);
CREATE INDEX IF NOT EXISTS idx_google_id ON users(googleId) WHERE googleId IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_banned_users ON users(isBanned);
CREATE INDEX IF NOT EXISTS idx_email_verify_token ON users(emailVerificationToken) WHERE emailVerificationToken IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_password_reset_token ON users(passwordResetToken) WHERE passwordResetToken IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_email_verified ON users(isEmailVerified);

-- Create function to update updatedAt
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updatedAt = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updatedAt
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

