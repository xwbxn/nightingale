package cpu

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type CPUStats struct {
	agent.PluginConfig
	CollectPerCPU bool `toml:"collect_per_cpu"`
}

func create() agent.Plugin {
	v := &CPUStats{}
	return v
}

var plugin_name = "cpu"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
