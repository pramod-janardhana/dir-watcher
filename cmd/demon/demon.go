package demon

import (
	"github.com/knadh/koanf/providers/file"
	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/internal/process"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

func Run(fileEventDB, jobRunsDB database.Database) {
	{
		path, err := config.JonSchedulerConfigPath()
		if err != nil {
			zlog.Fatalf("could not get config.jobscheduler.windows.json file path: %v", err)
		}

		Watch(file.Provider(path), func() {
			config.ServiceConf.CronJobScheduler = *config.LoadJobSchedulerConfig()
			process.LoadScheduler(jobRunsDB, fileEventDB)
		})

		path, err = config.FileWatcherConfigPath()
		if err != nil {
			zlog.Fatalf("could not get config.filewatcher.windows.json file path: %v", err)
		}

		Watch(file.Provider(path), func() {
			config.ServiceConf.FileWatcher = *config.LoadFileWatcherConfig()
			process.LoadFileWatcher(fileEventDB)
		})
	}

	// file watcher configuration and registration
	process.LoadFileWatcher(fileEventDB)

	// cron jobs scheduler configuration and registration
	process.LoadScheduler(jobRunsDB, fileEventDB)
}
