package handlers

import (
	config "SDAS/config"
	source_service "SDAS/services/source_service"
)

func AddSourceHandler(Name, Type, Content string) error {
	switch Type {
	case "rtsp":
		source := config.SOURCE{
			Type:    Type,
			Name:    Name,
			Content: Content,
		}
		err := source_service.AddRtspSource(&source)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveSourceHandler(Name string) {
	source_service.RemoveSource(Name)
}

func ListSourceHandler() []config.SOURCE {
	return source_service.ListSources()
}
