package main

import (
	_ "embed"
	"github.com/aavshr/panda/internal/store"
	"log"
	"os"
	"strings"
)

//go:embed internal/store/schema/init.sql
var dbSchemaInit string

//go:embed internal/store/schema/migrations.sql
var dbSchemaMigrations string

const (
	DefaultDataDirPath = "/.local/share/panda/data"
	DefaultDatabaseName = "panda.db"
)

func main() {
	isDev := strings.ToLower(os.Getenv("PANDA_ENV")) == "dev"

	dataDirPath := DefaultDataDirPath
	databaseName := DefaultDatabaseName
	if isDev {
		devDataDirPath := os.Getenv("PANDA_DATA_DIR_PATH")
		devDatabaseName := os.Getenv("PANDA_DATABASE_NAME")
		if devDataDirPath != "" {
			dataDirPath = devDataDirPath
		}
		if devDatabaseName != "" {
			databaseName = devDatabaseName
		}
	}

	_, err := store.New(store.Config{
		DataDirPath: dataDirPath,
		DatabaseName: databaseName,
	}, &dbSchemaInit, &dbSchemaMigrations)
	if err != nil {
		log.Fatal("failed to initialize store: ", err)
	}
	return
}