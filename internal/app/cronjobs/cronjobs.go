package cronjobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/dto/cronjobs"
	"github.com/pramod-janardhana/dir-watcher/internal/helpers/jobs"
	serror "github.com/pramod-janardhana/dir-watcher/internal/server/error"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
	"github.com/pramod-janardhana/dir-watcher/repository/datastore"
	"github.com/pramod-janardhana/dir-watcher/repository/entity"
)

type Controller struct {
	jobRunsDB entity.JobRunsStore
}

func NewController(jobRunsDB database.Database) Controller {
	return Controller{
		jobRunsDB: datastore.NewJobRunsStore(jobRunsDB),
	}
}

func (c Controller) GetByStartTime(startsAtUTC time.Time) (*cronjobs.GetByStartTimeRes, error) {
	model := entity.JobRuns{}
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s='%s'",
		model.TableName(), model.ColumnStartTime(), startsAtUTC,
	)

	jobRuns, err := c.jobRunsDB.Query(query)
	if err != nil {
		zlog.Errorf("error fetching job run details: %s", err.Error())
		return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
	}

	if len(jobRuns) == 0 {
		zlog.Infof("not job runs for the given start time: %s", startsAtUTC)
		return nil, serror.NewError(http.StatusNotFound, fmt.Sprintf("not job runs for the given start time: %s", startsAtUTC))
	}

	jobrun := jobRuns[0]

	runDetails := entity.CronJobDetails{}
	if err := json.Unmarshal(jobrun.Details, &runDetails); err != nil {
		zlog.Errorf("error unmarshalling jobrun details: %v", err)
		zlog.Errorf("json: %s", jobrun.Details)
		return nil, serror.NewError(http.StatusInternalServerError, "error parsing response, please contact support")
	}

	toReturn := cronjobs.GetByStartTimeRes{
		JobRunDetails: cronjobs.JobRunDetails{
			RunId:        jobrun.RunId,
			JobType:      int(jobrun.JobType),
			StartTime:    jobrun.StartTime.UTC(),
			EndTime:      jobrun.EndTime.UTC(),
			TotalRunTime: fmt.Sprintf("%f min", jobrun.EndTime.Sub(jobrun.StartTime).Minutes()),
			Status:       int(jobrun.Status),
			Details: entity.CronJobDetails{
				MagicStringCount:  runDetails.MagicStringCount,
				DeletedFilesCount: runDetails.DeletedFilesCount,
				CreatedFilesCount: runDetails.CreatedFilesCount,
			},
		},
	}

	return &toReturn, nil
}

func (c Controller) UpdateConfig(req *cronjobs.UpdateConfigReq) (*cronjobs.UpdateConfigRes, error) {
	// verifying cron expression
	{
		scheduler := jobs.GetJobScheduler()
		if scheduler == nil {
			return nil, serror.NewError(http.StatusInternalServerError, "scheduler not available, please contact support")
		}

		testJob := &jobs.Job{
			CronExpresion: req.CronExpresion,
			Handler:       func() {},
		}

		if err := scheduler.AddJob(testJob); err != nil {
			zlog.Errorf("could not get add cron job with cron expression(%s): %s", req.CronExpresion, err.Error())
			return nil, serror.NewError(http.StatusBadRequest, "could not update cron job, please verify cron expression and try again")
		}

		scheduler.RemoveJob(testJob)
	}

	// updating the config file
	{
		path, err := config.JonSchedulerConfigPath()
		if err != nil {
			zlog.Errorf("could not get config.jobscheduler.windows.json file path: %v", err)
			return nil, serror.NewError(http.StatusInternalServerError, "could not find the config file, please contact support")
		}

		// read the config file
		data, err := ioutil.ReadFile(path)
		if err != nil {
			zlog.Errorf("error reading config file: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// unmarshal the JSON data into a struct
		var conf = config.CronJobScheduler{}
		if err := json.Unmarshal(data, &conf); err != nil {
			zlog.Errorf("error unmarshalling JSON data: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// updating the fields in the struct
		{
			if req.CronExpresion != "" {
				conf.CronExpression = req.CronExpresion
			}

			if req.MagicString != "" {
				conf.MagicString = req.MagicString
			}
		}

		// marshal the struct back to json
		updatedJSON, err := json.MarshalIndent(conf, "", "	")
		if err != nil {
			zlog.Errorf("error marshalling json data: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// writing the updated JSON data back to the file
		if err := ioutil.WriteFile(path, updatedJSON, 0644); err != nil {
			zlog.Errorf("error writing updated JSON to config file: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}
	}

	toReturn := cronjobs.UpdateConfigRes("restarting service with updated config")
	return &toReturn, nil
}
