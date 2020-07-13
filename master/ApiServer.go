package master

import (
	"encoding/json"
	"github.com/Frank-liang/crond/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务的接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

//保存任务接口
func handleJobSave(resp http.ResponseWriter, req *http.Request) {

	var (
		err error
		postJob string
		job common.Job
		oldJob *common.Job
		bytes []byte
	)
	//解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
    //获取表中的job字段
	postJob = req.PostForm.Get("job")

	//反序列化
	if err = json.Unmarshal([]byte(postJob),&job); err != nil {
	    goto ERR
	}

	//保存到etcd
	if oldJob,err = G_jobMgr.SaveJob(&job); err != nil{
		goto ERR
	}
	//返回正常应答
	if bytes, err = common.BuildResponse(0,"success",oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	//返回错误应答
	if bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil {
		resp.Write(bytes)
	}
}

//删除任务接口
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		name string
		oldJob *common.Job
		bytes []byte

	)

	//POST: a=2&b=2&c=3
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//delete job name
	name = req.PostForm.Get("name")

	//delete job
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	//response OK
	if bytes, err = common.BuildResponse(0,"success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if  bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil {
		resp.Write(bytes)
	}

}

//list all cron jobs
func handleJobList(resp http.ResponseWriter, req *http.Request){
	var (
		jobList []*common.Job
		err error
		bytes []byte
	)
	if jobList, err = G_jobMgr.ListJobs(); err != nil{
		goto ERR
	}

	//response OK
	if bytes, err = common.BuildResponse(0,"success", jobList); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if  bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil {
		resp.Write(bytes)
	}

}

// kill job forced
func handleJobKill(resp http.ResponseWriter, req *http.Request){
	var (
		err error
		name string
		bytes []byte
	)

	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("name")

	if err = G_jobMgr.killjob(name); err != nil {
		goto  ERR

	}
	//response OK
	if bytes, err = common.BuildResponse(0,"success", nil); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if  bytes, err = common.BuildResponse(-1,err.Error(),nil); err == nil {
		resp.Write(bytes)
	}



}

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list",handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		ReadTimeout:   time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:  time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	return
}
