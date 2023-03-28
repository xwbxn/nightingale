package mem

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type MemStats struct {
	cpconf.PluginConfig
	CollectPlatformFields bool `toml:"collect_platform_fields"`
}

func create() cpconf.Plugin {
	v := &MemStats{}
	return v
}

var plugin_name = "disk"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
