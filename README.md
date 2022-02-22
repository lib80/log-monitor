## 简介
log-monitor 用于对日志文件实时监控，通过指定的关键字过滤出需要的信息，并推送到机器人接收器。
- 目前机器人消息接收器支持钉钉和云之家
- 可根据过滤关键字对多个日志文件进行分组
## 安装
- 安装Go并设置环境变量GOPATH
- `go install github.com/lib80/log-monitor`
- `cd $GOPATH/src/github.com/lib80/log-monitor`
- `make`
## 使用
- 编写yml格式的配置文件，参考示例文件`config.example.yml`
- 启动`./bin/log-monitor -f log-monitor.yml`
