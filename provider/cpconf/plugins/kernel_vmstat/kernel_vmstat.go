package kernel_vmstat

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type KernelVmstat struct {
	cpconf.PluginConfig
	WhiteList map[string]int `toml:"white_list"`
}

func create() cpconf.Plugin {
	v := &KernelVmstat{}
	return v
}

var plugin_name = "kernel_vmstat"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
