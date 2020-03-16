package conf

import "testing"

var cfgStr = `
goroutine: 4

elastic:
  url: "http://127.0.0.1:9092"
  user: "elastic"
  password: "abc"
  interval: 200
  threshold: 1
  
indices:
- name: "salary"
  app: "xxx"
  lv: "INFO"
  msg: "xxsss"
  
alert:
  url: "http://127.0.0.1:9093/api/v1/alert"
  from: "ElasticLog"
  level: "C"
`

func TestLoad(t *testing.T) {
	cfg, err := Load(cfgStr)
	if err != nil {
		t.Error(err)
	}

	t.Log(cfg.Indices[0])
}
