package net_response

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type Instance struct {
	agent.InstanceConfig

	Targets     []string       `toml:"targets"`
	Protocol    string         `toml:"protocol"`
	Timeout     agent.Duration `toml:"timeout"`
	ReadTimeout agent.Duration `toml:"read_timeout"`
	Send        string         `toml:"send"`
	Expect      string         `toml:"expect"`
}

type NetResponse struct {
	agent.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func create() agent.Plugin {
	v := &NetResponse{}
	return v
}

var plugin_name = "net_response"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
