package kernel_vmstat

import (
	"github.com/didi/nightingale/v5/src/webapi/agent"
	"github.com/toolkits/pkg/logger"
)

type KernelVmstat struct {
	agent.PluginConfig
	WhiteList map[string]int `toml:"white_list"`
}

func create() agent.Plugin {
	v := &KernelVmstat{}
	return v
}

var plugin_name = "kernel_vmstat"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
