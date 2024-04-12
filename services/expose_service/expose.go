package services

import (
	"SDAS/kitex_gen/api"
	"errors"
	"sync"
)

var Exposes sync.Map

type ExposeEntity interface {
	Start() error
	Stop()
	GetConfig() *api.Expose
}

func BuildExposeEntity(Name, Type, SourceName string, Content map[string]string) (ExposeEntity, error) {
	switch Type {
	case "pull":
		entity, err := NewExposeEntityPull(Name, SourceName, Content)
		if err != nil {
			return nil, err
		}
		return entity, err
	default:
		return nil, errors.New("type" + Name + " is not supported")
	}
}

func ListExposes() []*api.Expose {
	list := make([]*api.Expose, 0)
	Exposes.Range(func(key, value interface{}) bool {
		list = append(list, value.(ExposeEntity).GetConfig())
		return true
	})
	return list
}
