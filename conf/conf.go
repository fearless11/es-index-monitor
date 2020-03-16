package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config 顶层配置
type Config struct {
	Goroutine int        `yaml:"goroutine"`
	Elastic   *Elastic   `yaml:"elastic"`
	Indices   []*Indices `yaml:"indices"`
	Alert     *Alert     `yaml:"alert"`
}

type Elastic struct {
	URL       []string `yaml:"url"`
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
	Interval  int      `yaml:"interval"`
	Threshold int      `yaml:"threshold"`
}

type Indices struct {
	Name string `yaml:"name"`
	App  string `yaml:"app"`
	Lv   string `yaml:"lv"`
	Msg  string `yaml:"msg,omitempty"`
}

// Alert 告警
type Alert struct {
	URL   string `yaml:"url"`
	From  string `yaml:"from"`
	Level string `yaml:"level"`
}

// LoadFile 读yaml文件
func LoadFile(filename string) *Config {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := Load(string(content))
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

// Load 解析yaml文件
func Load(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
