package store

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

const (
	dataDirPathPermission = 0755
)

type Config struct {
	DataDirPath string
	DatabaseName string
}

type Store struct {
	db *sqlx.DB
}

func New(config Config, schemaInit, migrations *string, ) (*Store, error) {
	if err := os.MkdirAll(config.DataDirPath, dataDirPathPermission); err != nil {
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
	_, err = db.Exec(*schemaInit)
	if err != nil {
		return nil, fmt.Errorf("could not init schemas, db.Exec: %w", err)
	}
	if migrations != nil  && *migrations != ""{
		_, err = db.Exec(*migrations)
		if err != nil {
			return nil, fmt.Errorf("could not run migrations, db.Exec: %w", err)
		}
	}
	return &Store{db: db}, nil
}