package net

import (
	"github.com/didi/nightingale/v5/src/webapi/agent"
	"github.com/toolkits/pkg/logger"
)

type NetIOStats struct {
	agent.PluginConfig
	CollectProtocolStats bool     `toml:"collect_protocol_stats"`
	Interfaces           []string `toml:"interfaces"`
}

func create() agent.Plugin {
	v := &NetIOStats{}
	return v
}

var plugin_name = "net"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
