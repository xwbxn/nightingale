package system

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type SystemStats struct {
	cpconf.PluginConfig
	CollectUserNumber bool `toml:"collect_user_number"`
}

func create() cpconf.Plugin {
	v := &SystemStats{}
	return v
}

var plugin_name = "system"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
