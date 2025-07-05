CREATE TABLE users (
    guid UUID PRIMARY KEY
);

CREATE TABLE refresh_tokens (
    guid UUID PRIMARY KEY REFERENCES users(guid) ON DELETE CASCADE,

    token_hash TEXT NOT NULL,
    session_id TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);
