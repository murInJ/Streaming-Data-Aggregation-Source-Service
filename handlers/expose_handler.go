package handlers

import (
	config "SDAS/config"
	expose_service "SDAS/services/expose_service"
)

func AddExposeHandler(Name, Type, Content, SourceName string) error {
	switch Type {
	case "rtsp":
		expose := config.EXPOSE{
			Type:       Type,
			Name:       Name,
			Content:    Content,
			SourceName: SourceName,
		}
		err := expose_service.AddRtspExpose(&expose)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveExposeHandler(Name string) {
	expose_service.RemoveExpose(Name)
}

func ListExposeHandler() []config.EXPOSE {
	return expose_service.ListExposes()
}
