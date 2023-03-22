package netstat

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type NetStats struct {
	agent.PluginConfig

	DisableConnectionStats bool `toml:"disable_connection_stats"`
	TcpExt                 bool `toml:"tcp_ext"`
	IpExt                  bool `toml:"ip_ext"`
}

func create() agent.Plugin {
	v := &NetStats{}
	return v
}

var plugin_name = "netstat"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
