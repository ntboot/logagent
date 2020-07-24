package tail
import (
	"github.com/hpcloud/tail"
)

func InitTail(logpath string)chan *tail.Line{
	t, err := tail.TailFile(
		logpath, tail.Config{
		Follow: true,
		ReOpen:true,
		MustExist:false,
		Poll:true,
		Location:&tail.SeekInfo{Offset:0, Whence:2},
		})
	if err !=nil{
		return nil
	}

	return t.Lines
}

