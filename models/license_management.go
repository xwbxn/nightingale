// Package models  许可配置
// date : 2023-10-29 15:12
// desc : 许可配置
package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// LicenseConfig  许可配置。
// 说明:
// 表名:license_config
// group: LicenseConfig
// version:2023-10-29 15:12
type LicenseConfig struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int         comment:主键          version:2023-10-29 15:12
	Days      int64          `gorm:"column:DAYS" json:"days" `                                 //type:*int         comment:剩余天数      version:2023-10-29 15:12
	Nodes     int64          `gorm:"column:NODES" json:"nodes" `                               //type:*int         comment:剩余节点数    version:2023-10-29 15:12
	Frequency string         `gorm:"column:FREQUENCY" json:"frequency" `                       //type:string       comment:提醒频率      version:2023-10-29 15:12
	Email     string         `gorm:"column:EMAIL" json:"email" `                               //type:string       comment:邮箱          version:2023-10-29 15:12
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人        version:2023-10-29 15:12
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间      version:2023-10-29 15:12
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人        version:2023-10-29 15:12
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间      version:2023-10-29 15:12
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间      version:2023-10-29 15:12
}

type License struct {
	Id             int64  `json:"id"`
	SerialNumber   string `json:"serial_number"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	PermissionNode int64  `json:"permission_node"`
	UsedNode       int64  `json:"used_node"`
	TargetVersion  string `json:"target_version"`
}

var PASSWORD = []byte("mypassword")

// TableName 表名:license_config，许可配置。
// 说明:
func (l *LicenseConfig) TableName() string {
	return "license_config"
}

// 按id查询
func LicenseConfigGetById(ctx *ctx.Context, id int64) (*LicenseConfig, error) {
	var obj *LicenseConfig
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func LicenseConfigGets(ctx *ctx.Context) (*LicenseConfig, error) {
	var lst []*LicenseConfig
	err := DB(ctx).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	return lst[0], nil
}

// 增加许可配置
func (l *LicenseConfig) Add(ctx *ctx.Context) error {
	// 这里写LicenseConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(l).Error
}

// 更新许可配置
func (l *LicenseConfig) Update(ctx *ctx.Context, updateFrom LicenseConfig) error {
	// 这里写LicenseConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(l).Updates(updateFrom).Error
}

// 更新许可配置Map
func (l *LicenseConfig) LicenseConfigUpdateMap(ctx *ctx.Context, updateFrom LicenseConfig) error {
	// 这里写LicenseConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(l).Updates(updateFrom).Error
}

func LicenseConfigStatistics(ctx *ctx.Context) (*Statistics, error) {
	session := DB(ctx).Model(&LicenseConfig{}).Select("count(*) as total", "max(UPDATED_AT) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

// 查询所有
func LicenseConfigGetsAll(ctx *ctx.Context) ([]*LicenseConfig, error) {
	var lst []*LicenseConfig
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

func LicenseGets() ([]*License, error) {
	licenses := make([]*License, 0)
	fileInfos, err := ioutil.ReadDir("etc/license/crt/")
	if err != nil {
		return licenses, errors.New("文件读取失败")
	}

	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), ".bak") || strings.Contains(fileInfo.Name(), "readme.txt") {
			continue
		}
		var license License
		crtData, err := ioutil.ReadFile("etc/license/crt/" + fileInfo.Name())
		if err != nil {
			return licenses, errors.New("文件读取失败")
		}
		keyData, err := ioutil.ReadFile("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key")
		if err != nil {
			return licenses, errors.New("文件读取失败")
		}
		// 解密证书并获取证书
		certBlock, _ := pem.Decode(crtData)
		// if certBlock == nil || certBlock.Type != "CERTIFICAET" {
		// 	ginx.Bomb(http.StatusOK, "证书文件格式错误")
		// 	return
		// }

		certificate, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return licenses, errors.New("证书格式化失败")
		}
		id, _ := strconv.ParseInt(strings.Split(strings.Split(fileInfo.Name(), ".")[0], "-")[1], 10, 64)
		license.Id = id
		license.StartTime = certificate.NotBefore.Unix()
		license.EndTime = certificate.NotAfter.Unix()

		var value []byte
		for _, val := range certificate.Extensions {
			if val.Id.String() == "1.2.3.4" {

				_, err := asn1.Unmarshal(val.Value, &value)
				if err != nil {
					return licenses, errors.New("额外属性解析失败")
				}
				break
			}
		}

		// 解密私钥PEM并获取私钥
		keytBlock, _ := pem.Decode(keyData)
		// if keytBlock == nil || keytBlock.Type != "RSA PRIVATE KEY" {
		// 	ginx.Bomb(http.StatusOK, "私钥文件格式错误")
		// }
		decryptedPrivateKey, err := x509.DecryptPEMBlock(keytBlock, PASSWORD)
		if err != nil {
			return licenses, errors.New("私钥解密失败")
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(decryptedPrivateKey)
		if err != nil {
			return licenses, errors.New("私钥格式化失败")
		}
		decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, value)
		if err != nil {
			return licenses, errors.New("解密失败")
		}
		md := make(map[string]interface{})
		err = json.Unmarshal(decryptedText, &md)
		if err != nil {
			return licenses, errors.New("额外信息解析失败")
		}
		// logger.Debug(md)
		license.SerialNumber = md["serialNumber"].(string)
		license.PermissionNode = int64(md["permissionNode"].(float64))
		licenses = append(licenses, &license)
	}
	sort.Slice(licenses, func(i, j int) bool {
		return licenses[i].Id < licenses[j].Id
	})
	return licenses, nil
}

func LicenseCacheGets() ([]*License, error) {
	licenses, err := LicenseGets()
	if err != nil {
		return nil, err
	}
	return licenses, nil
	// var license License
	// license.Id = 0
	// for _, val := range licenses {
	// 	if val.StartTime < time.Now().Unix() && val.EndTime > time.Now().Unix() && val.Id > license.Id {
	// 		license.SerialNumber = val.SerialNumber
	// 		license.StartTime = val.StartTime
	// 		license.EndTime = val.EndTime
	// 		license.PermissionNode = val.PermissionNode
	// 		license.Id = val.Id
	// 		license.TargetVersion = val.TargetVersion
	// 	}
	// }
	// res := make([]*License, 0)
	// res = append(res, &license)
	// return res, nil
}

func LicenseStatistics() (*Statistics, error) {
	licenses, err := LicenseGets()
	if err != nil {
		return nil, err
	}
	var maxDate int64 = 0
	for _, val := range licenses {
		if val.Id > maxDate {
			maxDate = val.Id
		}
	}

	var stats Statistics
	stats.Total = int64(len(licenses))
	stats.LastUpdated = maxDate
	return &stats, nil
}
