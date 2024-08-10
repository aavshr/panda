package ui

import (
	"fmt"
	"os"
	"testing"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/store"
	"golang.org/x/term"
)

func TestModel_View(t *testing.T) {
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
		t.Fatalf("failed to get terminal size: %v", err)
	}

	m, _ := New(&Config{
		InitThreadsLimit: 10,
		MaxThreadsLimit:  100,
		Width:            width - 5,
		Height:           height - 10,
	}, mockStore)
	_ = m.Init()
	fmt.Println(m.View())
}
