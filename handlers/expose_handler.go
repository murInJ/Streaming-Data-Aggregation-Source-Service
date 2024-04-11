package handlers

import (
	config "SDAS/config"
	expose_service "SDAS/services/expose_service"
	"errors"
	"fmt"
)

func AddExposeHandler(Name, Type, Content, SourceName string) error {

	expose := config.EXPOSE{
		Type:       Type,
		Content:    Content,
		Name:       Name,
		SourceName: SourceName,
	}
	if _, ok := expose_service.Exposes.Load(expose.Name); ok {
		return fmt.Errorf("expose %s already exist", expose.Name)
	}
	switch Type {
	case "pull":
		err := expose_service.AddPullExpose(expose.Name, expose.SourceName, expose.Content) //此处的Content为msgtype
		if err != nil {
			return err
		}
		// case "rtsp":
		// 	expose := config.EXPOSE{
		// 		Type:       Type,
		// 		Name:       Name,
		// 		Content:    Content,
		// 		SourceName: SourceName,
		// 	}
		// 	err := expose_service.AddRtspExpose(&expose)
		// 	if err != nil {
		// 		return err
		// 	}
	}

	return nil
}

func RemoveExposeHandler(Name string) {
	expose_service.RemoveExpose(Name)
}

func ListExposeHandler() []config.EXPOSE {
	return expose_service.ListExposes()
}

func PlayExposePull(name string) error {
	i, ok := expose_service.Exposes.Load(name)
	if !ok {
		return errors.New("expose not found")
	}
	v := i.(expose_service.EXPOSE_ENTITY)
	*v.GetControlChannel() <- expose_service.PLAY
	return nil

}

func PauseExposePull(name string) error {
	if i, ok := expose_service.Exposes.Load(name); ok {
		v := i.(expose_service.EXPOSE_ENTITY)
		*v.GetControlChannel() <- expose_service.PAUSE
		return nil
	}
	return errors.New("expose not found")
}

func GetExposePullData(name string) (string, error) {
	for {
		if i, ok := expose_service.Exposes.Load(name); ok {
			v := i.(expose_service.EXPOSE_ENTITY)
			data, err := v.GetData()
			if err != nil {
				return "", err
			}
			return data.(string), nil
		}
	}
}
