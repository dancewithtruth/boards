CREATE TABLE IF NOT EXISTS email_verifications(
    id UUID PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    user_id UUID REFERENCES users(id),
    is_verified BOOLEAN,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);