package diskio

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type DiskIO struct {
	cpconf.PluginConfig
	Devices []string `toml:"devices"`
}

func create() cpconf.Plugin {
	v := &DiskIO{}
	return v
}

var plugin_name = "diskio"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
