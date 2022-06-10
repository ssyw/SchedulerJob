package main

import (
	"scheduler/Router"
	"scheduler/Worker"
)



func main()  {
	go Worker.WorkerRun() //计划作业
	Router.Run("8000")//http监听
}

