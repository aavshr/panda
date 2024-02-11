package db

// Thread represents a chat thread that contains messages
type Thread struct {
	ID        string `db:"id"`
	Name string  `db:"t_name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	ExternalMessageStore bool `db:"external_message_store"`
}

// Message represents a chat message which is part of a thread
type Message struct {
	ID   string `db:"id"`
	Role string `db:"m_role"`
	Content string `db:"content"`
	CreatedAt string `db:"created_at"`
	ThreadID string `db:"thread_id"`
}