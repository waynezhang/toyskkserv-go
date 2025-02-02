package scheduler

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/tcp"
)

func StartUpdateWatcher(cfg *config.Config) {
	go func() {
		slog.Info("Update schedule", "schedule", cfg.UpdateSchedule)

		wait := make(chan struct{})
		for {
			last := loadLastRun(cfg.DictionaryDirectory)
			if last != nil {
				slog.Info("Last update", "time", last)
			} else {
				slog.Info("Last update", "time", "never")
			}

			now := time.Now()
			d := getNextDuration(last, now, cfg.UpdateSchedule)
			if d == -1 {
				break
			}

			slog.Info("Next run", "duration", d)
			time.AfterFunc(d, func() {
				run(cfg)
				saveLastRun(cfg.DictionaryDirectory, now)
				wait <- struct{}{}
			})
			<-wait
		}
	}()
}

func getNextDuration(last *time.Time, now time.Time, interval string) time.Duration {
	// daily, weekly, monthly, debug, disabled
	var d time.Duration

	switch interval {
	case "daily":
		d = time.Hour * 24
	case "weekly":
		d = time.Hour * 24 * 7
	case "monthly":
		d = time.Hour * 24 * 30
	case "debug":
		d = time.Second * 30
	case "disabled":
		return -1
	default:
		slog.Error("Invalid schedule", "schedule", interval)
		return -1
	}

	if last == nil {
		// first time
		return 0
	}

	next := last.Add(d)
	if next.Before(now) {
		return 0
	}

	return next.Sub(now)
}

func loadLastRun(dir string) *time.Time {
	path := filepath.Join(dir, ".lastupdate")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil
	}

	t := time.UnixMilli(val)
	return &t
}

func saveLastRun(dir string, t time.Time) {
	path := filepath.Join(dir, ".lastupdate")

	n := t.UnixMilli()
	s := strconv.FormatInt(n, 10)
	os.WriteFile(path, []byte(s), 0644)
}

func run(cfg *config.Config) {
	slog.Info("Start updating")
	tcp.SendReloadCommand(cfg.ListenAddr)
	slog.Info("Finished updating")
}
