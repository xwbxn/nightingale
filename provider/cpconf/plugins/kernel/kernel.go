package kernel

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type KernelStats struct {
	cpconf.PluginConfig
}

func create() cpconf.Plugin {
	v := &KernelStats{}
	return v
}

var plugin_name = "kernel"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
