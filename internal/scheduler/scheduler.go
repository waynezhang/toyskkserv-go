package scheduler

import (
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/files"
)

func StartUpdateWatcher(cfg *config.Config) {
	cron := map[string]string{
		"daily":   "0 0 * * *",
		"weekly":  "0 0 * * 1",
		"monthly": "0 0 1 * *",
		"debug":   "* * * * *",
	}[cfg.UpdateSchedule]
	if cron == "" {
		slog.Info("Update schedule is disabled", "schedule", cfg.UpdateSchedule)
		return
	}
	slog.Info("Update schedule", "cron", cron)

	sh, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("Failed to start update watcher", "err", err)
		return
	}

	job := gocron.CronJob(cron, false)
	sh.NewJob(job, gocron.NewTask(func() {
		slog.Info("Checking updates")
		files.UpdateDictionaries(cfg.Dictionaries, cfg.DictionaryDirectory)
	}))
	sh.Start()
}
