package cpconf_test

import (
	"reflect"
	"testing"

	"github.com/ccfos/nightingale/v6/provider/cpconf"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/cpu"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/disk"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/diskio"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/elasticsearch"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/kafka"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/kernel"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/kernel_vmstat"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/mem"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/mysql"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/net"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/net_response"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/netstat"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/ping"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/redis"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/sockstat"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/sysctl_fs"
	_ "github.com/ccfos/nightingale/v6/provider/cpconf/plugins/system"
)

func TestMarshal(t *testing.T) {
	for name, creator := range cpconf.Plugins {
		println(name)
		conf := creator()
		toml, _ := cpconf.ToToml(conf)
		println(toml)
		json, _ := cpconf.ToJson(conf)
		println(json)
	}
}

func TestUnmarshal(t *testing.T) {
	data := `{"Interval":0,"Devices":null}`
	conf := cpconf.FromJson("cpu", data)
	println(reflect.TypeOf(conf).Elem().Name())
	data = `interval = 0
	collect_protocol_stats = false`
	conf = cpconf.FromToml("net", data)
	println(reflect.TypeOf(conf).Elem().Name())
}
