ALTER TABLE tasks ADD COLUMN due_at datetime NULL;
ALTER TABLE tasks ADD COLUMN assignee_id varchar(36) NULL;
ALTER TABLE tasks ADD COLUMN enable_notice boolean NOT NULL DEFAULT false;
ALTER TABLE tasks ADD COLUMN reminder_at datetime NULL;
ALTER TABLE tasks ADD COLUMN type_reminder int NOT NULL DEFAULT 0;