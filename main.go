package main

import (
	_ "embed"
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui"
	"github.com/aavshr/panda/internal/ui/store"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	"log"
	"os"
	//"log"
	//"os"
	//"strings"
)

//go:embed internal/db/schema/init.sql
var dbSchemaInit string

//go:embed internal/db/schema/migrations.sql
var dbSchemaMigrations string

const (
	DefaultDataDirPath  = "/.local/share/panda/data"
	DefaultDatabaseName = "panda.db"
)

func main() {
	/*
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

		_, err := db.New(db.Config{
			DataDirPath: dataDirPath,
			DatabaseName: databaseName,
		}, &dbSchemaInit, &dbSchemaMigrations)
		if err != nil {
			log.Fatal("failed to initialize db: ", err)
		}
	*/

	testThreads := []*db.Thread{
		{
			ID:        "1",
			Name:      "Thread 1",
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-01-02",
		},
		{
			ID:        "2",
			Name:      "Thread 2",
			CreatedAt: "2024-01-03",
			UpdatedAt: "2024-01-02",
		},
	}
	mockStore := store.NewMock(testThreads)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	m, err := ui.New(&ui.Config{
		InitThreadsLimit: 10,
		MaxThreadsLimit:  100,
		Width:            width - 10,
		Height:           height - 10,
	}, mockStore)
	if err != nil {
		log.Fatal("ui.New: ", err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
