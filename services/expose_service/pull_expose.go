package services

import (
	config "SDAS/config"
	"SDAS/kitex_gen/api"
	source "SDAS/services/source_service"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
)

/**content
Op
**/

type ExposeEntityPull struct {
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

func NewExposeEntityPull(name, sourceName string, content map[string]string) (*ExposeEntityPull, error) {
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
	entity := &ExposeEntityPull{
		Name:           name,
		Type:           "pull",
		ControlChannel: &control_channel,
		SourceChannel:  c,
		Status:         config.CLOSE,
		SourceName:     sourceName,
		Content:        content,
	}
	return entity, nil

}

func (e *ExposeEntityPull) GetConfig() *api.Expose {
	return &api.Expose{
		Name:       e.Name,
		Type:       e.Type,
		SourceName: e.SourceName,
		Content:    e.Content,
	}
}

func (e *ExposeEntityPull) Start() error {
	if e.Status != config.CLOSE && e.Status != config.ERR {
		return errors.New("source already started")
	}
	go e.goroutinePullExpose()
	for {
		switch e.Status {
		case config.OPEN:
			return nil
		case config.ERR:
			err := errors.New("pull expose start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func (e *ExposeEntityPull) Stop() {
	e.once.Do(func() {
		*e.ControlChannel <- config.CLOSE
		close(*e.ControlChannel)
		v, ok := source.Sources.Load(e.SourceName)
		if ok {
			v.(source.SourceEntity).ReleaseOutChannel()
		}
	})
}

func (e *ExposeEntityPull) goroutinePullExpose() {
	defer func() {
		e.Wg.Done()
	}()
	e.Status = config.OPEN
	klog.Infof("expose[pull]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case config.CLOSE:
				e.Status = config.CLOSE
				klog.Infof("expose[pull]: %s closed.", e.Name)
				return
			case config.PLAY:
				e.Status = config.PLAY
				klog.Infof("expose[pull]: %s play.", e.Name)
			case config.PAUSE:
				e.Status = config.PAUSE
				klog.Infof("expose[pull]: %s pause.", e.Name)
			}
		default:
		}

		if e.Status == config.PLAY {
			select {
			case msg := <-*e.SourceChannel:
				resp := &api.PullExposeStreamResponse{
					Code:      0,
					Message:   "data",
					SourceMsg: msg,
				}
				if sendErr := e.Stream.Send(resp); sendErr != nil {
					e.Status = config.ERR
					klog.Error(sendErr)
					return
				}
			default:
			}
		}

		if e.Status == config.PAUSE {
			runtime.Gosched()
		}

	}
}
