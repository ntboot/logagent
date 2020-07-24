package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"go.etcd.io/etcd/clientv3"
	"time"

)
var (
	etcdAddress string = "10.10.30.112:2379"
)
type DetailConfig struct {
	Topic string `json:"topic"`
	Logfile string `json:"logfile"`
}

type Config struct {
	Info []DetailConfig `json:"info"`
}

//反射config.json类容
func reflex(configpath string)*Config{
	var info, _ = ioutil.ReadFile(configpath)
	var demo = &Config{}
	json.Unmarshal(info, demo)

	return demo
}

//实例化etcd客户端
func InitEtcdClient(etcdhost string)(*clientv3.Client, error){
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:            []string{etcdhost},
			DialTimeout:          5*time.Second,
		})

	if err!= nil{
		return nil, fmt.Errorf("etcd客户端实例化失败...")
	}

	return cli, nil
}

//推送配置
func PushConfig(cli *clientv3.Client, configfile string){
	cli, err := InitEtcdClient(etcdAddress)
	defer cli.Close()
	if err != nil{
		return
	}

	//set config to etcd
	var configInfo,_ = json.Marshal(reflex(configfile))
	cli.Put(context.Background(), "logagentInfo", string(configInfo))

}


