INSERT INTO threads (id, t_name, external_message_store, created_at, updated_at) VALUES
    ('t0', 'cat and mouse', 'f', datetime('now'), datetime('now')),
    ('t1', 'pure functions', 'f', datetime('now', '-1 day'), datetime('now', '-1 day')),
    ('t2', 'who was Douglas Engelbart', 'f', datetime('now', '-2 day'), datetime('now', '-2 day'));

INSERT INTO messages (id, m_role, content, created_at, thread_id) VALUES
    ('t0m0', 'user', 'do cats and mice really hate each other', datetime('now'), '0'),
    ('t0m1', 'assistant', 'yes they do', datetime('now'), '0'),
    ('t1m0', 'user', 'what are pure functions', datetime('now', '-1 day'), '1'),
    ('t1m1', 'assistant', 'pure functions are functions that completely deterministic in their inputs', datetime('now', '-1 day'), '1'),
    ('t1m2', 'assistant', 'they do not modify the state of the program', datetime('now', '-1 day'), '1'),
    ('t2m0', 'user', 'he was an american engineer and inventor', datetime('now', '-2 day'), '2'),
    ('t2m1', 'user', 'he invented the mouse', datetime('now', '-2 day'), '2'),
    ('t2m2', 'user', 'he was a pioneer in the field of human computer interaction', datetime('now', '-2 day'), '2');

INSERT INTO virtual_thread_names(thread_id, thread_name) VALUES
    ('t0', 'cat and mouse'),
    ('t1', 'pure functions'),
    ('t2', 'who was Douglas Engelbart');

INSERT INTO virtual_message_content(message_id, thread_id, message_content) VALUES
    ('t0m0', 't0', 'do cats and mice really hate each other'),
    ('t0m1', 't0', 'yes they do'),
    ('t1m0', 't1', 'what are pure functions'),
    ('t1m1', 't1', 'pure functions are functions that completely deterministic in their inputs'),
    ('t1m2', 't1', 'they do not modify the state of the program'),
    ('t2m0', 't2', 'he was an american engineer and inventor'),
    ('t2m1', 't2', 'he invented the mouse'),
    ('t2m2', 't2', 'he was a pioneer in the field of human computer interaction');
