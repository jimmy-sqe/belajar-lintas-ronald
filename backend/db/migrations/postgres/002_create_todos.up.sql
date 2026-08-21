CREATE TABLE IF NOT EXISTS todos (
    id              UUID PRIMARY KEY,
    title           VARCHAR(200) NOT NULL,
    description     VARCHAR(2000),
    due_date        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID NOT NULL,
    modified_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_by     UUID NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_todos_created_by ON todos (created_by);
CREATE INDEX IF NOT EXISTS idx_todos_created_by_created_at ON todos (created_by, created_at DESC);
