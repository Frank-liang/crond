package main

import "runtime"

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func mian() {

	//初始化线程
	initEnv()

}
