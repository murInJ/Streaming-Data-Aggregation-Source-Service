package services

import (
	"SDAS/config"
	"SDAS/kitex_gen/api"
	source "SDAS/services/source_service"
	"SDAS/utils"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	jsoniter "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack/v5"
	"image"
	"runtime"
	"sync"
)

/**content
url string
**/

type ExposeEntityHttpPush struct {
	ControlChannel *chan int
	SourceChannel  *chan *api.SourceMsg
	Status         int
	Name           string
	Type           string
	SourceName     string
	Content        map[string]string
	once           sync.Once
	Stream         api.SDAS_PullExposeStreamServer
	Wg             *sync.WaitGroup
}

func NewExposeEntityHttpPush(name, sourceName string, content map[string]string) (*ExposeEntityHttpPush, error) {
	control_channel := make(chan int)
	i, ok := source.Sources.Load(sourceName)
	if !ok {
		err := fmt.Errorf("source %s not found", sourceName)
		klog.Error(err)
		return nil, err
	}
	sourceEntity := i.(source.SourceEntity)
	if !sourceEntity.IsExpose() {
		return nil, fmt.Errorf("source is not expose")
	}
	c, err := sourceEntity.RequestOutChannel()
	if err != nil {
		return nil, err
	}
	entity := &ExposeEntityHttpPush{
		Name:           name,
		Type:           "httpPush",
		ControlChannel: &control_channel,
		SourceChannel:  c,
		Status:         config.CLOSE,
		SourceName:     sourceName,
		Content:        content,
	}
	return entity, nil

}

func (e *ExposeEntityHttpPush) GetConfig() *api.Expose {
	return &api.Expose{
		Name:       e.Name,
		Type:       e.Type,
		SourceName: e.SourceName,
		Content:    e.Content,
	}
}

func (e *ExposeEntityHttpPush) Start() error {
	if e.Status != config.CLOSE && e.Status != config.ERR {
		return errors.New("source already started")
	}
	go e.goroutineHttpPushExpose()
	for {
		switch e.Status {
		case config.OPEN:
			return nil
		case config.ERR:
			err := errors.New("httpPush expose start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func (e *ExposeEntityHttpPush) Stop() {
	e.once.Do(func() {
		*e.ControlChannel <- config.CLOSE
		close(*e.ControlChannel)
		v, ok := source.Sources.Load(e.SourceName)
		if ok {
			v.(source.SourceEntity).ReleaseOutChannel()
		}
	})
}

func (e *ExposeEntityHttpPush) goroutineHttpPushExpose() {
	e.Status = config.OPEN
	klog.Infof("expose[httpPush]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case config.CLOSE:
				e.Status = config.CLOSE
				klog.Infof("expose[httpPush]: %s closed.", e.Name)
				return
			}
		default:
		}

		select {
		case msg := <-*e.SourceChannel:
			var img image.RGBA
			err := msgpack.Unmarshal(msg.Data, &img)

			if err != nil {
				e.Status = config.ERR
				klog.Error(err)
				return
			}
			base64, err := utils.RGBAToBase64(&img)

			if err != nil {
				e.Status = config.ERR
				klog.Error(err)
				return
			}

			body := struct {
				DataType string
				Ntp      int64
				Data     string
			}{
				DataType: msg.DataType,
				Ntp:      msg.Ntp,
				Data:     base64,
			}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			marshal, err := json.Marshal(&body)
			if err != nil {
				e.Status = config.ERR
				klog.Error(err)
				return
			}
			// 创建一个HTTP请求
			err = utils.PostJson(e.Content["url"], marshal)
			if err != nil {
				e.Status = config.ERR
				klog.Error(err)
				return
			}

		default:
		}

	}
}
