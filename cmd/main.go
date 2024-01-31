package main

import (
	"os"

	"github.com/pramod-janardhana/dir-watcher/cmd/demon"
	"github.com/pramod-janardhana/dir-watcher/cmd/restserver"
	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/database/sqlite3"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

func dbSetup(dbPath string, initScriptPath string) (database.Database, error) {
	fileEventDB, err := sqlite3.Open(dbPath)
	if err != nil {
		zlog.Fatalf("could not open database: %s", err.Error())
		return nil, err
	}

	f, err := os.Open(initScriptPath)
	if err != nil {
		zlog.Fatalf("could not open database init file: %s", err.Error())
		return nil, err
	}

	if err := sqlite3.IntiDB(fileEventDB, f); err != nil {
		zlog.Fatalf("could not initialize the database: %s", err.Error())
		return nil, err
	}

	return fileEventDB, nil
}

func main() {

	// loading configuration
	config.LoadConfig()

	// TODO: configure different loggers for different services
	// initializing logger
	zlog.InitLogger(zlog.NewZapLogger(os.Stdout))

	// initializing database
	var fileEventDB, jobRunsDB database.Database
	{
		var err error = nil
		fileEventDB, err = dbSetup(
			config.ServiceConf.FileEventDB.DBPath,
			config.ServiceConf.FileEventDB.InitScript,
		)
		if err != nil {
			zlog.Fatalf("could not setup file event database: %s", err.Error())
		}

		jobRunsDB, err = dbSetup(
			config.ServiceConf.JobRunsDB.DBPath,
			config.ServiceConf.JobRunsDB.InitScript,
		)
		if err != nil {
			zlog.Fatalf("could not setup job runs database: %s", err.Error())
		}

		zlog.Infof("complted all database setups")
	}

	// starting the services
	{
		demon.Run(fileEventDB, jobRunsDB)
		restserver.Run(jobRunsDB)
	}
}
