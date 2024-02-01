# Introduction
**Supported only on Windows**

DirWatcher is long running background time looped scheduled task which exposes the results via REST API.
This has two main components:
1. REST API Server
2. scheduled backuground task  

## Tech stack
1. Golang
2. Gin for API Server
3. Sqlite database

# Highlights design
![dri-watcher design](https://github.com/pramod-janardhana/dir-watcher/blob/main/artifacts/dir-watcher.svg)
## Components
### [1] Application starter
Application starter is responsible for starting the Process Handler and REST API Server.

### [2] Process Handler
Process Handler is responsible for:
1. Starting File Watcher
2. Starting Cron Job Scheduler
3. Registering watcher for file_watcher and cron_job_scheduler config file to notify the respective components.

#### [2.1] File Watcher
File Watcher is responsible for:
1. Listning for registered events on the configured watch directory.
2. Record the event in the file_event database.

#### [2.2] Cron Job Scheduler
Cron Job Scheduler is responsible for:
1. Starting Cron Job according to the configuration.
2. Record "start_time", "end_time", "total_run_time", "files_added", "files_removed", "magic_string_count" in the job_run database.

### [3] REST API Server
REST API Server is responsible for starting and configuruing the Gin HTTP Server.

# Database schema

## file event database
![file_event database schema](https://github.com/pramod-janardhana/dir-watcher/blob/main/artifacts/file_event_db.png)
### file_event table
1. id is the primary key
2. deleted_files_count is the number of files deleted after the previous cron job run.
3. created_files_count is the number of created files after the previous cron job run.

## job run database
![job_run database schema](https://github.com/pramod-janardhana/dir-watcher/blob/main/artifacts/job_runs_db.png)
### job_run table
1. run_id is the primary key
2. job_type is the type of job
   * 1 - cron job
3. start_time is start the time of the job run
4. end_time is the end time of the job run
5. status is the current status of the job run
   * 1 - in-progress
   * 2 - failed
   * 3 - completed
6. details is the job_type specific details in json format
   * job_type 1: magicStringCount, deletedFilesCount, createdFilesCount


# Project setup
## Configuration files

### [1] general service config file
NOTE: This file should named as "config.general.windows.json" and should be present in ${programdata}\DirWatcher\ directory.
```
{
  "service": {
    "name": "DirWatcher",
    "version": "v1.0",
    "fileEventDB":{
      "dbPath": "<path to file_event.sqlite>",
      "initScript": "<path to file_event.sql>"
    },
    "jobRunsDB":{
      "dbPath": "<path to job_runs.sqlite>",
      "initScript": "<path to job_runs.sql>"
    },
    "server": {
      "port": "8000",
      "allowedOrigins": ["*"],
      "allowedMethods": ["GET", "PUT"],
      "allowedHeaders": ["Authorization", "Content-Type", "Origin", "ResponseType"]
    },
    "logging":{
      "path":"dirwatcher.log"
    }
  }
}
```

**service.name** is the name of the service.

**service.version** is the version of the service.

**service.fileEventDB.dpPath** is the path to the file_event.sqlite file (if the .sqlite file does not exist then the application will create it in the provided path).
**service.fileEventDB.initScript** is the path to the file_event.sql file.

**service.jobRunsDB.dpPath** is the path to the job_runs.sqlite file (if the .sqlite file does not exist then the application will create it in the provided path),
**service.jobRunsDB.initScript** is the path to the job_runs.sql file.

**service.server.port** is the port of the REST API server.
**service.server.allowedOrigins** are the list of allowed origins.
**service.server.allowedMethods** are the list of allowed methods.
**service.server.allowedHeaders** is the list of allowed headers.

**service.logging.path** is the path to the log file.

### file watcher service config file
NOTE: This file should named as "config.filewatcher.windows.json" and should be present in ${programdata}\DirWatcher\ directory.
```
{
	"dirOrFileToWatch": "<path to file or folder to be watched>",
	"frequencyInSecond": <pooling ferquency to check for file events in seconds>
}
```

**dirOrFileToWatch** is path to file or folder to be watched.

**frequencyInSecond** is the pooling ferquency to check for file events in seconds

### file watcher service config file
NOTE: This file should named as "config.jobscheduler.windows.json" and should be present in ${programdata}\DirWatcher\ directory.
```
{
	"cronExpression": "<valid cron expression>",
	"magicString": "<magic string>"
}
```

**cronExpression** is the cron expression (provide the expression in such a way that not more then 1 cron job is running at a given point of time to ensure consistent and accurate results. size should provided considring the size and number of files being watched).

**magicString** is the magic string to be scanned in the files being watched.

## Running the project
Run the command from the root of the project
```go run cmd/main.go```

# API collection
[Postman collection](https://github.com/pramod-janardhana/dir-watcher/blob/main/artifacts/dir-watcher%20apis.postman_collection.json)

# Limitations
1. Supports only windows OS.
2. Only add and remove file events are supported. File events like rename, move etc are not supported.
3. Does not perform strict validation on config files

# TODO:
1. Add autentication and authorization to REST APIs.
2. Provide APIs:
    * to get more details job run.
    * manually start and stop jobs.