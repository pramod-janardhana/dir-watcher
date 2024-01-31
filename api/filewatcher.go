package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	fileWatcherDtos "github.com/pramod-janardhana/dir-watcher/dto/filewatcher"
	"github.com/pramod-janardhana/dir-watcher/internal/app/filewatcher"
	serror "github.com/pramod-janardhana/dir-watcher/internal/server/error"
	"github.com/pramod-janardhana/dir-watcher/internal/server/interceptor"
)

type fileWatcher struct {
	controller filewatcher.Controller
}

func NewFileWatcher() *fileWatcher {
	return &fileWatcher{
		controller: filewatcher.NewController(),
	}
}

func (c fileWatcher) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req *fileWatcherDtos.UpdateConfigReq = &fileWatcherDtos.UpdateConfigReq{}
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
