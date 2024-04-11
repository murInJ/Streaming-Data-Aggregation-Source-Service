package services

import (
	config "SDAS/config"
	source "SDAS/services/source_service"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/cloudwego/kitex/pkg/klog"
)

func addPullExpose(name string, sourceName string, msgType string) error {
	entity, err := BuildExposeEntityPullByMsgType(name, sourceName, msgType)
	if err != nil {
		return err
	}
	entity.Start()
	for {
		if entity.GetStatus() == OPEN || entity.GetStatus() == PLAY || entity.GetStatus() == PAUSE {
			Exposes.Store(name, entity)
			return nil
		} else if entity.GetStatus() == ERR {
			err := errors.New("rtsp source start error")
			return err
		} else {
			runtime.Gosched()
		}
	}
}

func AddPullExpose(name string, sourceName string, msgType string) error {
	err := addPullExpose(name, sourceName, msgType)
	if err != nil {
		return err
	}
	refreshExpose()
	return nil
}

type ExposeEntityPull[T source.SOURCE_MESSAGE] struct {
	ControlChannel *chan int
	InputChannel   *chan T
	SeriesChannel  *chan string
	Status         int
	Name           string
	Type           string
	SourceName     string
	Expose         *config.EXPOSE
}

func BuildExposeEntityPullByMsgType(name string, sourceName string, TypeName string) (EXPOSE_ENTITY, error) {
	switch TypeName {
	case "rtspMsg":
		e, err := NewExposeEntityPull[source.MessageRtsp](name, sourceName)
		return e, err
	default:
		return nil, errors.New("not support type")
	}
}

func NewExposeEntityPull[T source.SOURCE_MESSAGE](name string, sourceName string) (*ExposeEntityPull[T], error) {
	control_channel := make(chan int)
	SeriesChannel := make(chan string, 1024)
	i, ok := source.Sources.Load(sourceName)
	if !ok {
		err := fmt.Errorf("source not found")
		klog.Error(err)
		return nil, err
	}
	v := reflect.ValueOf(i)
	entity := &ExposeEntityPull[T]{
		Name:           name,
		Type:           "pull",
		ControlChannel: &control_channel,
		InputChannel:   v.MethodByName("GetOutChannel").Call([]reflect.Value{})[0].Interface().(*chan T),
		Status:         CLOSE,
		SourceName:     sourceName,
		SeriesChannel:  &SeriesChannel,
	}
	return entity, nil

}

func (e *ExposeEntityPull[T]) GetExposeString() (string, error) {
	b, err := json.Marshal(e.Expose)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *ExposeEntityPull[T]) GetName() string {
	return e.Name
}

func (e *ExposeEntityPull[T]) GetType() string {
	return e.Type
}

func (e *ExposeEntityPull[T]) GetSourceName() string {
	return e.SourceName
}

func (e *ExposeEntityPull[T]) GetStatus() int {
	return e.Status
}

func (e *ExposeEntityPull[T]) GetControlChannel() *chan int {
	return e.ControlChannel
}

func (e *ExposeEntityPull[T]) GetData() (any, error) {
	data, ok := <-*e.SeriesChannel
	if !ok {
		return nil, errors.New("channel closed")
	} else {
		return data, nil
	}
}

func (e *ExposeEntityPull[T]) Start() {
	go e.goroutine_pull_expose()
}

func (e *ExposeEntityPull[T]) Stop() {
	*e.ControlChannel <- CLOSE
	close(*e.ControlChannel)
	close(*e.SeriesChannel)
}

func (e *ExposeEntityPull[T]) goroutine_pull_expose() {
	e.Status = OPEN
	klog.Infof("expose[pull]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case CLOSE:
				e.Status = CLOSE
				klog.Infof("expose[pull]: %s closed.", e.Name)
				return
			case PLAY:
				e.Status = PLAY
				klog.Infof("expose[pull]: %s play.", e.Name)
			case PAUSE:
				e.Status = PAUSE
				klog.Infof("expose[pull]: %s pause.", e.Name)
			}
		default:
		}

		if e.Status == PLAY {
			select {
			case msg := <-*e.InputChannel:
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err := enc.Encode(msg)
				if err != nil {
					klog.Errorf("expose[pull]: %s encode error: %v\n", e.Name, err)
					continue
				}

				*e.SeriesChannel <- buf.String()

			default:
			}
		}

		if e.Status == PAUSE {
			runtime.Gosched()
		}

	}
}

func (e *ExposeEntityPull[T]) Play() {
	*e.ControlChannel <- PLAY
}

func (e *ExposeEntityPull[T]) Pause() {
	*e.ControlChannel <- PAUSE
}
