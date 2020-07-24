package kafka
import (
	"Lightclass/logagent/parseConfig"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
)


func dirs(file interface{}, depth int)string{
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	parentdir := strings.Replace(dir, "\\", "/", -1)

	newdir := strings.Split(parentdir, "/")[:len(strings.Split(parentdir, "/")) - (depth -1 )]
	return strings.Join(newdir, "/")

}
func InitKafka()(sarama.SyncProducer ,error){
	// 实例化配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//获取kafka broker地址
	configinfo, err := parseConfig.Parse(filepath.Join(parseConfig.Dirs(1), "config.ini"), "kafka")
	if err != nil{
		return nil, fmt.Errorf("配置解析出错")
	}

	kafkahost, _ := configinfo.(*ini.Section).GetKey("host")
	kafkaport, _ := configinfo.(*ini.Section).GetKey("port")

	client, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%s", kafkahost.Value(), kafkaport.Value())}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return nil, fmt.Errorf("实例化client出错")
	}

	//构造一个消息体
	// 构造一个消息

	return client, nil






}