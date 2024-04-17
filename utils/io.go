package utils

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"io"
	"net/http"
	"os"
)

func DownloadFileFromUrl(url string, savePath string) error {
	klog.Infof("DownloadFileFromUrl url: %s, savePath: %s", url, savePath)
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 获取远程文件
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	// 创建一个文件用于保存下载的内容
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将下载的内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CheckFileOrDirExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFileOrDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		if fi, err := os.Lstat(path); err == nil && fi.IsDir() {
			return nil
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("无法创建文件或文件夹 '%s': %v", path, err)
		}
	} else {
		// 其他错误
		return fmt.Errorf("检查路径 '%s' 状态时发生错误: %v", path, err)
	}
	return nil
}
