-- Create task table
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(20) PRIMARY KEY,
    task_id VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    label VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL DEFAULT 'medium',
    html TEXT,
    lexical TEXT,
    due_at TIMESTAMPTZ NULL,
    assignee_id VARCHAR(36) NULL,
    enable_notice BOOLEAN NOT NULL DEFAULT false,
    reminder_at TIMESTAMPTZ NULL,
    type_reminder INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by  VARCHAR(20),
    updated_by  VARCHAR(20),
    deleted_by  VARCHAR(20)
);

-- Create indexes
CREATE INDEX  IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX  IF NOT EXISTS idx_tasks_label ON tasks(label);
CREATE INDEX  IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX  IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);
