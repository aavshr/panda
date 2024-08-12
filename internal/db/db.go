package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Config struct {
	DataDirPath  string
	DatabaseName string
}

type Store struct {
	db *sqlx.DB
}

func New(config Config, schemaInit, migrations *string) (*Store, error) {
	// 0755 = rwxr-xr-x
	if err := os.MkdirAll(config.DataDirPath, 0755); err != nil {
		return nil, fmt.Errorf("could not make data dir, os.MkdirAll: %w", err)
	}
	f, err := os.Create(filepath.Join(config.DataDirPath, config.DatabaseName))
	if err != nil {
		return nil, fmt.Errorf("could not create database file, os.OpenFile: %w", err)
	}
	defer f.Close()

	db, err := sqlx.Open("sqlite3", filepath.Join(config.DataDirPath, config.DatabaseName))
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}
	if schemaInit != nil && *schemaInit != "" {
		_, err = db.Exec(*schemaInit)
		if err != nil {
			return nil, fmt.Errorf("could not init schemas, db.Exec: %w", err)
		}
	}
	if migrations != nil && *migrations != "" {
		_, err = db.Exec(*migrations)
		if err != nil {
			return nil, fmt.Errorf("could not run migrations, db.Exec: %w", err)
		}
	}
	return &Store{db: db}, nil
}

func (s *Store) Begin() (*sqlx.Tx, error) {
	return s.db.Beginx()
}

func (s *Store) ListLatestThreadsPaginated(offset, limit int) ([]*Thread, error) {
	var threads []*Thread
	err := s.db.Select(&threads, "SELECT * FROM threads ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return threads, fmt.Errorf("db.Select: %w", err)
	}
	return threads, nil
}

func (s *Store) CreateThreadTx(tx *sqlx.Tx, thread *Thread) error {
	query := `INSERT INTO threads (id, t_name, created_at, updated_at, external_message_store) 
			VALUES (:id, :t_name, :created_at, :updated_at, :external_message_store)`
	if _, err := tx.Exec(query, thread); err != nil {
		return fmt.Errorf("tx.NamedExec: %w", err)
	}
	query = `INSERT INTO virtual_thread_names (thread_id, thread_name) VALUES ($1, $2)`
	if _, err := tx.Exec(query, thread.ID, thread.Name); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (s *Store) UpdateThreadNameTx(tx *sqlx.Tx, threadID, name string) error {
	if _, err := tx.Exec("UPDATE threads SET t_name = $1 WHERE id = $2", name, threadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	if _, err := tx.Exec("UPDATE virtual_thread_names SET thread_name = $1 WHERE thread_id = $2", name, threadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (s *Store) DeleteThreadTx(tx *sqlx.Tx, threadID string) error {
	if _, err := tx.Exec("DELETE FROM threads WHERE id = $1", threadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM virtual_thread_names WHERE thread_id = $1", threadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM virtual_message_content WHERE thread_id = $1", threadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (s *Store) DeleteAllThreadsTx(tx *sqlx.Tx) error {
	if _, err := tx.Exec("DELETE FROM threads"); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM virtual_thread_names"); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM virtual_message_content"); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (s *Store) CreateMessageTx(tx *sqlx.Tx, message *Message) error {
	query := `INSERT INTO messages (id, m_role, content, created_at, thread_id) 
	VALUES (:id, :m_role, :content, :created_at, :thread_id)`
	if _, err := tx.NamedExec(query, message); err != nil {
		return fmt.Errorf("tx.NamedExec: %w", err)
	}
	query = `UPDATE threads SET updated_at = DATETIME('now') WHERE id = $1`
	if _, err := tx.Exec(query, message.ThreadID); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (s *Store) ListMessagesByThreadIDPaginated(threadID string, offset, limit int) ([]*Message, error) {
	var messages []*Message
	err := s.db.Select(&messages, "SELECT * FROM messages WHERE thread_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", threadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not select messages, db.Select: %w", err)
	}
	return messages, nil
}

func (s *Store) SearchThreadNamesPaginated(term string, offset, limit int) ([]*Thread, error) {
	var threads []*Thread
	query := `SELECT T.* FROM virtual_thread_names VTN INNER JOIN threads T ON VTN.thread_id = T.id WHERE VTN.thread_name MATCH $1 ORDER BY rank LIMIT $2 OFFSET $3`
	if err := s.db.Select(&threads, query, term, limit, offset); err != nil {
		return nil, fmt.Errorf("could not select from virtual thread names, db.Select: %w", err)
	}
	return threads, nil
}

func (s *Store) SearchMessageContentPaginated(term string, offset, limit int) ([]*Message, error) {
	var threads []*Message
	query := `SELECT M.* FROM virtual_message_content VMC INNER JOIN messages M ON VMC.message_id = M.id WHERE VMC.message_content MATCH $1 ORDER BY rank LIMIT $2 OFFSET $3`
	if err := s.db.Select(&threads, query, term, limit, offset); err != nil {
		return nil, fmt.Errorf("could not select from virtual message content, db.Select: %w", err)
	}
	return threads, nil
}
