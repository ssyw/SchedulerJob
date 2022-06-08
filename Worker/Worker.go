package Worker

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Worker struct {
}

type Job struct {
	Interval int `yaml:"interval"`
	Url string	`yaml:"url"`
	RunTimes uint	`yaml:"run_times"`
	TimeOut int `yaml:"time_out"`
}


var WorkerTag = 1  //当前作业的版本
var WorkerRunning = true //作业状态
var taskMap = make(map[string]Job)//内存作业列表

func WorkerRun()  {
	wg := sync.WaitGroup{}
	for key,value := range taskMap{
		wg.Add(1)//协程+1
		go doWork(key,value,WorkerTag,&wg)
	}
	wg.Wait()//所有协程结束
	WorkerRunning = false//标记作业状态为停止
	if WorkerTag != 0{
		fmt.Println("restart")
		LoadConfig()//重新载入作业列表
		WorkerRunning = true//标记作业状态为开始
		WorkerRun()
	}
	fmt.Println("worker down")
}


func init()  {
	LoadConfig()
}

func LoadConfig(){
	content,err := ioutil.ReadFile("Config/cronConf.yaml")
	if err!= nil{
		panic("未获取到配置文件")
	}

	taskConfig := &taskMap
	err = yaml.Unmarshal(content,taskConfig)
	if err!= nil{
		panic("配置文件加载错误")
	}
}

func doWork(job string,content Job,tag int,wg *sync.WaitGroup)  {

	if content.Interval == 0 {
		if content.RunTimes == 0 {//一次性作业
			content.RunTimes ++
			//TODO 如一次性作业为唯一作业时，可把这个次数更新到Yaml文件里
			fmt.Printf("%s 任务：%s 运行参数：%s 运行间隔：%v 运行次数：%v\n",time.Now().Format("2006-01-02 15:04:05"),job,content.Url,content.Interval,content.RunTimes)
			fmt.Println("once stopped")
		}
		wg.Done()
	}else{
		timeTicker := time.NewTicker(time.Second*time.Duration(content.Interval))//开启计时器

		for range timeTicker.C{
			content.RunTimes++
			fmt.Printf("%s 任务名:%s 触发\n",time.Now().Format("2006-01-02 15:04:05"),job)
			ch := make(chan struct{},1);
			go func() {
				fmt.Printf("%s 任务名:%s 运行中\n",time.Now().Format("2006-01-02 15:04:05"),job)
				time.Sleep(time.Second*30)//模拟超时
				ch<- struct{}{}
			}()

			//超时控制
			select {
				case res :=<-ch:
					fmt.Printf("%s 任务:%s  运行间隔：%v 执行结果%v\n",time.Now().Format("2006-01-02 15:04:05"),job,content.Interval,res)
				case <-time.After(time.Duration(content.TimeOut)*time.Second):
					fmt.Printf("%s 任务:%s  运行间隔：%v 运行超时\n",time.Now().Format("2006-01-02 15:04:05"),job,content.Interval)
			}
			if tag != WorkerTag {//当前作业版本与内存作业版本不致时，重启整个作业
				timeTicker.Stop()//停止计时器
				wg.Done()//协程结束
				fmt.Println("stopped")
			}
		}
	}
}
