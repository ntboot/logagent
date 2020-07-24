package parseConfig

import (
	"context"

	"fmt"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
	"go.etcd.io/etcd/clientv3"

)
var ConfigMonitor = make(chan string, 10)

func Dirs(depth int)string{
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	parentdir := strings.Replace(dir, "\\", "/", -1)

	newdir := strings.Split(parentdir, "/")[:len(strings.Split(parentdir, "/")) - (depth -1 )]
	return strings.Join(newdir, "/")

}
//用config.ini 解析配置
func Parse(configpath string, configType string)(interface{}, error){
	cfg, err := ini.Load(configpath)
	if err!=nil{
		return nil, fmt.Errorf("")
	}
	return cfg.Section(configType), nil

}

//用etcd解析配置
func GetEtcdConfig(cli *clientv3.Client)(configinfo string){
	Getter, err := cli.Get(context.Background(), "logagentInfo")
	if err != nil{
		return
	}
	for _, item := range Getter.Kvs{
		configinfo = string(item.Value)

	}

	return
}

//实例化一个watcher, 热加载etcd中数据
func EtcdConfigWather(cli *clientv3.Client){
	watchChan := cli.Watch(context.Background(), "logagentInfo")
	for wc := range watchChan{
		for _, events := range wc.Events{
			ConfigMonitor <- string(events.Kv.Value)
		}
	}
}