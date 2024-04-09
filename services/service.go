package services

import (
	source "SDAS/services/source_service"
	"sync"
)

var Exposes sync.Map
var Pipeline any

func InitServices() {
	source.InitSource()
}
