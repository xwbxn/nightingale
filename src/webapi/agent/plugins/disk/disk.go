package disk

import (
	"github.com/didi/nightingale/v5/src/webapi/agent"
	"github.com/toolkits/pkg/logger"
)

type DiskStats struct {
	agent.PluginConfig
	MountPoints       []string `toml:"mount_points"`
	IgnoreFS          []string `toml:"ignore_fs"`
	IgnoreMountPoints []string `toml:"ignore_mount_points"`
}

func create() agent.Plugin {
	v := &DiskStats{}
	return v
}

var plugin_name = "disk"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
