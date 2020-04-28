package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Goroutine int        `yaml:"goroutine"`
	Ticker    int        `yaml:"ticker"`
	Interval  int        `yaml:"interval"`
	Overcount int        `yaml:"overcount"`
	Elastic   *Elastic   `yaml:"elastic"`
	Indices   []*Indices `yaml:"indices"`
	Alert     *Alert     `yaml:"alert"`
}

type Elastic struct {
	URL      []string `yaml:"url"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
}

type Indices struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Interval  int    `yaml:"interval,omitempty"`
	Overcount int    `yaml:"overcount,omitempty"`
	Include   string `yaml:"include"`
	Exclude   string `yaml:"exclude,omitempty"`
}

type Alert struct {
	URL   string `yaml:"url"`
	From  string `yaml:"from"`
	Level string `yaml:"level"`
}

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

func Load(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	for _, index := range cfg.Indices {
		if index.Interval == 0 {
			index.Interval = cfg.Interval
		}
		if index.Overcount == 0 {
			index.Overcount = cfg.Overcount
		}
	}

	return cfg, nil
}
