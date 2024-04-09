package main

import (
	api "SDAS/kitex_gen/api/sdas"
	"fmt"
	"net"
	"os"

	config "SDAS/config"
	s "SDAS/services"
	u "SDAS/utils"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

func main() {
	/**
	 * LOG
	 */
	f, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		klog.Error(err)
	}
	defer f.Close()
	klog.SetOutput(f)
	/**
	 * init
	 */
	config.LoadConfig("./config.json")
	u.InitUtils()
	s.InitServices()

	/**
	 * SERVER
	 */
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.Server.Host, config.Conf.Server.Port))
	svr := api.NewServer(new(SDASImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		klog.Error(err)
	}
}
