
CREATE TYPE difficulty AS ENUM ('EASY', 'MEDIUM', 'HARD', 'EXTREME');
CREATE TYPE task_status AS ENUM ('DRAFT', 'PUBLISHED', 'ARCHIVED');
CREATE TYPE school_class AS ENUM ('SEVEN', 'EIGHT', 'NINE', 'TEN', 'ELEVEN');
CREATE TYPE answer_type AS ENUM ('TEXT', 'NUMBER', 'FORMULA');


CREATE TABLE topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    parent_id UUID NULL,
    school_class school_class NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    body_md TEXT NOT NULL,
    difficulty difficulty NOT NULL DEFAULT 'EASY',
    status task_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    topic_id UUID NOT NULL,
    author_id UUID NOT NULL,

    official_solution TEXT,
    correct_answer TEXT,
    answer_type answer_type NOT NULL DEFAULT 'TEXT',
    image_url TEXT,

    CONSTRAINT fk_topic FOREIGN KEY (topic_id)
        REFERENCES topics (id) ON DELETE CASCADE,

    CONSTRAINT fk_author FOREIGN KEY (author_id)
        REFERENCES users (id) ON DELETE CASCADE
);


CREATE INDEX idx_topics_parent ON topics(parent_id);
CREATE INDEX idx_tasks_topic ON tasks(topic_id);
CREATE INDEX idx_tasks_author ON tasks(author_id);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_topics
BEFORE UPDATE ON topics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_update_tasks
BEFORE UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
