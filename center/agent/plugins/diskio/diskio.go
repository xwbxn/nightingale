package diskio

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type DiskIO struct {
	agent.PluginConfig
	Devices []string `toml:"devices"`
}

func create() agent.Plugin {
	v := &DiskIO{}
	return v
}

var plugin_name = "diskio"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
