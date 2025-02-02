package scheduler

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextDuration(t *testing.T) {
	now := time.Date(2024, 2, 1, 12, 0, 0, 0, time.Local)

	tests := []struct {
		name     string
		last     *time.Time
		interval string
		want     time.Duration
	}{
		{
			name:     "daily with no last run",
			last:     nil,
			interval: "daily",
			want:     0,
		},
		{
			name:     "daily with last run yesterday",
			last:     timePtr(now.Add(-24 * time.Hour)),
			interval: "daily",
			want:     0,
		},
		{
			name:     "weekly with no last run",
			last:     nil,
			interval: "weekly",
			want:     0,
		},
		{
			name:     "weekly with last run 3 days ago",
			last:     timePtr(now.Add(-72 * time.Hour)),
			interval: "weekly",
			want:     96 * time.Hour,
		},
		{
			name:     "monthly with no last run",
			last:     nil,
			interval: "monthly",
			want:     0,
		},
		{
			name:     "monthly with last run 15 days ago",
			last:     timePtr(now.Add(-360 * time.Hour)),
			interval: "monthly",
			want:     360 * time.Hour,
		},
		{
			name:     "debug mode",
			last:     timePtr(now.Add(-15 * time.Second)),
			interval: "debug",
			want:     15 * time.Second,
		},
		{
			name:     "disabled mode",
			last:     timePtr(now),
			interval: "disabled",
			want:     -1,
		},
		{
			name:     "invalid schedule",
			last:     timePtr(now),
			interval: "invalid",
			want:     -1,
		},
		{
			name:     "overdue daily update",
			last:     timePtr(now.Add(-48 * time.Hour)),
			interval: "daily",
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNextDuration(tt.last, now, tt.interval)
			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}

func TestLastRunFileOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scheduler_test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	testTime := time.Date(2024, 2, 1, 12, 0, 0, 0, time.Local)

	saveLastRun(tmpDir, testTime)

	loaded := loadLastRun(tmpDir)
	assert.Equal(t, testTime, *loaded)

	nonExistentDir := filepath.Join(tmpDir, "nonexistent")
	loaded = loadLastRun(nonExistentDir)
	assert.Nil(t, loaded)

	corruptFile := filepath.Join(tmpDir, ".lastupdate")
	err = os.WriteFile(corruptFile, []byte("not-a-number"), 0644)
	assert.Nil(t, err)

	loaded = loadLastRun(tmpDir)
	assert.Nil(t, loaded)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
