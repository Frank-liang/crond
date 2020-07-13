package master

import (
	"context"
	"encoding/json"
	"github.com/Frank-liang/crond/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	//单类
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints: []string{G_config.EtcdEndPoint},
		DialTimeout: time.Duration(G_config.EtcdDailTImeOut) * time.Millisecond,
	}

	//建立链接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
		lease: lease,
	}

	return
}

//保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job)(oldJob *common.Job, err error){
	//把任务保存到/cron/job/任务名 ->json
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job
	)

	//etcd保存key
	jobKey = common.JOB_SAVE_DIR + job.Name
	//任务信息json
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	//保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV()); err != nil {
		return
	}

	//如果是更新,返回旧值
	if putResp.PrevKv != nil {
		//对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value,&oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

//删除任务
func(jobMgr *JobMgr) DeleteJob(name string)(oldJob *common.Job, err error){
	var(
		jobKey string
		delResp *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	jobKey = common.JOB_SAVE_DIR + name

	if delResp,err = jobMgr.kv.Delete(context.TODO(),jobKey,clientv3.WithPrevKV()); err != nil{
		return
	}

    //返回被删除的任务信息
    if len(delResp.PrevKvs) != 0 {
    	if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
    		err = nil
    		return
		}
		oldJob = &oldJobObj
	}

	return
}

func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error){
	var (
		dirKey string
		getResp *clientv3.GetResponse
		kvPairs *mvccpb.KeyValue
		job *common.Job
	)

	dirKey = common.JOB_SAVE_DIR

	if getResp, err = jobMgr.kv.Get(context.TODO(),dirKey,clientv3.WithPrevKV()); err != nil {
		return
	}

	jobList = make([]*common.Job,0)
	// get all keys and values
	for _, kvPairs = range  getResp.Kvs{
		job = &common.Job{}
		if err = json.Unmarshal(kvPairs.Value,job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList,job)

	}
return
}

func (jobMgr *JobMgr) killjob(name string)(err error) {
	var (
		killerkey string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
	)

     killerkey = common.JOB_KILL_DIR + name

     //让worker监听到一次put操作，创建一个租约让其稍后自动过期即可
     if leaseGrantResp, err =jobMgr.lease.Grant(context.TODO(),1); err != nil {
     	return
	 }

	 //租约ID
	 leaseId = leaseGrantResp.ID

	 //设置killer 标记
	 if _,err = jobMgr.kv.Put(context.TODO(),killerkey,"",clientv3.WithLease(leaseId)); err != nil{
	 	return
	 }

return
}