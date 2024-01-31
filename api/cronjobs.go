package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pramod-janardhana/dir-watcher/database"
	cronjobDtos "github.com/pramod-janardhana/dir-watcher/dto/cronjobs"
	"github.com/pramod-janardhana/dir-watcher/internal/app/cronjobs"
	serror "github.com/pramod-janardhana/dir-watcher/internal/server/error"
	"github.com/pramod-janardhana/dir-watcher/internal/server/interceptor"
)

type cronJobs struct {
	controller cronjobs.Controller
}

func NewCronJobs(jobRunsDB database.Database) *cronJobs {
	return &cronJobs{
		controller: cronjobs.NewController(jobRunsDB),
	}
}

func (c cronJobs) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startAt := ctx.Query("startAt")
		if startAt != "" {
			format := "2006-01-02 15:04:05 0000 UTC"
			t, err := time.Parse(format, startAt)
			if err != nil {
				interceptor.SendErrRes(
					ctx,
					fmt.Sprintf("startAt must be of the form '%s'", format),
					http.StatusBadRequest,
				)
				return
			}

			res, err := c.controller.GetByStartTime(t.UTC())
			if err != nil {
				serr, _ := err.(serror.Error)
				interceptor.SendErrRes(
					ctx,
					serr.ErrMessage,
					serr.StatusCode,
				)
				return
			}

			interceptor.SendSuccessRes(ctx, res, http.StatusOK)
			return
		}

		interceptor.SendErrRes(
			ctx,
			"include startAt timestamp to get job run details",
			http.StatusBadRequest,
		)
		return
	}
}

func (c cronJobs) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *cronjobDtos.UpdateConfigReq = &cronjobDtos.UpdateConfigReq{}
		if err := ctx.ShouldBind(req); err != nil {
			interceptor.SendErrRes(ctx, "could not parse the request, please verify request body and try again", http.StatusBadRequest)
			return
		}

		res, err := c.controller.UpdateConfig(req)
		if err != nil {
			serr, _ := err.(serror.Error)
			interceptor.SendErrRes(
				ctx,
				serr.ErrMessage,
				serr.StatusCode,
			)
			return
		}

		interceptor.SendSuccessRes(ctx, res, http.StatusOK)
		return
	}
}
