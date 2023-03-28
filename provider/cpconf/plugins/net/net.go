package net

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type NetIOStats struct {
	cpconf.PluginConfig
	CollectProtocolStats bool     `toml:"collect_protocol_stats"`
	Interfaces           []string `toml:"interfaces"`
}

func create() cpconf.Plugin {
	v := &NetIOStats{}
	return v
}

var plugin_name = "net"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
