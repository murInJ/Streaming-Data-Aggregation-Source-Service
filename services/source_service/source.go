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

var Sources sync.Map

type SOURCE_ENTITY interface {
	Start()
	Stop()
	GetSourceString() (string, error)
	GetName() string
	GetType() string
}

func InitSource() {
	for _, source_config := range config.Conf.Sources {
		switch source_config.Type {
		case "rtsp":
			err := addRtspSource(&source_config)
			if err != nil {
				klog.Error("Init source error: ", err)

			}
		}
	}

}

func addRtspSource(source_config *config.SOURCE) error {
	name := source_config.Name
	var source config.SOURCE_RTSP
	json.Unmarshal([]byte(source_config.Content), &source)
	entity, err := NewSourceEntityRtsp(name, &source)
	if err != nil {
		return err
	}
	entity.Start()
	for {
		switch entity.Status {
		case OPEN:
			Sources.Store(name, entity)
			return nil
		case ERR:
			err := errors.New("rtsp source start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func refreshSource() {
	NewSources := []config.SOURCE{}

	Sources.Range(func(key, value interface{}) bool {
		v := reflect.ValueOf(value)

		res := v.MethodByName("GetSourceString").Call([]reflect.Value{})
		sourceString := res[0].Interface().(string)

		source := config.SOURCE{
			Name:    v.MethodByName("GetName").Call([]reflect.Value{})[0].Interface().(string),
			Type:    v.MethodByName("GetType").Call([]reflect.Value{})[0].Interface().(string),
			Content: sourceString,
		}
		NewSources = append(NewSources, source)
		return true
	})

	config.Conf.Sources = NewSources
}

func AddRtspSource(source_config *config.SOURCE) error {
	err := addRtspSource(source_config)
	if err != nil {
		return err
	}
	refreshSource()
	config.SaveConfigJSON("./config.json")
	return nil
}

func RemoveSource(name string) {
	if i, ok := Sources.Load(name); ok {
		v := reflect.ValueOf(i)
		v.MethodByName("Stop").Call([]reflect.Value{})
		Sources.Delete(name)
		refreshSource()
		config.SaveConfigJSON("./config.json")
	}
}

func ListSources() []config.SOURCE {
	return config.Conf.Sources
}
