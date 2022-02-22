package config

import (
	"fmt"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	cfg, err := LoadConfigFile("../config.example.yml")
	if err != nil {
		t.Errorf("加载配置文件出错：%v", err)
		return
	}
	fmt.Println(cfg)
}
