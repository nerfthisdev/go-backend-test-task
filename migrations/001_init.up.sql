CREATE TABLE auth (
    guid         UUID NOT NULL PRIMARY KEY,
    user_agent   TEXT NOT NULL,
    ip_address   TEXT NOT NULL,
    token        TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL
);
