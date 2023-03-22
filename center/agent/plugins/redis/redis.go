package redis

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type Command struct {
	Command []interface{} `toml:"command"`
	Metric  string        `toml:"metric"`
}

type Instance struct {
	agent.InstanceConfig

	Address  string    `toml:"address"`
	Username string    `toml:"username"`
	Password string    `toml:"password"`
	PoolSize int       `toml:"pool_size"`
	Commands []Command `toml:"commands"`
}

type Redis struct {
	agent.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func create() agent.Plugin {
	v := &Redis{}
	return v
}

var plugin_name = "redis"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
