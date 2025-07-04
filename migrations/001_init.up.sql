CREATE TABLE IF NOT EXISTS refresh_tokens (
    guid UUID PRIMARY KEY,
    token_hash TEXT NOT NULL,
    session_id TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);
