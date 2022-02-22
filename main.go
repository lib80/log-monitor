package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"

	"github.com/lib80/log-monitor/config"
	"github.com/lib80/log-monitor/notify"
)

var hostname, _ = os.Hostname()

//monitor 监控日志文件，将满足过滤条件的信息推送给接收器
func monitor(wg *sync.WaitGroup, path string, includeKeyWords []string, receivers []notify.Receiver) {
	defer wg.Done()
	t, err := tail.TailFile(path, tail.Config{ReOpen: true, MustExist: true, Follow: true})
	if err != nil {
		log.Printf("[%v]监控文件出错：%v\n", path, err)
		os.Exit(1)
	}

	for line := range t.Lines {
		for _, kw := range includeKeyWords {
			if strings.Contains(line.Text, kw) {
				now := time.Now().Format("2006-01-02 15:04:05")
				msg := fmt.Sprintf("[日志异常告警]\n异常信息：%v\n日志文件：%v\n主机：%v\n时间：%v\n", line.Text, path, hostname, now)
				for _, receiver := range receivers {
					if err := receiver.Fire(msg); err != nil {
						log.Printf("发送到接收器出错：%v\n", err)
						os.Exit(1)
					}
				}
				break
			}
		}
	}
}

func main() {
	//加载配置文件
	var configFile string
	flag.StringVar(&configFile, "f", "config.example.yml", "config file path")
	flag.Parse()
	cfg, err := config.LoadConfigFile(configFile)
	if err != nil {
		log.Printf("加载配置文件出错：%v\n", err)
		os.Exit(1)
	}

	//创建接收器
	receivers := make([]notify.Receiver, 0)
	for _, rcvName := range cfg.Receivers {
		receiver, err := notify.NewReceiver(rcvName, cfg)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		receivers = append(receivers, receiver)
	}

	//开启多个协程并发监控日志文件
	wg := &sync.WaitGroup{}
	for _, monitorItem := range cfg.MonitorTargets {
		for _, path := range monitorItem.Paths {
			wg.Add(1)
			go monitor(wg, path, monitorItem.IncludeKeyWords, receivers)
		}
	}

	wg.Wait()
}
