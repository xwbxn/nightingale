package router

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

type ConfigFormat string

const (
	YamlFormat ConfigFormat = "yaml"
	TomlFormat ConfigFormat = "toml"
	JsonFormat ConfigFormat = "json"
)

type ConfigWithFormat struct {
	Config string       `json:"config"`
	Format ConfigFormat `json:"format"`
}

type httpRemoteProviderResponse struct {
	// version is signature/md5 of current Config, server side should deal with the Version calculate
	Version string `json:"version"`

	// ConfigMap (InputName -> Config), if version is identical, server side can set Config to nil
	Configs map[string]map[string]*ConfigWithFormat `json:"configs"`
}

//	 example response:
//	 {
//	  "version": "111",
//	  "configs": {
//	    "mysql": {
//		  "name": {
//	        "config": "# # collect interval\n# interval = 15\n\n[[ instances ]]\naddress = \"172.33.44.55:3306\"\nusername = \"111\"\npassword = \"2222\"\nlabels = { instance = \"mysql2\"}\nextra_innodb_metrics =true",
//	        "format": "toml"
//		  }
//	    }
//	  }
//	}
func (rt *Router) CategrafConfigGet(c *gin.Context) {
	ident := ginx.QueryStr(c, "agent_hostname")
	version := ginx.QueryStr(c, "version", "")

	provider := rt.AssetCache.GetByIdent(ident)
	if len(provider) == 0 {
		ginx.Dangerous("config not found", 404)
	}
	resp := convertToResponse(rt, provider)
	if version != resp.Version {
		models.AssetSetStatus(rt.Ctx, ident, 1)
	}
	c.JSON(200, resp)
}

// 计算配置版本，对所有key排序后进行hash，如有变化则会生成新的版本号
func calcVersion(conf map[string]map[string]*ConfigWithFormat) string {
	hash := fnv.New32a()
	keys := make([]string, 0, len(conf))
	for k := range conf {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		hash.Write([]byte(key))
		subKeys := make([]string, 0, len(conf[key]))
		for k := range conf[key] {
			subKeys = append(subKeys, k)
		}
		sort.Strings(subKeys)

		for _, subKey := range subKeys {
			hash.Write([]byte(subKey))
			subValue := conf[key][subKey]
			hash.Write([]byte(subValue.Format))
			hash.Write([]byte(subValue.Config))
		}
	}

	hashValue := hash.Sum32()

	return strconv.Itoa(int(hashValue))
}

func convertToResponse(rt *Router, assets []*models.Asset) httpRemoteProviderResponse {
	resp := httpRemoteProviderResponse{
		Configs: make(map[string]map[string]*ConfigWithFormat),
	}

	for _, asset := range assets {
		if asset.Plugin == "host" {
			// 对于主机监控，解析toml文件, 拆分为多个插件配置
			configs := make(map[string]map[string]interface{})
			toml.Unmarshal([]byte(asset.Configs), &configs)

			for k, v := range configs {
				if v == nil {
					v = make(map[string]interface{})
				}
				labels := make(map[string]string)
				labels["asset_type"] = asset.Type
				labels["asset_id"] = fmt.Sprintf("%d", asset.Id)
				group := rt.BusiGroupCache.GetByBusiGroupId(asset.GroupId)
				if group != nil {
					if group.LabelEnable == 1 {
						labels["busigroup"] = group.LabelValue
					} else {
						labels["busigroup"] = group.Name
					}
				}
				labels["instance"] = asset.Label

				//将tags写入到label里，给指标打标
				for _, tags := range strings.Fields(asset.Tags) {
					attr := strings.Split(tags, "=")
					if len(attr) > 1 {
						labels[attr[0]] = attr[1]
					}
				}

				v["labels"] = labels

				cfg, err := toml.Marshal(v)
				if err != nil {
					logger.Warningf("parse config %s error: %s", k, err)
					continue
				}

				conf := &ConfigWithFormat{
					Format: TomlFormat,
					Config: string(cfg),
				}
				checksum := fmt.Sprintf("%s-%x", k, md5.Sum(cfg))
				resp.Configs[k] = map[string]*ConfigWithFormat{checksum: conf}
			}
		} else {
			//其他插件不用解析
			labels := make(map[string]string)
			labels["asset_type"] = asset.Type
			labels["asset_id"] = fmt.Sprintf("%d", asset.Id)
			group := rt.BusiGroupCache.GetByBusiGroupId(asset.GroupId)
			if group != nil {
				if group.LabelEnable == 1 {
					labels["busigroup"] = group.LabelValue
				} else {
					labels["busigroup"] = group.Name
				}
			}
			labels["instance"] = asset.Label

			//将tags写入到label里，给指标打标
			for _, tags := range strings.Fields(asset.Tags) {
				attr := strings.Split(tags, "=")
				if len(attr) > 1 {
					labels[attr[0]] = attr[1]
				}
			}

			var configs map[string]interface{}
			toml.Unmarshal([]byte(asset.Configs), &configs)
			configs["instances"].([]interface{})[0].(map[string]interface{})["labels"] = labels

			cfg, err := toml.Marshal(configs)
			if err != nil {
				logger.Warningf("parse config %s error: %s", asset.Type, err)
				continue
			}
			conf := &ConfigWithFormat{
				Format: TomlFormat,
				Config: string(cfg),
			}
			checksum := fmt.Sprintf("%s-%x", asset.Plugin, md5.Sum(cfg))
			if resp.Configs[asset.Plugin] == nil {
				resp.Configs[asset.Plugin] = make(map[string]*ConfigWithFormat)
			}
			resp.Configs[asset.Plugin][checksum] = conf
		}
	}

	resp.Version = calcVersion(resp.Configs)

	return resp
}
