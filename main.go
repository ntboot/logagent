package main

import (
	"Lightclass/logagent/etcd"
	"Lightclass/logagent/kafka"
	"Lightclass/logagent/parseConfig"
	"Lightclass/logagent/tail"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	tailf "github.com/hpcloud/tail"
	"go.etcd.io/etcd/clientv3"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//定义etcd 结构体
type Etcd struct {
	host string
	port int32
	configFile string
}

//定义获取etcd logagent配置接口
func (e *Etcd)GetEtcdConfig()(string, error){
	cli, err := etcd.InitEtcdClient(fmt.Sprintf("%s:%d", e.host, e.port ))
	if err!=nil{
		return "", fmt.Errorf("初始化etcd client失败...")
	}

	configinfo := parseConfig.GetEtcdConfig(cli)
	return configinfo, nil
}


func dirs(file interface{}, depth int)string{
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	parentdir := strings.Replace(dir, "\\", "/", -1)
	newdir := strings.Split(parentdir, "/")[:len(strings.Split(parentdir, "/")) - (depth -1 )]
	return strings.Join(newdir, "/")

}

var (
	msg = &sarama.ProducerMessage{}
	ConfigInfo = &etcd.Config{}
	LogTailChannel chan *tailf.Line
	lastTailList []chan *tailf.Line
	wg sync.WaitGroup
)

//动态异步从etcd中获取多个TailTask
func TailLogToKafka(ctx context.Context,  logStream chan *tailf.Line, topic string){
	client,  _ := kafka.InitKafka()
	defer func() {
		wg.Done()
		return
	}()
	defer client.Close()

	for line := range logStream {
		if line.Text == "quiting..."{
			//fmt.Printf("%s准备退出\n", topic)
			break
		}

		msg.Value = sarama.StringEncoder(line.Text)
		msg.Topic = topic
		fmt.Printf("%s, %s", msg.Topic, msg.Value)
		_,_,err:=client.SendMessage(msg)
		if err !=nil{
			fmt.Println(err.Error())
		}

	}
}


func main(){
	var etcdobj = Etcd{
		host:       "10.10.30.112",
		port:       2379,
		configFile: "",
	}
	cli, _:= etcd.InitEtcdClient(fmt.Sprintf("%s:%d", etcdobj.host, etcdobj.port))
	var ctx, _ = context.WithCancel(context.Background())

	//默认读取etcd中配置，启动tailTask
	initConfig := parseConfig.GetEtcdConfig(cli)
	json.Unmarshal([]byte(initConfig), ConfigInfo)

	for _, item := range ConfigInfo.Info {
		LogTailChannel = tail.InitTail(item.Logfile)
		lastTailList = append(lastTailList, LogTailChannel)
		wg.Add(1)
		go TailLogToKafka(ctx, LogTailChannel, item.Topic)
	}

	go func(cli *clientv3.Client){
		//启动一个watcher 配置的groutine
		go parseConfig.EtcdConfigWather(cli)

		for {
			select {
				//加载etcd watcher通道数据
				case test:= <- parseConfig.ConfigMonitor:

					if len(lastTailList) != 0{
						for _,singeTail := range lastTailList{
							//发送给tail 通道quiting...信息，让TaskTail的groutine退出
							singeTail <- &tailf.Line{Text: "quiting..."}
						}
					}

					json.Unmarshal([]byte(test), ConfigInfo)

					//这里先将tail 通道切片清空，为了下面append 新的配置的tail通道
					lastTailList = lastTailList[0:0]
					for _, item := range ConfigInfo.Info {
						//循环比较哪个item发生改变了
						LogTailChannel = tail.InitTail(item.Logfile)
						lastTailList = append(lastTailList, LogTailChannel)
						wg.Add(1)
						go TailLogToKafka(ctx, LogTailChannel, item.Topic)
					}
				default:
					break
			}
		}
	}(cli)
	wg.Wait()
}
