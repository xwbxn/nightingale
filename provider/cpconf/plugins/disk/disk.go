package disk

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type DiskStats struct {
	cpconf.PluginConfig
	MountPoints       []string `toml:"mount_points"`
	IgnoreFS          []string `toml:"ignore_fs"`
	IgnoreMountPoints []string `toml:"ignore_mount_points"`
}

func create() cpconf.Plugin {
	v := &DiskStats{}
	return v
}

var plugin_name = "disk"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
