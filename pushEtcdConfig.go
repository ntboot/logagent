// +build ignore

package main

import (
	"Lightclass/logagent/etcd"
	"fmt"
	"Lightclass/logagent/parseConfig"
	"path/filepath"
)

//定义etcd 结构体
type EtcdAddress struct {
	host string
	port int32
	configFile string
}

var etcdobj = EtcdAddress{
	host:       "10.10.30.112",
	port:       2379,
	configFile: "",
}

func main(){
	cli, _:= etcd.InitEtcdClient(fmt.Sprintf("%s:%d", etcdobj.host, etcdobj.port))
	configFilePath := filepath.Join(parseConfig.Dirs(1), "etcd/config.json")

	//推送最新的配置至etcd
	etcd.PushConfig(cli, configFilePath)

}
