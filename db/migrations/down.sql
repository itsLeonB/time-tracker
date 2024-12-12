DROP INDEX IF EXISTS idx_task_logs_task_id;
DROP INDEX IF EXISTS idx_tasks_number;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP TABLE IF EXISTS task_logs;
DROP TYPE task_log_action;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;