package Router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scheduler/Worker"
)

func Run(port string) {

	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK,"%s","ok")
	})
	router.GET("/reload", func(context *gin.Context) {
		Worker.WorkerTag++//更新作业版本
		if Worker.WorkerRunning== false {//作业为停止状态时调起作业
			defer Worker.WorkerRun()
		}
		context.String(http.StatusOK,"%s","ok")
	})
	router.GET("/stop", func(context *gin.Context) {
		Worker.WorkerTag = 0//停止作业
		context.String(http.StatusOK,"%s","ok")
	})
	router.Run(":" + port)
}
