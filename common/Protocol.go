package common

import (
	"encoding/json"
	"fmt"
)

//定时任务
type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//http接口应答
type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

//应答方法
func BuildResponse(errno int, msg string,data interface{})(resp []byte, err error){
	var (
	response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	fmt.Println(response)

	//序列化json
	resp, err = json.Marshal(response)
	return

}