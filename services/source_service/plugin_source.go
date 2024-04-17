package services

import (
	"SDAS/config"
	"SDAS/kitex_gen/api"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"path/filepath"
	"plugin"
	"runtime"
	"sync"
	"sync/atomic"
)

/**content
	inSource_entry_map map[string]string
	exist_outSource_map map[string]string
	outSource_expose_map map[string]string
	args json_string
**/

func NewSourceEntityPlugin(name string, expose bool, content map[string]string) ([]*SourceEntityPlugin, error) {
	var entities []*SourceEntityPlugin
	var insourceEntryMap map[string]string
	var existOutsourceMap map[string]string
	var outsourceExposeMap map[string]string
	err := json.Unmarshal([]byte(content["inSource_entry_map"]), &insourceEntryMap)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(content["exist_outSource_map"]), &existOutsourceMap)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(content["outSource_expose_map"]), &outsourceExposeMap)
	if err != nil {
		return nil, err
	}

	sources := make(map[string]*chan *api.SourceMsg)
	for insource, _ := range insourceEntryMap {
		v, ok := Sources.Load(insource)
		if !ok {
			return nil, errors.New("source " + insource + " not exist")
		}
		outChannel, err := v.(SourceEntity).RequestOutChannel()
		if err != nil {
			return nil, err
		}
		sources[insource] = outChannel
	}

	controlChannel := make(chan int)
	rootEntity := &SourceRootEntityPlugin{
		Name:              name,
		Expose:            expose,
		Type:              "plugin",
		Content:           content,
		ControlChannel:    &controlChannel,
		Status:            config.CLOSE,
		Sources:           sources,
		InsourceEntryMap:  insourceEntryMap,
		ExistOutsourceMap: existOutsourceMap,
		Args:              content["args"],
	}
	children := make(map[string]*SourceEntityPlugin)
	for _, outSource := range existOutsourceMap {
		outputChannel := make(chan *api.SourceMsg, 1024)
		var e bool
		if expose {
			e = outsourceExposeMap[outSource] == "true"
		} else {
			e = false
		}
		entity := &SourceEntityPlugin{
			root: rootEntity,

			OutputChannel: &outputChannel,
			Name:          outSource,
			Type:          "plugin",
			Expose:        e,
		}
		entity.requested.Store(false)
		children[outSource] = entity
		entities = append(entities, entity)
	}
	rootEntity.children = children
	return entities, nil
}

type SourceRootEntityPlugin struct {
	Name              string
	Status            int
	Content           map[string]string
	Expose            bool
	Type              string
	once              sync.Once
	children          map[string]*SourceEntityPlugin
	ControlChannel    *chan int
	Sources           map[string]*chan *api.SourceMsg
	InsourceEntryMap  map[string]string
	ExistOutsourceMap map[string]string
	Args              string
}

func (e *SourceRootEntityPlugin) Start() error {
	if e.Status != config.CLOSE && e.Status != config.ERR {
		return errors.New("source already started")
	}
	go e.goroutinePluginSource()
	for {
		switch e.Status {
		case config.OPEN:
			return nil
		case config.ERR:
			err := errors.New("plugin source start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func (e *SourceRootEntityPlugin) Stop() {
	e.once.Do(func() {
		*e.ControlChannel <- config.CLOSE
		for _, child := range e.children {
			Sources.Delete(child.Name)
		}
		close(*e.ControlChannel)
	})
}

func (e *SourceRootEntityPlugin) goroutinePluginSource() {
	p, err := e.startupPlugin(e.Args)
	if err != nil {
		e.Status = config.ERR
		klog.Errorf("source[plugin]: %s open failed. %v", e.Name, err)
		return
	}
	e.Status = config.OPEN
	klog.Infof("source[plugin]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case config.CLOSE:
				for _, outsource := range e.ExistOutsourceMap {
					outChannel := e.children[outsource].OutputChannel
					close(*outChannel)
				}
				e.Status = config.CLOSE
				klog.Infof("source[plugin]: %s closed.\n", e.Name)
				return
			}
		default:

		}

		Entry := make(map[string]map[string]any)
		for insource, channel := range e.Sources {

			msg := <-*channel
			if msg == nil {
				for _, outsource := range e.ExistOutsourceMap {
					outChannel := e.children[outsource].OutputChannel
					close(*outChannel)
				}
				e.Status = config.CLOSE
				klog.Infof("source[plugin]: %s closed.since %s have been closed\n", e.Name, insource)
				return
			}
			Entry[insource] = map[string]any{
				"Data":     msg.Data,
				"DataType": msg.DataType,
				"Ntp":      msg.Ntp,
			}
		}

		Exist, err := p(Entry)

		if err != nil {
			klog.Error(err)
			e.Status = config.ERR
			continue
		}
		for exist, outsource := range e.ExistOutsourceMap {
			outChannel := e.children[outsource].OutputChannel
			*outChannel <- &api.SourceMsg{
				Data:     Exist[exist]["Data"].([]byte),
				DataType: Exist[exist]["DataType"].(string),
				Ntp:      Exist[exist]["Ntp"].(int64),
			}
		}
	}
}

func (e *SourceRootEntityPlugin) startupPlugin(args string) (func(map[string]map[string]any) (map[string]map[string]any, error), error) {
	// FUTURE TODO:从网络下载，github目前没有稳定的下载环境
	//err := utils.CreateFileOrDir(config.Conf.Server.PluginPath)
	//if err != nil {
	//	return nil, err
	//}
	//fileName := fmt.Sprintf("%s.so", e.Name)
	//err = utils.DownloadFileFromUrl(fmt.Sprintf("https://github.com/murInJ/SDAS-plugin/tree/main/plugins/%s", fileName), path.Join(config.Conf.Server.PluginPath, fileName))
	//if err != nil {
	//	return nil, err
	//}
	p, err := e.getPlugin(e.Name, args)
	return p, err
}

func (e *SourceRootEntityPlugin) getPlugin(name string, args string) (func(map[string]map[string]any) (map[string]map[string]any, error), error) {
	pluginPath := filepath.Join(config.Conf.Server.PluginPath, fmt.Sprintf("%s.so", name))

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	symbol, err := p.Lookup("Init")
	if err != nil {
		return nil, err
	}
	fInit := symbol.(func(string) error)
	err = fInit(args)
	if err != nil {
		return nil, err
	}

	symbol, err = p.Lookup("Pipeline")
	if err != nil {
		return nil, err
	}
	f := symbol.(func(map[string]map[string]any) (map[string]map[string]any, error))

	return f, nil
}

type SourceEntityPlugin struct {
	root          *SourceRootEntityPlugin
	OutputChannel *chan *api.SourceMsg
	Name          string
	Type          string
	Expose        bool
	requested     atomic.Bool
}

func (e *SourceEntityPlugin) GetConfig() *api.Source {
	return &api.Source{
		Type:    e.Type,
		Name:    e.Name,
		Expose:  e.Expose,
		Content: e.root.Content,
	}
}

func (e *SourceEntityPlugin) RequestOutChannel() (*chan *api.SourceMsg, error) {
	if e.requested.CompareAndSwap(false, true) {
		return e.OutputChannel, nil
	} else {
		return nil, errors.New("request out channel already in use")
	}
}

func (e *SourceEntityPlugin) ReleaseOutChannel() {
	e.requested.Store(false)
}

func (e *SourceEntityPlugin) Start() error {
	err := e.root.Start()
	if err != nil {
		return err
	}
	return nil
}

func (e *SourceEntityPlugin) Stop() {
	e.root.Stop()
}

func (e *SourceEntityPlugin) GetName() string {
	return e.Name
}
