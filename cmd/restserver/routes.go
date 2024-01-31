package restserver

import (
	"github.com/gin-gonic/gin"
	"github.com/pramod-janardhana/dir-watcher/api"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

func registerRoutes(engine *gin.Engine, jobRunsDB database.Database) {
	zlog.Debugf("starting route registration for the REST server")

	baseRoute := engine.Group("/api/v1")
	{
		baseRoute.HEAD("/health", api.Health)
		baseRoute.GET("/health", api.Health)

		cronJobRoutes := baseRoute.Group("/cronjob")
		{
			cronJobAPIs := api.NewCronJobs(jobRunsDB)

			cronJobRoutes.GET("", cronJobAPIs.Get())
			cronJobConfigRoutes := cronJobRoutes.Group("/config")
			{
				cronJobConfigRoutes.PUT("", cronJobAPIs.Update())
			}
		}

		fileWatcherRoutes := baseRoute.Group("/filewatcher")
		{
			fileWatcherAPIs := api.NewFileWatcher()
			fileWatcherConfigRoutes := fileWatcherRoutes.Group("/config")
			{
				fileWatcherConfigRoutes.PUT("", fileWatcherAPIs.Update())
			}
		}
	}
}
