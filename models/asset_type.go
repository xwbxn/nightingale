package models

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"text/template"

	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
)

type AssetType struct {
	Name       string                    `json:"name"`
	Plugin     string                    `json:"plugin"`
	Metrics    []*AssetTypeMetric        `json:"metrics"`
	Category   string                    `json:"category"`
	Form       []*BaseProp               `json:"form"`
	IpUnique   bool                      `json:"ip_unique"`
	BaseProps  []*BaseProp               `json:"base_props" yaml:"base_props"`
	ExtraProps map[string]*ExtraPropPart `json:"extra_props" yaml:"extra_props"`

	Dashboard string `json:"-"`
}

type AssetTypeMetric struct {
	Name    string      `json:"name"`
	Metrics string      `json:"metrics" yaml:"metrics"`
	Value   model.Value `json:"-" yaml:"-"`
	Unit    string      `json:"unit"`
	Label   string      `json:"label"`
	Rules   []struct {
		Name      string
		Serverity int64
		Relation  string
		Value     int64
	}
}

type AssetConfigs struct {
	Config []*AssetType `json:"config"`
}

type Unit string

const (
	Byte Unit = "byte"
	Bit  Unit = "bit"
	Time Unit = "time"
)

type BaseProp struct {
	Name     string `json:"name" yaml:"name"`
	Label    string `json:"label"`
	Required string `json:"required"`
	Type     string `json:"type"`
	Unit     Unit   `json:"unit"`
	Options  []*struct {
		Label string `json:"label" yaml:"label"`
		Value string `json:"value" yaml:"value"`
	} `json:"options" yaml:"options"`
}

type ExtraPropPart struct {
	Label string `json:"label" yaml:"label"`
	Sort  int64  `json:"sort" yaml:"sort"`
	Props []*struct {
		Name       string      `json:"name" yaml:"name"`
		Label      string      `json:"label" yaml:"label"`
		Type       string      `json:"type" yaml:"type"`
		ItemsLimit int64       `json:"items_limit" yaml:"items_limit"`
		Items      []*BaseProp `json:"items" yaml:"items"`
	} `json:"props" yaml:"props"`
}

func AssetTypeGetsAll() ([]*AssetType, error) {
	fp := path.Join("etc", "assets.yaml")
	var assetConfig AssetConfigs
	err := file.ReadYaml(fp, &assetConfig)
	return assetConfig.Config, err
}

func AssetTypeGet(name string) (*AssetType, error) {
	all, err := AssetTypeGetsAll()
	if err != nil {
		return nil, err
	}
	for _, t := range all {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%s not found", name))
}

func (at *AssetType) GenConfig(f map[string]interface{}) (string, error) {
	filepath := path.Join("etc", "default", fmt.Sprintf("%s.toml", at.Plugin))

	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		logger.Errorf("unable to open config: %s", filepath)
		return "", err
	}
	var content bytes.Buffer
	tpl.Execute(&content, f)
	return content.String(), nil
}
