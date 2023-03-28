package elasticsearch

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type Elasticsearch struct {
	cpconf.PluginConfig
	Instances []*Instance `toml:"instances"`
}

type Instance struct {
	cpconf.InstanceConfig

	Local                bool            `toml:"local"`
	Servers              []string        `toml:"servers"`
	HTTPTimeout          cpconf.Duration `toml:"http_timeout"`
	ClusterHealth        bool            `toml:"cluster_health"`
	ClusterHealthLevel   string          `toml:"cluster_health_level"`
	ClusterStats         bool            `toml:"cluster_stats"`
	IndicesInclude       []string        `toml:"indices_include"`
	IndicesLevel         string          `toml:"indices_level"`
	NodeStats            []string        `toml:"node_stats"`
	Username             string          `toml:"username"`
	Password             string          `toml:"password"`
	NumMostRecentIndices int             `toml:"num_most_recent_indices"`
}

func create() cpconf.Plugin {
	v := &Elasticsearch{}
	return v
}

var plugin_name = "elasticsearch"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
