package config

import (
	"log"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

type Service struct {
	General
	CronJobScheduler
	FileWatcher
}

type Server struct {
	Host           string
	Port           string
	AllowedOrigins []string
	AllowMethods   []string
	AllowHeaders   []string
}

type General struct {
	Server
	Logging
	FileEventDB Database
	JobRunsDB   Database
}

type Database struct {
	DBPath     string
	InitScript string
}

type Logging struct {
	Path string
}

type CronJobScheduler struct {
	CronExpression string `json:"cronExpression"`
	MagicString    string `json:"magicString"`
}

type FileWatcher struct {
	DirOrFileToWatch  string `json:"dirOrFileToWatch"`
	FrequencyInSecond int64  `json:"frequencyInSecond"`
}

// global configuration object
var cfg *koanf.Koanf = koanf.New(".")

var ServiceConf *Service = nil

func LoadGeneralConfig() *General {
	path, err := generalConfigPath()
	if err != nil {
		log.Fatalf("could not get config.general.windows.json file path: %v", err)
	}
	if err := cfg.Load(file.Provider(path), json.Parser()); err != nil {
		log.Fatalf("failed to load config.general.windows")
	}

	return &General{
		Server: Server{
			Host:           cfg.String("service.server.host"),
			Port:           cfg.String("service.server.port"),
			AllowedOrigins: cfg.Strings("service.server.allowedOrigins"),
			AllowMethods:   cfg.Strings("service.server.allowMethods"),
			AllowHeaders:   cfg.Strings("service.server.allowHeaders"),
		},
		FileEventDB: Database{
			DBPath:     os.ExpandEnv(cfg.String("service.fileEventDB.dbPath")),
			InitScript: os.ExpandEnv(cfg.String("service.fileEventDB.initScript")),
		},
		JobRunsDB: Database{
			DBPath:     os.ExpandEnv(cfg.String("service.jobRunsDB.dbPath")),
			InitScript: os.ExpandEnv(cfg.String("service.jobRunsDB.initScript")),
		},
		Logging: Logging{
			Path: cfg.String("service.logging.path"),
		},
	}
}

func LoadFileWatcherConfig() *FileWatcher {
	path, err := FileWatcherConfigPath()
	if err != nil {
		log.Fatalf("could not get config.filewatcher.windows.json file path: %v", err)
	}
	if err := cfg.Load(file.Provider(path), json.Parser()); err != nil {
		log.Fatalf("failed to load config.filewatcher.windows")
	}

	return &FileWatcher{
		DirOrFileToWatch:  cfg.String("dirOrFileToWatch"),
		FrequencyInSecond: cfg.Int64("frequencyInSecond"),
	}
}

func LoadJobSchedulerConfig() *CronJobScheduler {
	path, err := JonSchedulerConfigPath()
	if err != nil {
		log.Fatalf("could not get config.jobscheduler.windows.json file path: %v", err)
	}
	if err := cfg.Load(file.Provider(path), json.Parser()); err != nil {
		log.Fatalf("failed to load config.jobscheduler.windows")
	}

	return &CronJobScheduler{
		CronExpression: cfg.String("cronExpression"),
		MagicString:    cfg.String("magicString"),
	}
}

func LoadConfig() {
	ServiceConf = &Service{
		General:          *LoadGeneralConfig(),
		CronJobScheduler: *LoadJobSchedulerConfig(),
		FileWatcher:      *LoadFileWatcherConfig(),
	}
}
