package sockstat

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type SockStat struct {
	cpconf.PluginConfig
	Protocols []string `toml:"protocols"`
}

func create() cpconf.Plugin {
	v := &SockStat{}
	return v
}

var plugin_name = "sockstat"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
