package mysql

import (
	"github.com/ccfos/nightingale/v6/provider/cpconf"
	"github.com/toolkits/pkg/logger"
)

type QueryConfig struct {
	Mesurement    string          `toml:"mesurement"`
	LabelFields   []string        `toml:"label_fields"`
	MetricFields  []string        `toml:"metric_fields"`
	FieldToAppend string          `toml:"field_to_append"`
	Timeout       cpconf.Duration `toml:"timeout"`
	Request       string          `toml:"request"`
}

type Instance struct {
	cpconf.InstanceConfig

	Address        string `toml:"address"`
	Username       string `toml:"username"`
	Password       string `toml:"password"`
	Parameters     string `toml:"parameters"`
	TimeoutSeconds int64  `toml:"timeout_seconds"`

	Queries       []QueryConfig `toml:"queries"`
	GlobalQueries []QueryConfig `toml:"-"`

	ExtraStatusMetrics              bool `toml:"extra_status_metrics"`
	ExtraInnodbMetrics              bool `toml:"extra_innodb_metrics"`
	GatherProcessListProcessByState bool `toml:"gather_processlist_processes_by_state"`
	GatherProcessListProcessByUser  bool `toml:"gather_processlist_processes_by_user"`
	GatherSchemaSize                bool `toml:"gather_schema_size"`
	GatherTableSize                 bool `toml:"gather_table_size"`
	GatherSystemTableSize           bool `toml:"gather_system_table_size"`
	GatherSlaveStatus               bool `toml:"gather_slave_status"`
}

type MySQL struct {
	cpconf.PluginConfig
	Instances []*Instance   `toml:"instances"`
	Queries   []QueryConfig `toml:"queries"`
}

func create() cpconf.Plugin {
	v := &MySQL{}
	return v
}

var plugin_name = "disk"

func init() {
	cpconf.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
