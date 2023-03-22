package sysctl_fs

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type SysctlFS struct {
	agent.PluginConfig
}

func create() agent.Plugin {
	v := &SysctlFS{}
	return v
}

var plugin_name = "sysctl_fs"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
