package sysctl_fs

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type SysctlFS struct {
	cpconf.PluginConfig
}

func create() cpconf.Plugin {
	v := &SysctlFS{}
	return v
}

var plugin_name = "sysctl_fs"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
