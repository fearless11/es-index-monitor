package conf

import "testing"

var cfgStr = `
goroutine: 4
interval: 5
overcount: 50

elastic:
  url: ["http://127.0.0.1:9092"]
  user: "elastic"
  password: "abc"

indices:
- name: "crm-frontend-app-nginx"
  type: "eq"
  include: "status=500"
  interval: 5
  overcount: 50
- name: "crm-frontend-app-nginx"
  type: "gt"
  interval: 5
  overcount: 50
  include: "request_time>60"
  exclude: "status=101;status=499"
  
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
