package agent_test

import (
	"reflect"
	"testing"

	"github.com/didi/nightingale/v5/src/webapi/agent"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/cpu"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/disk"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/diskio"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/elasticsearch"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/kafka"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/kernel"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/kernel_vmstat"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/mem"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/mysql"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/net"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/net_response"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/netstat"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/ping"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/redis"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/sockstat"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/sysctl_fs"
	_ "github.com/didi/nightingale/v5/src/webapi/agent/plugins/system"
)

func TestMarshal(t *testing.T) {
	for name, creator := range agent.Plugins {
		println(name)
		conf := creator()
		toml, _ := agent.ToToml(conf)
		println(toml)
		json, _ := agent.ToJson(conf)
		println(json)
	}
}

func TestUnmarshal(t *testing.T) {
	data := `{"Interval":0,"Devices":null}`
	conf := agent.FromJson("cpu", data)
	println(reflect.TypeOf(conf).Elem().Name())
	data = `interval = 0
	collect_protocol_stats = false`
	conf = agent.FromToml("net", data)
	println(reflect.TypeOf(conf).Elem().Name())
}
