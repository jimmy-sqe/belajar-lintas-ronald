CREATE TABLE IF NOT EXISTS users (
    id              CHAR(36) PRIMARY KEY,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(255),
    created_at      DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    created_by      CHAR(36),
    modified_at     DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    modified_by     CHAR(36)
);

CREATE INDEX idx_users_email ON users (email);
