package services

import (
	"SDAS/kitex_gen/api"
	"errors"
	"sync"
)

var Sources sync.Map

type SourceEntity interface {
	Start() error
	Stop()
	GetConfig() *api.Source
	RequestOutChannel() (*chan *api.SourceMsg, error)
	ReleaseOutChannel()
	GetName() string
}

func BuildSourceEntity(Name, Type string, Expose bool, Content map[string]string) ([]SourceEntity, error) {
	switch Type {
	case "rtsp":
		entity, err := NewSourceEntityRtsp(Name, Expose, Content)
		if err != nil {
			return nil, err
		}
		return []SourceEntity{entity}, nil
	case "plugin":
		enities, err := NewSourceEntityPlugin(Name, Expose, Content)
		if err != nil {
			return nil, err
		}
		var sourceEntities []SourceEntity
		for _, entity := range enities {
			sourceEntities = append(sourceEntities, entity)
		}

		return sourceEntities, nil
	default:
		return nil, errors.New("type" + Name + " is not supported")
	}
}

func AddSource(Name, Type string, Expose bool, Content map[string]string) error {
	if _, ok := Sources.Load(Name); ok {
		return errors.New("Source already exists")
	}
	entities, err := BuildSourceEntity(Name, Type, Expose, Content)
	if err != nil {
		return err
	}
	for _, entity := range entities {
		err = entity.Start()
		if err != nil {
			return err
		}
		Sources.Store(entity.GetName(), entity)
	}

	return nil
}

func RemoveSource(name string) {
	if i, ok := Sources.Load(name); ok {
		entity := i.(SourceEntity)
		entity.Stop()
		Sources.Delete(name)
	}
}

func ListSources() []*api.Source {
	list := make([]*api.Source, 0)
	Sources.Range(func(key, value interface{}) bool {
		list = append(list, value.(SourceEntity).GetConfig())
		return true
	})
	return list
}
