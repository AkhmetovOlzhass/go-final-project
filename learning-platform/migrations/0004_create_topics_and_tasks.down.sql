DROP TRIGGER IF EXISTS trg_update_tasks ON tasks;
DROP TRIGGER IF EXISTS trg_update_topics ON topics;
DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS topics;

DROP TYPE IF EXISTS difficulty;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS school_class;
DROP TYPE IF EXISTS answer_type;
