package services

import (
	config "SDAS/config"
	"encoding/json"
	"errors"
	"image"
	"reflect"
	"runtime"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	clone "github.com/huandu/go-clone/generic"
)

var Sources sync.Map

type SOURCE_ENTITY[T SOURCE_MESSAGE] interface {
	Start()
	Stop()
	GetSourceString() (string, error)
	GetName() string
	GetType() string
	GetOutChannel() *chan T
	send_to_out_channel(msg T)
}

type SOURCE_MESSAGE interface {
	GetNTP() int64
	GetImage() (image.Image, bool)
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
	configSources := []config.SOURCE{}
	Sources.Range(func(key, value interface{}) bool {
		v := reflect.ValueOf(value)

		res := v.MethodByName("GetSourceString").Call([]reflect.Value{})
		sourceString := res[0].Interface().(string)

		source := config.SOURCE{
			Name:    v.MethodByName("GetName").Call([]reflect.Value{})[0].Interface().(string),
			Type:    v.MethodByName("GetType").Call([]reflect.Value{})[0].Interface().(string),
			Content: sourceString,
		}
		if source.Type == "rtsp" {
			NewSources = append(NewSources, source)
		}

		return true
	})
	cp_conf := clone.Clone(config.Conf)
	cp_conf.Sources = configSources
	config.SaveConfigJSON("./config.json", cp_conf)
	config.Conf.Sources = NewSources
}

func AddRtspSource(source_config *config.SOURCE) error {
	err := addRtspSource(source_config)
	if err != nil {
		return err
	}
	refreshSource()
	return nil
}

func RemoveSource(name string) {
	if i, ok := Sources.Load(name); ok {
		v := reflect.ValueOf(i)
		v.MethodByName("Stop").Call([]reflect.Value{})
		Sources.Delete(name)
		refreshSource()
	}
}

func ListSources() []config.SOURCE {
	return config.Conf.Sources
}
