-- Drop task table
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_label;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP TABLE IF EXISTS tasks;
