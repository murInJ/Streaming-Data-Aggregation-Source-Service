package services

import (
	expose "SDAS/services/expose_service"
	pipeline "SDAS/services/pipeline_service"
	source "SDAS/services/source_service"
	"sync"
)

var Exposes sync.Map

func InitServices() {
	source.InitSource()
	pipeline.InitPipeline()
	expose.InitExpose()
}
