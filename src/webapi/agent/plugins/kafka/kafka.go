package kafka

import (
	"github.com/didi/nightingale/v5/src/webapi/agent"
	"github.com/toolkits/pkg/logger"
)

type Kafka struct {
	agent.PluginConfig
	Instances []*Instance `toml:"instances"`
}

type Instance struct {
	agent.InstanceConfig

	LogLevel string `toml:"log_level"`

	// Address (host:port) of Kafka server.
	KafkaURIs []string `toml:"kafka_uris,omitempty"`

	// Connect using SASL/PLAIN
	UseSASL bool `toml:"use_sasl,omitempty"`

	// Only set this to false if using a non-Kafka SASL proxy
	UseSASLHandshake *bool `toml:"use_sasl_handshake,omitempty"`

	// SASL user name
	SASLUsername string `toml:"sasl_username,omitempty"`

	// SASL user password
	SASLPassword string `toml:"sasl_password,omitempty"`

	// The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism
	SASLMechanism string `toml:"sasl_mechanism,omitempty"`

	// Connect using TLS
	UseTLS bool `toml:"use_tls,omitempty"`

	// The optional certificate authority file for TLS client authentication
	CAFile string `toml:"ca_file,omitempty"`

	// The optional certificate file for TLS client authentication
	CertFile string `toml:"cert_file,omitempty"`

	// The optional key file for TLS client authentication
	KeyFile string `toml:"key_file,omitempty"`

	// If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
	InsecureSkipVerify bool `toml:"insecure_skip_verify,omitempty"`

	// Kafka broker version
	KafkaVersion string `toml:"kafka_version,omitempty"`

	// if you need to use a group from zookeeper
	UseZooKeeperLag bool `toml:"use_zookeeper_lag,omitempty"`

	// Address array (hosts) of zookeeper server.
	ZookeeperURIs []string `toml:"zookeeper_uris,omitempty"`

	// Metadata refresh interval
	MetadataRefreshInterval string `toml:"metadata_refresh_interval,omitempty"`

	// If true, all scrapes will trigger kafka operations otherwise, they will share results. WARN: This should be disabled on large clusters
	AllowConcurrent *bool `toml:"allow_concurrency,omitempty"`

	// Maximum number of offsets to store in the interpolation table for a partition
	MaxOffsets int `toml:"max_offsets,omitempty"`

	// How frequently should the interpolation table be pruned, in seconds
	PruneIntervalSeconds int `toml:"prune_interval_seconds,omitempty"`

	// Regex filter for topics to be monitored
	TopicsFilter string `toml:"topics_filter_regex,omitempty"`

	// Regex filter for consumer groups to be monitored
	GroupFilter string `toml:"groups_filter_regex,omitempty"`
}

func create() agent.Plugin {
	v := &Kafka{}
	return v
}

var plugin_name = "kafka"

func init() {
	agent.AddPlugin(plugin_name, create)
	logger.Infof("[+] %s registerd", plugin_name)
}
