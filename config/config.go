package config

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DingDing struct {
	AccessToken string   `yaml:"access_token" validate:"required"`
	Secret      string   `yaml:"secret" validate:"required"`
	AtMobiles   []string `yaml:"at_mobiles" validate:"unique"`
	IsAtAll     bool     `yaml:"is_at_all"`
}

type YunZhiJia struct {
	Webhook string   `yaml:"webhook" validate:"required,url"`
	AtNames []string `yaml:"at_names" validate:"unique"`
	IsAtAll bool     `yaml:"is_at_all"`
}

type MonitorItem struct {
	Paths           []string `yaml:"paths" validate:"required,gt=0,unique"`
	IncludeKeyWords []string `yaml:"include_key_words" validate:"required,gt=0,unique"`
}

type Config struct {
	DingDing       *DingDing      `yaml:"dingding"`
	YunZhiJia      *YunZhiJia     `yaml:"yunzhijia" `
	Receivers      []string       `yaml:"receivers" validate:"required,gt=0,unique"`
	MonitorTargets []*MonitorItem `yaml:"monitor_targets" validate:"required,gt=0"`
}

//Validate 参数合法性校验
func Validate(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}

func LoadConfigFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if err := Validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
