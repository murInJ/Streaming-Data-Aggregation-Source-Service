package services

import (
	"SDAS/client"
	"SDAS/config"
	"SDAS/kitex_gen/api"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
)

/**content
	Url
	SourceName
**/

type SourceEntityRemote struct {
	ControlChannel *chan int
	OutputChannel  *chan *api.SourceMsg
	Status         int
	Content        map[string]string
	Name           string
	Type           string
	Expose         bool
	once           sync.Once
	requested      atomic.Bool
}

func NewSourceEntityRemote(name string, expose bool, content map[string]string) (*SourceEntityRemote, error) {

	controlChannel := make(chan int)
	outputChannel := make(chan *api.SourceMsg, 1024)

	entity := &SourceEntityRemote{
		Name:           name,
		Type:           "remote",
		ControlChannel: &controlChannel,
		OutputChannel:  &outputChannel,
		Status:         config.CLOSE,
		Content:        content,
		Expose:         expose,
	}
	entity.requested.Store(false)
	return entity, nil

}

func (e *SourceEntityRemote) GetConfig() *api.Source {
	return &api.Source{
		Type:    e.Type,
		Name:    e.Name,
		Expose:  e.Expose,
		Content: e.Content,
	}
}

func (e *SourceEntityRemote) RequestOutChannel() (*chan *api.SourceMsg, error) {
	if e.requested.CompareAndSwap(false, true) {
		return e.OutputChannel, nil
	} else {
		return nil, errors.New("request out channel already in use")
	}
}

func (e *SourceEntityRemote) ReleaseOutChannel() {
	e.requested.Store(false)
}

func (e *SourceEntityRemote) Start() error {
	if e.Status != config.CLOSE && e.Status != config.ERR {
		return errors.New("source already started")
	}
	go e.goroutineRemoteSource()
	for {
		switch e.Status {
		case config.OPEN:
			return nil
		case config.ERR:
			err := errors.New("remote source start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func (e *SourceEntityRemote) Stop() {
	e.once.Do(func() {
		*e.ControlChannel <- config.CLOSE
		close(*e.OutputChannel)
		close(*e.ControlChannel)
	})
}

func (e *SourceEntityRemote) GetName() string {
	return e.Name
}

func (e *SourceEntityRemote) IsExpose() bool {
	return e.Expose
}

func (e *SourceEntityRemote) goroutineRemoteSource() {
	c, err := e.startupRemote()
	if err != nil {
		e.Status = config.ERR
		klog.Errorf("source[remote]: %s open failed. %v", e.Name, err)
		return
	}
	e.Status = config.OPEN
	klog.Infof("source[remote]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case config.CLOSE:
				err := c.Close()
				if err != nil {
					klog.Error(err.Error())
				}
				e.Status = config.CLOSE
				klog.Infof("source[remote]: %s closed.\n", e.Name)
				return
			}
		default:
		}

		msg, err := c.RecvPullExposeStream()
		if err == io.EOF {
			continue
		} else if err != nil {
			klog.Error(err.Error())
			err := c.Close()
			if err != nil {
				klog.Error(err.Error())
			}
			e.Status = config.ERR
			klog.Infof("source[remote]: %s closed.\n", e.Name)
			return
		}
		if msg != nil {
			*e.OutputChannel <- msg
		}
	}
}

func (e *SourceEntityRemote) startupRemote() (*client.SDASClient, error) {
	c, err := client.NewSDASClient(e.Content["url"], true, true)
	if err != nil {
		return nil, err
	}
	err = c.SendPullExposeStream(e.Name, e.Content["source_name"], config.OPEN)
	if err != nil {
		return nil, err
	}
	err = c.SendPullExposeStream(e.Name, e.Content["source_name"], config.PLAY)
	if err != nil {
		return nil, err
	}
	return c, nil
}
