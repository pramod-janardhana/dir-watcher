package process

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/internal/helpers/jobs"
	"github.com/pramod-janardhana/dir-watcher/internal/helpers/scanner"
	"github.com/pramod-janardhana/dir-watcher/internal/helpers/watcher"
	workerpool "github.com/pramod-janardhana/dir-watcher/internal/helpers/worker_pool"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
	"github.com/pramod-janardhana/dir-watcher/repository/datastore"
	"github.com/pramod-janardhana/dir-watcher/repository/entity"
)

func LoadScheduler(jobRunsDB, fileEventDB database.Database) {
	// stop the previous scheduler
	if js := jobs.GetJobScheduler(); js != nil {
		js.Stop()
	}

	fn := func() {
		jobRunsStore := datastore.NewJobRunsStore(jobRunsDB)
		jobRun := entity.JobRuns{
			RunId:     time.Now().UTC().String(),
			JobType:   entity.CRON_JOB,
			StartTime: time.Now().UTC(),
			EndTime:   time.Now().UTC(),
			Status:    entity.IN_PROGRESS,
			Details:   []byte(""),
		}

		// adding job run details to database
		if err := jobRunsStore.Upset(jobRun); err != nil {
			zlog.Errorf("error in cron job, failed to register the run details on job_run table: %s", err.Error())
			zlog.Debugf("jod run details: %+v", jobRun)
			return
		}

		// Configuring workers to process the request
		var wpool *workerpool.WorkerPool
		{
			wpool = workerpool.NewWorkerPool(10, func(task workerpool.Task) (any, error) {
				filepath := task.Details().(string)
				f, err := os.Open(filepath)
				if err != nil {
					return nil, err
				}

				return scanner.NewScanner(config.ServiceConf.CronJobScheduler.MagicString).Scan(f), nil
			})

			go wpool.Start()
		}

		var magicStringCount, deletedFilesCount, createdFilesCount int64 = 0, 0, 0
		fileEventStore := datastore.NewFileEventStore(fileEventDB)

		var abortSend, abortReceive = make(chan struct{}), make(chan struct{})
		var jobAborted atomic.Int32
		jobAborted.Store(0)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case feedback, ok := <-wpool.TaskFeedbackChan():
					if !ok {
						zlog.Infof("task feedback channel closed, stoping feedback prcessor")
						return
					}
					if !feedback.Success {
						zlog.Errorf("[feedback processor] worker failed to complete the task %d: %s", feedback.TaskId, feedback.Err.Error())
						zlog.Debugf("[feedback processor] sending abort signal to terminal the current job run")
						abortSend <- struct{}{}
						return
					} else {
						magicStringCount = magicStringCount + feedback.Data.(int64)
					}
				case <-abortReceive:
					zlog.Debugf("[feedback processor] received abort signal. Terminaing the current job run")
					jobAborted.Store(1)
					return
				}
			}
		}(&wg)

	beakLoop:
		for limit, offset := 20, 0; ; offset = offset + limit {
			select {
			case <-abortSend:
				zlog.Debugf("[task publisher] received abort signal. Terminaing the current job run")
				jobAborted.Store(1)
				break beakLoop
			default:
				fileEvents, err := fileEventStore.Get(limit, offset)
				if err != nil {
					zlog.Errorf("[task publisher] failed to get events from file event table: %s", err.Error())
					zlog.Debugf("[task publisher] sending abort signal to terminal the current job run")
					abortReceive <- struct{}{}

					break beakLoop
				}

				if len(fileEvents) == 0 {
					zlog.Debugf("completed processing")
					// done <- struct{}{}
					break beakLoop
				}

				for _, fileEvent := range fileEvents {
					switch fileEvent.Event {
					case int8(watcher.FILE_CREATED):
						createdFilesCount += 1
						wpool.TaskChan() <- workerpool.NewTask(rand.Int(), fileEvent.Path)
					case int8(watcher.FILE_DELETED):
						deletedFilesCount += 1
					case int8(watcher.FILE_MODIFIED):
						wpool.TaskChan() <- workerpool.NewTask(rand.Int(), fileEvent.Path)
					default:
						zlog.Errorf("unsupported event type %d for file %s", fileEvent.Event, fileEvent.Path)
					}
				}
			}
		}

		zlog.Infof("[task publisher] completed")

		wpool.Flush()
		zlog.Debugf("worker pool flushed")

		wg.Wait()

		jobRun.EndTime = time.Now().UTC()

		if jobAborted.Load() == 0 {
			zlog.Infof("current job run has completed")
			jobRun.Status = entity.COMPLETED
			jobRun.Details, _ = json.Marshal(entity.CronJobDetails{
				MagicStringCount:  magicStringCount,
				DeletedFilesCount: deletedFilesCount,
				CreatedFilesCount: createdFilesCount,
			})
		} else {
			zlog.Infof("current job run has failed")
			jobRun.Status = entity.FAILED
			jobRun.Details = bytes.NewBufferString("").Bytes()
		}

		if err := jobRunsStore.Upset(jobRun); err != nil {
			zlog.Errorf("error in cron job, failed to register the run details on job_run table: %s", err.Error())
			zlog.Debugf("jod run details: %+v", jobRun)
			return
		}

		// we don't have to mark the current run as failed but the next run might get wrong details.
		if err := fileEventStore.Truncate(); err != nil {
			zlog.Debugf("error truncating the file event table: %s", err.Error())
		}

		zlog.Debugf("job run with id %s completed", jobRun.RunId)
	}

	zlog.Debugf("starting cron jobs scheduler")
	scheduler := jobs.NewJobScheduler()
	scheduler.AddJob(&jobs.Job{
		CronExpresion: config.ServiceConf.CronJobScheduler.CronExpression,
		Handler:       fn,
	})

	scheduler.Start()
	zlog.Infof("started job scheduler")
}

func LoadFileWatcher(fileEventDB database.Database) {
	if fw := watcher.GetFileWatcher(); fw != nil {
		zlog.Debugf("stopping the previous file watcher")
		fw.Close()
	}

	if err := datastore.NewFileEventStore(fileEventDB).Truncate(); err != nil {
		zlog.Debugf("error truncating file event table before starting file watcher: %s", err.Error())
	}

	// file watcher configuration and registration
	zlog.Debugf("starting the file watcher")
	fw := watcher.NewFileWatcher(watcher.Config{
		PathToWatch: config.ServiceConf.FileWatcher.DirOrFileToWatch,
		Ops:         []watcher.Op{watcher.FILE_CREATED, watcher.FILE_MODIFIED, watcher.FILE_DELETED},
		Frequency:   time.Duration(config.ServiceConf.FileWatcher.FrequencyInSecond) * time.Second,
	})

	fwHandler := func(event watcher.Event) {
		fileEventStore := datastore.NewFileEventStore(fileEventDB)
		if !event.IsDir() {
			e := entity.FileEvent{Path: event.Path, Timestamp: event.ModTime().UTC()}
			switch watcher.Op(event.Op) {
			case watcher.FILE_CREATED:
				e.Event = int8(watcher.FILE_CREATED)
			case watcher.FILE_DELETED:
				e.Event = int8(watcher.FILE_DELETED)
			case watcher.FILE_MODIFIED:
				e.Event = int8(watcher.FILE_MODIFIED)
			default:
				zlog.Errorf("unsupported file watcher event %d for file %s", event.Op, event.Path)
				return
			}

			// TODO: add batching support
			if err := fileEventStore.Upset(e); err != nil {
				log.Println("error in file handle", err)
				zlog.Errorf("error adding event to file_event table: %s", err.Error())
				return
			}

			zlog.Debugf("add event (type: %d, file: %s) to file_event table", event.Op, event.Path)
		}
	}

	// register files being watched in database
	{
		for path, f := range fw.WatchedFiles() {
			e := watcher.NewEvent(watcher.FILE_CREATED, path, path, f)
			fwHandler(*e)
		}

		zlog.Debugf("registered all file being watched to file event database")
	}

	go func() {
		for {
			select {
			case event := <-fw.Event():
				fwHandler(watcher.Event{Event: event})
			case err := <-fw.Error():
				zlog.Errorf("errror in file watcher: %s", err)
			case <-fw.Closed():
				return
			}
		}
	}()

	go func() {
		// Start the watching process - it'll check for changes according to the polling frequency.
		if err := fw.Start(); err != nil {
			zlog.Errorf("failed to start file watcher: %s", err.Error())
		}
	}()
}
