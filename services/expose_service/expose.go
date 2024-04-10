package services

import (
	config "SDAS/config"
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
)

var Exposes sync.Map

type EXPOSE_ENTITY interface {
	Start()
	Stop()
	GetExposeString() (string, error)
	GetName() string
	GetType() string
	GetSourceName() string
}

func InitExpose() {
	for _, expose_config := range config.Conf.Exposes {
		switch expose_config.Type {
		case "rtsp":
			err := addRtspExpose(&expose_config)
			if err != nil {
				klog.Error("Init expose error: ", err)

			}
		}
	}

}

func addRtspExpose(expose_config *config.EXPOSE) error {
	name := expose_config.Name
	var expose config.EXPOSE_RTSP
	json.Unmarshal([]byte(expose_config.Content), &expose)
	entity, err := NewExposeEntityRtsp(name, &expose, expose_config.SourceName)
	if err != nil {
		return err
	}
	entity.Start()
	for {
		switch entity.Status {
		case OPEN:
			Exposes.Store(name, entity)
			return nil
		case ERR:
			err := errors.New("rtsp source start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func refreshExpose() {
	NewExposes := []config.EXPOSE{}

	Exposes.Range(func(key, value interface{}) bool {
		v := reflect.ValueOf(value)

		res := v.MethodByName("GetExposeString").Call([]reflect.Value{})
		exposeString := res[0].Interface().(string)

		expose := config.EXPOSE{
			Name:       v.MethodByName("GetName").Call([]reflect.Value{})[0].Interface().(string),
			Type:       v.MethodByName("GetType").Call([]reflect.Value{})[0].Interface().(string),
			Content:    exposeString,
			SourceName: v.MethodByName("GetSourceName").Call([]reflect.Value{})[0].Interface().(string),
		}
		NewExposes = append(NewExposes, expose)
		return true
	})

	config.Conf.Exposes = NewExposes
}

func AddRtspExpose(expose_config *config.EXPOSE) error {
	err := addRtspExpose(expose_config)
	if err != nil {
		return err
	}
	refreshExpose()
	config.SaveConfigJSON("./config.json")
	return nil
}

func RemoveExpose(name string) {
	if i, ok := Exposes.Load(name); ok {
		v := reflect.ValueOf(i)
		v.MethodByName("Stop").Call([]reflect.Value{})
		Exposes.Delete(name)
		refreshExpose()
		config.SaveConfigJSON("./config.json")
	}
}

func ListExposes() []config.EXPOSE {
	return config.Conf.Exposes
}
