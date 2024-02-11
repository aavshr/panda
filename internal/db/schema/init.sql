CREATE TABLE IF NOT EXISTS threads (
    id TEXT PRIMARY KEY,
    t_name TEXT NOT NULL,
    external_message_store BOOL DEFAULT 'f',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    m_role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    thread_id TEXT REFERENCES threads(id) ON DELETE CASCADE
);

CREATE VIRTUAL TABLE virtual_thread_names USING fts5(
    thread_name,
    thread_id UNINDEXED
);

CREATE VIRTUAL TABLE virtual_message_content USING fts5(
    message_content,
    message_id UNINDEXED,
    thread_id UNINDEXED
);