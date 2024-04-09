package services

import (
	config "SDAS/config"
	"encoding/json"
	"errors"
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
	entity := NewEntityRtsp(name, &source)
	entity.Start()
	for {
		switch entity.Status {
		case OPEN:
			Sources.Store(name, entity)
			return nil
		case ERR:
			err := errors.New("rtsp source start error")
			return err
		}
	}
}

func refreshSource() {
	NewSources := []config.SOURCE{}

	Sources.Range(func(key, value interface{}) bool {
		entity := *value.(*SOURCE_ENTITY)
		sourceString, err := entity.GetSourceString()
		if err != nil {
			klog.Error(err)
			return true
		}
		source := config.SOURCE{
			Name:    sourceString,
			Type:    entity.GetType(),
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
		entity := *i.(*SOURCE_ENTITY)
		entity.Stop()
		Sources.Delete(name)
		refreshSource()
		config.SaveConfigJSON("./config.json")
	}
}

func ListSources() []config.SOURCE {
	return config.Conf.Sources
}
