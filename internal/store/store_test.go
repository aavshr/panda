package store

import (
	_ "embed"
	"os"
	"testing"
)

var (
	//go:embed schema/init.sql
	schemaInit string
	//go:embed testdata/populate.sql
	populateTestData string
)

func TestIntegrationStoreSearch(t *testing.T) {
	if testing.Short(){
		t.Skip("skipping integration test")
	}

	tmpDirPath, err := os.MkdirTemp("/tmp", "store")
	if err != nil {
		t.Fatalf("failed to create temp for storage dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDirPath); err != nil {
			t.Logf("WARNING: failed to remove temp dir: %v", err)
		}
	}()

	store, err := New(Config{
		DataDirPath: tmpDirPath,
		DatabaseName: "test.db",
	}, &schemaInit, &populateTestData)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	testCases := []struct{
		searchTerm string
		expectedThreadIDs map[string]struct{}
		expectedMessageIDs map[string]struct{}
	}{
		{
			"mouse",
			map[string]struct{}{
				"t0": {},
			},
			map[string]struct{}{
				"t2m1": {},
			},
		},
		{
			"functions",
			map[string]struct{}{
				"t1": {},
			},
			map[string]struct{}{
				"t1m1": {},
				"t1m0": {},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.searchTerm, func(t *testing.T) {
			threads, err := store.SearchThreadNamesPaginated(tc.searchTerm, 0, 10)
			if err != nil {
				t.Errorf("failed to search thread names: %v", err)
			}
			for _, thread := range threads {
				thread := thread
				if _, ok := tc.expectedThreadIDs[thread.ID]; !ok {
					t.Errorf("unexpected thread with id '%s' and name '%s'", thread.ID, thread.Name)
				}
			}
			messages, err := store.SearchMessageContentPaginated(tc.searchTerm, 0, 10)
			if err != nil {
				t.Errorf("failed to search message content: %v", err)
			}
			for _, message := range messages {
				message := message
				if _, ok := tc.expectedMessageIDs[message.ID]; !ok {
					t.Errorf("unexpected message with id '%s' and content '%s'", message.ID, message.Content)
				}
			}
		})
	}
}