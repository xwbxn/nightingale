package kernel

import (
	"github.com/didi/nightingale/v5/src/webapi/agent"
	"github.com/toolkits/pkg/logger"
)

type KernelStats struct {
	agent.PluginConfig
}

func create() agent.Plugin {
	v := &KernelStats{}
	return v
}

var plugin_name = "kernel"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
