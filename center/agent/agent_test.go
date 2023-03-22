package agent_test

import (
	"reflect"
	"testing"

	"github.com/ccfos/nightingale/v6/center/agent"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/cpu"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/disk"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/diskio"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/elasticsearch"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/kafka"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/kernel"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/kernel_vmstat"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/mem"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/mysql"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/net"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/net_response"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/netstat"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/ping"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/redis"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/sockstat"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/sysctl_fs"
	_ "github.com/ccfos/nightingale/v6/center/agent/plugins/system"
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
