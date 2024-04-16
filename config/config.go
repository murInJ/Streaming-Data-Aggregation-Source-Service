package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	OPEN  = 0
	PLAY  = 1
	PAUSE = 2
	CLOSE = 3
	ERR   = 4
)

type Config struct {
	Server SERVER `json:"server"`
}

var Conf *Config

func NewDefaultConfig() {
	Conf = &Config{
		Server: SERVER{
			Node:       "test_sdas",
			Host:       "0.0.0.0",
			Port:       8088,
			PluginPath: "./plugins",
		},
	}
}

func SaveConfigJSON(filename string, conf *Config) error {
	content, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to marshal config: %w", err)
		klog.Error(err)
		return err
	}
	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		err = fmt.Errorf("failed to write file: %w", err)
		klog.Error(err)
		return err
	}
	klog.Info("save config to %s", filename)
	return nil
}
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		// 如果文件不存在，创建并使用默认配置
		NewDefaultConfig()
		err := SaveConfigJSON(filename, Conf)
		if err != nil {
			err = fmt.Errorf("failed to save default config: %w", err)
			klog.Error(err)
			return err
		}
	} else if err != nil {
		// 读取文件失败，返回错误
		return fmt.Errorf("failed to read file: %w", err)
	} else {
		var config *Config
		// 成功读取文件，解析JSON内容
		err = json.Unmarshal(content, &config)
		if err != nil {
			// 解析JSON失败，返回错误
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		Conf = config
	}
	return nil
}
