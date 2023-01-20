package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/yum2npm/yum2npm/pkg/handlers"
)

func setupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(gin.Recovery())

	r.GET("/", handlers.IndexHandler(config.Repos))
	r.GET("/repos", handlers.GetRepos(config.Repos))
	r.GET("/repos/:repo/packages", handlers.GetPackages(&repodata))
	r.GET("/repos/:repo/packages/:package", handlers.GetPackage(&repodata))
	r.GET("/repos/:repo/modules", handlers.GetModules(&modules))
	r.GET("/repos/:repo/modules/:module/packages", handlers.GetModulePackages(&modules))
	r.GET("/repos/:repo/modules/:module/packages/:package", handlers.GetModulePackage(&modules))

	return r
}
