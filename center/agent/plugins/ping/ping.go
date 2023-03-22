package ping

import (
	"github.com/ccfos/nightingale/v6/center/agent"
	"github.com/toolkits/pkg/logger"
)

type Instance struct {
	agent.InstanceConfig

	Targets      []string `toml:"targets"`
	Count        int      `toml:"count"`         // ping -c <COUNT>
	PingInterval float64  `toml:"ping_interval"` // ping -i <INTERVAL>
	Timeout      float64  `toml:"timeout"`       // ping -W <TIMEOUT>
	Interface    string   `toml:"interface"`     // ping -I/-S <INTERFACE/SRC_ADDR>
	IPv6         bool     `toml:"ipv6"`          // Whether to resolve addresses using ipv6 or not.
	Size         *int     `toml:"size"`          // Packet size
	Conc         int      `toml:"concurrency"`   // max concurrency coroutine
}

type Ping struct {
	agent.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func create() agent.Plugin {
	v := &Ping{}
	return v
}

var plugin_name = "ping"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
