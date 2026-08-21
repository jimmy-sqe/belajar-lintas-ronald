CREATE TABLE IF NOT EXISTS todos (
    id              CHAR(36) PRIMARY KEY,
    title           VARCHAR(200) NOT NULL,
    description     VARCHAR(2000),
    due_date        DATETIME(6),
    created_at      DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    created_by      CHAR(36) NOT NULL,
    modified_at     DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    modified_by     CHAR(36) NOT NULL
);

CREATE INDEX idx_todos_created_by ON todos (created_by);
CREATE INDEX idx_todos_created_by_created_at ON todos (created_by, created_at);
