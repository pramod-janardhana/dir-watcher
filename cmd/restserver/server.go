package restserver

import (
	"time"

	gcors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/internal/server/cors"
)

func Run(jobRunsDB database.Database) {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	cors.Add(engine, &gcors.Config{
		AllowMethods:     config.ServiceConf.Server.AllowMethods,
		AllowHeaders:     config.ServiceConf.Server.AllowHeaders,
		AllowOrigins:     config.ServiceConf.Server.AllowedOrigins,
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})

	registerRoutes(engine, jobRunsDB)

	// must be last line
	engine.Run(":" + config.ServiceConf.Server.Port)
}
