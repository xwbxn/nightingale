package mem

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type MemStats struct {
	agent.PluginConfig
	CollectPlatformFields bool `toml:"collect_platform_fields"`
}

func create() agent.Plugin {
	v := &MemStats{}
	return v
}

var plugin_name = "disk"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
