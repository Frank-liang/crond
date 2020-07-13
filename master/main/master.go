package main

import (
	"flag"
	"fmt"
	"github.com/Frank-liang/crond/master"
	"runtime"
	"time"
)

var (
	confFile string
)


//解析命令行参数
func initArgs(){
	flag.StringVar(&confFile,"config","./master.json","指定master json")
	flag.Parse()
}


func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	var (
		err error
	)

	initArgs()

	//初始化线程
	initEnv()


	//加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

    //任务管理器
    if err = master.InitJobMgr(); err != nil {
    	goto ERR
	}


	//启动API http服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}


ERR:
	fmt.Println(err)

}
