package main

import (
	_ "embed"
	"log"
	"os"

	"strings"

	"github.com/aavshr/panda/internal/config"
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/llm/openai"
	"github.com/aavshr/panda/internal/ui"
	"github.com/aavshr/panda/internal/ui/store"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

//go:embed internal/db/schema/init.sql
var dbSchemaInit string

//go:embed internal/db/schema/migrations.sql
var dbSchemaMigrations string

const (
	DefaultDatabaseName = "panda.db"
)

func initMockStore() *store.Mock {
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
	testMessages := []*db.Message{
		{
			ID:        "1",
			Role:      "user",
			ThreadID:  "1",
			Content:   "Thread 1\nMessage 1",
			CreatedAt: "2024-01-01",
		},
		{
			ID:        "2",
			Role:      "assistant",
			ThreadID:  "1",
			Content:   "Thread 1\nMessage 2",
			CreatedAt: "2024-01-02",
		},
		{
			ID:        "3",
			Role:      "user",
			ThreadID:  "2",
			Content:   "Thread 2\nMessage 1",
			CreatedAt: "2024-01-01",
		},
		{
			ID:        "4",
			Role:      "assistant",
			ThreadID:  "2",
			Content:   "Thread 2\nMessage 2",
			CreatedAt: "2024-01-02",
		},
	}

	return store.NewMock(testThreads, testMessages)
}

func main() {
	isDev := strings.ToLower(os.Getenv("PANDA_ENV")) == "dev"
	dataDirPath := config.GetDataDir()
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

	dbStore, err := db.New(db.Config{
		DataDirPath:  dataDirPath,
		DatabaseName: databaseName,
	}, &dbSchemaInit, &dbSchemaMigrations)
	if err != nil {
		log.Fatal("failed to initialize db: ", err)
	}

	openaiLLM := openai.New("")

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	m, err := ui.New(&ui.Config{
		InitThreadsLimit: 10,
		MaxThreadsLimit:  100,
		MessagesLimit:    50,
		Width:            width - 8,
		Height:           height - 10,
	}, dbStore, openaiLLM)
	if err != nil {
		log.Fatal("ui.New: ", err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
