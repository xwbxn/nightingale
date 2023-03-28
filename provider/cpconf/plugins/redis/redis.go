package redis

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type Command struct {
	Command []interface{} `toml:"command"`
	Metric  string        `toml:"metric"`
}

type Instance struct {
	cpconf.InstanceConfig

	Address  string    `toml:"address"`
	Username string    `toml:"username"`
	Password string    `toml:"password"`
	PoolSize int       `toml:"pool_size"`
	Commands []Command `toml:"commands"`
}

type Redis struct {
	cpconf.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func create() cpconf.Plugin {
	v := &Redis{}
	return v
}

var plugin_name = "redis"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
