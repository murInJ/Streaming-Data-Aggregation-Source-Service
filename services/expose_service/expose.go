package services

import (
	config "SDAS/config"
	"reflect"
	"sync"

	clone "github.com/huandu/go-clone/generic"
)

var (
	OPEN  = 0
	PLAY  = 1
	PAUSE = 2
	CLOSE = 3
	ERR   = 4
)
var Exposes sync.Map

type EXPOSE_ENTITY interface {
	Start()
	Stop()
	GetExposeString() (string, error)
	GetName() string
	GetType() string
	GetSourceName() string
	GetStatus() int
	GetControlChannel() *chan int
	GetData() (any, error)
}

func InitExpose() {
	for _, expose_config := range config.Conf.Exposes {
		switch expose_config.Type {
		// case "rtsp":
		// 	err := addRtspExpose(&expose_config)
		// 	if err != nil {
		// 		klog.Error("Init expose error: ", err)

		// 	}
		}
	}

}

func refreshExpose() {
	NewExposes := []config.EXPOSE{}
	configExposes := []config.EXPOSE{}

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
		if expose.Type == "rtsp" {
			configExposes = append(configExposes, expose)
		}
		return true
	})
	cp_conf := clone.Clone(config.Conf)
	cp_conf.Exposes = configExposes
	config.SaveConfigJSON("./config.json", cp_conf)
	config.Conf.Exposes = NewExposes
}

func RemoveExpose(name string) {
	if i, ok := Exposes.Load(name); ok {
		Exposes.Delete(name)
		v := reflect.ValueOf(i)
		v.MethodByName("Stop").Call([]reflect.Value{})
		refreshExpose()
	}
}

func ListExposes() []config.EXPOSE {
	return config.Conf.Exposes
}
