package services

import (
	config "SDAS/config"
	"sync"
)

var Pipelines sync.Map

type PIPELINE_ENTITY interface {
	Start()
	Stop()
	GetPipelineString() (string, error)
	GetName() string
	GetType() string
}

func InitPipeline() {
	for _, pipeline_config := range config.Conf.Pipelines {
		switch pipeline_config.Type {
		case "functional":
			// err := addRtspSource(&source_config)
			// if err != nil {
			// 	klog.Error("Init source error: ", err)

			// }
		}
	}

}
