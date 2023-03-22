package system

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type SystemStats struct {
	agent.PluginConfig
	CollectUserNumber bool `toml:"collect_user_number"`
}

func create() agent.Plugin {
	v := &SystemStats{}
	return v
}

var plugin_name = "system"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
