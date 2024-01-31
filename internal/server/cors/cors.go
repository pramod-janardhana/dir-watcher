package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

func Add(engine *gin.Engine, config *cors.Config) {
	if config == nil {
		zlog.Infof("no cors config was provided")
		return
	}

	corsObject := cors.New(*config)
	engine.Use(corsObject)

	zlog.Debugf("add cord config")
}
