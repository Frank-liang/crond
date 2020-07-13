package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	EtcdEndPoint string `json:"etcdEndPoint"`
	EtcdDailTImeOut int `json:"etcdDailTimeOut"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error){
	var (
		content []byte
		conf Config
	)

	//读取json文件
	if content, err = ioutil.ReadFile(filename); err != nil{
		return
	}

	//解析json
	if err = json.Unmarshal(content,&conf); err != nil {
		return
	}

	G_config = &conf


	return

}

