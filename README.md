# Streaming-Data-Aggregation-Source-Service(流数据聚合源服务)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/murInJ/Streaming-Data-Aggregation-Source-Service.svg)](https://github.com/murInJ/Streaming-Data-Aggregation-Source-Service)
![GitHub Release](https://img.shields.io/github/v/release/murInJ/Streaming-Data-Aggregation-Source-Service)
[![GitHub contributors](https://img.shields.io/github/contributors/MurInJ/Streaming-Data-Aggregation-Source-Service.svg)](https://GitHub.com/MurInJ/Streaming-Data-Aggregation-Source-Service/graphs/contributors/)
<!-- ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/murInJ/Streaming-Data-Aggregation-Source-Service/go.yml) -->
该微服务的作用在于将多个分散的数据源进行转发，同时聚合到一起，进行处理后再对外提供服务。

`sh script/deploy_build.sh`执行部署环境的二进制和docker构建

## 现在支持的
1. Source
    - `rtsp` 从rtsp服务器拉流获取数据
    - `plugin` 从指定源获取数据，并通过指定插件的处理函数进行数据处理后，成为指定源。其中具体编写方法请参考[plugin仓库](https://github.com/murInJ/SDAS-plugin)
    - `remote` 从指定url的SDAS拉取数据
3. Expose
   - `pull stream` 通过rpc stream获取指定源的数据
   - `http push` 通过指定url进行http api接口推送数据
## Contributors
<a href="https://github.com/MurInj/Streaming-Data-Aggregation-Source-Service/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=MurInj/Streaming-Data-Aggregation-Source-Service" />
</a>