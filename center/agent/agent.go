package agent

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/toolkits/pkg/logger"
)

type Plugin interface{}

type Creator func() Plugin

var Plugins = map[string]Creator{}

type PluginConfig struct {
	Interval Duration `toml:"interval"`
}

func AddPlugin(name string, creator Creator) {
	Plugins[name] = creator
}

func CreateConf(name string) interface{} {
	creator, has := Plugins[name]
	if has {
		conf := creator()
		return conf
	} else {
		return nil
	}
}

func ToToml(c interface{}) (string, error) {

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ToJson(c interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(c)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func FromJson(name string, data string) interface{} {
	creator, has := Plugins[name]
	if !has {
		return nil
	}
	conf := creator()
	err := json.Unmarshal([]byte(data), conf)
	if err != nil {
		logger.Warningf("[-] json decode fail: %s, json: %s", err.Error(), data)
		return nil
	}
	return conf
}

func FromToml(name string, data string) interface{} {
	creator, has := Plugins[name]
	if !has {
		return nil
	}
	conf := creator()
	err := toml.Unmarshal([]byte(data), conf)
	if err != nil {
		logger.Warningf("[-] json decode fail: %s, json: %s", err.Error(), data)
		return nil
	}
	return conf
}

type InstanceConfig struct {
	IntervalTimes int64 `toml:"interval_times"`
}

// Duration is a time.Duration
type Duration time.Duration

// UnmarshalTOML parses the duration from the TOML config file
func (d *Duration) UnmarshalTOML(b []byte) error {
	// convert to string
	durStr := string(b)

	// Value is a TOML number (e.g. 3, 10, 3.5)
	// First try parsing as integer seconds
	sI, err := strconv.ParseInt(durStr, 10, 64)
	if err == nil {
		dur := time.Second * time.Duration(sI)
		*d = Duration(dur)
		return nil
	}
	// Second try parsing as float seconds
	sF, err := strconv.ParseFloat(durStr, 64)
	if err == nil {
		dur := time.Second * time.Duration(sF)
		*d = Duration(dur)
		return nil
	}

	// Finally, try value is a TOML string (e.g. "3s", 3s) or literal (e.g. '3s')
	durStr = strings.ReplaceAll(durStr, "'", "")
	durStr = strings.ReplaceAll(durStr, "\"", "")
	if durStr == "" {
		durStr = "0s"
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return err
	}

	*d = Duration(dur)
	return nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	return d.UnmarshalTOML(text)
}

// UnmarshalTOML parses the duration from the TOML config file
func (d *Duration) UnmarshalJSON(b []byte) error {
	// convert to string
	durStr := string(b)

	// Value is a TOML number (e.g. 3, 10, 3.5)
	// First try parsing as integer seconds
	sI, err := strconv.ParseInt(durStr, 10, 64)
	if err == nil {
		dur := time.Second * time.Duration(sI)
		*d = Duration(dur)
		return nil
	}
	// Second try parsing as float seconds
	sF, err := strconv.ParseFloat(durStr, 64)
	if err == nil {
		dur := time.Second * time.Duration(sF)
		*d = Duration(dur)
		return nil
	}

	// Finally, try value is a TOML string (e.g. "3s", 3s) or literal (e.g. '3s')
	durStr = strings.ReplaceAll(durStr, "'", "")
	durStr = strings.ReplaceAll(durStr, "\"", "")
	if durStr == "" {
		durStr = "0s"
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return err
	}

	*d = Duration(dur)
	return nil
}
