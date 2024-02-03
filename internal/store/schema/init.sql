CREATE TABLE IF NOT EXISTS threads (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    external_message_store BOOL DEFAULT 'f',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    thread_id TEXT REFERENCES threads(id) ON DELETE CASCADE
);

CREATE VIRTUAL TABLE searchable_content USING fts4(
    thread_name TEXT NOT NULL,
    message_content TEXT NOT NULL,
    thread_id TEXT REFERENCES threads(id) ON DELETE CASCADE,
    message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
    UNIQUE(thread_id, message_id)
);