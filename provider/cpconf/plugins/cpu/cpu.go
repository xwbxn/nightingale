package cpu

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type CPUStats struct {
	cpconf.PluginConfig
	CollectPerCPU bool `toml:"collect_per_cpu"`
}

func create() cpconf.Plugin {
	v := &CPUStats{}
	return v
}

var plugin_name = "cpu"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
