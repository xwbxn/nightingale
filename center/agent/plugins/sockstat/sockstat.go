package sockstat

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type SockStat struct {
	agent.PluginConfig
	Protocols []string `toml:"protocols"`
}

func create() agent.Plugin {
	v := &SockStat{}
	return v
}

var plugin_name = "sockstat"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
