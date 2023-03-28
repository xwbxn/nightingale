package net_response

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type Instance struct {
	cpconf.InstanceConfig

	Targets     []string        `toml:"targets"`
	Protocol    string          `toml:"protocol"`
	Timeout     cpconf.Duration `toml:"timeout"`
	ReadTimeout cpconf.Duration `toml:"read_timeout"`
	Send        string          `toml:"send"`
	Expect      string          `toml:"expect"`
}

type NetResponse struct {
	cpconf.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func create() cpconf.Plugin {
	v := &NetResponse{}
	return v
}

var plugin_name = "net_response"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
