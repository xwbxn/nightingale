// Package models  用户配置
// date : 2023-10-13 10:09
// desc : 用户配置
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// UserConfig  用户配置。
// 说明:
// 表名:user_config
// group: UserConfig
// version:2023-10-13 10:09
type UserConfig struct {
    Id              int64  `gorm:"column:ID;primaryKey" json:"id" `                 //type:BIGINT       comment:配置主键                   version:2023-10-13 10:09
    LogLever        int  `gorm:"column:LOG_LEVER" json:"log_lever" `               //type:*int         comment:日志等级                   version:2023-10-13 10:09
    HttpHost        string  `gorm:"column:HTTP_HOST" json:"http_host" `               //type:string       comment:HTTP监听地址               version:2023-10-13 10:09
    HttpPort        int  `gorm:"column:HTTP_PORT" json:"http_port" `               //type:*int         comment:HTTP监听端口               version:2023-10-13 10:09
    Captcha         int  `gorm:"column:CAPTCHA" json:"captcha" `                  //type:*int         comment:启用验证码                 version:2023-10-13 10:09
    ApiService      int  `gorm:"column:API_SERVICE" json:"api_service" `           //type:*int         comment:APIForService              version:2023-10-13 10:09
    AccessExpired   int64  `gorm:"column:ACCESS_EXPIRED" json:"access_expired" `     //type:*int         comment:token有效期（用户登录）    version:2023-10-13 10:09
    RefreshExpired  int64  `gorm:"column:REFRESH_EXPIRED" json:"refresh_expired" `   //type:*int         comment:token有效期（长期展示）    version:2023-10-13 10:09
    OpenRsa         int  `gorm:"column:OPEN_RSA" json:"open_rsa" `                 //type:*int         comment:启用加密                   version:2023-10-13 10:09
    LoginTitle      string  `gorm:"column:LOGIN_TITLE" json:"login_title" `           //type:string       comment:登录页标题                 version:2023-10-13 10:09
    LogoTop         string  `gorm:"column:LOGO_TOP" json:"logo_top" `                 //type:string       comment:系统顶部LOGO               version:2023-10-13 10:09
    LogoTitle       string  `gorm:"column:LOGO_TITLE" json:"logo_title" `             //type:string       comment:网页标题LOGO               version:2023-10-13 10:09
    CreatedBy       string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`             //type:string       comment:创建人                     version:2023-10-13 10:09
    CreatedAt       int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`             //type:*int         comment:创建时间                   version:2023-10-13 10:09
    UpdatedBy       string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`             //type:string       comment:更新人                     version:2023-10-13 10:09
    UpdatedAt       int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`             //type:*int         comment:更新时间                   version:2023-10-13 10:09
    DeletedAt       gorm.DeletedAt  `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`             //type:*time.Time   comment:删除时间                   version:2023-10-13 10:09
}

// TableName 表名:user_config，用户配置。
// 说明:
func (u *UserConfig) TableName() string {
	return "user_config"
}

// 条件查询
func UserConfigGets(ctx *ctx.Context, query string, limit, offset int) ([]UserConfig, error) {
	session := DB(ctx)
    // 分页
	if limit > -1 {
	    session = session.Limit(limit).Offset(offset).Order("id")
	}
	
	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
	if query != "" {
		q := "%" + query + "%"
		session = session.Where("id like ?", q)
	}
    
	var lst []UserConfig
	err := session.Find(&lst).Error
	
	return lst, err
}

type LogoConfig struct {
	LoginTitle      string  `gorm:"column:LOGIN_TITLE" json:"login_title" `           //type:string       comment:登录页标题                 version:2023-10-13 10:09
	LogoTop         string  `gorm:"column:LOGO_TOP" json:"logo_top" `                 //type:string       comment:系统顶部LOGO               version:2023-10-13 10:09
	LogoTitle       string  `gorm:"column:LOGO_TITLE" json:"logo_title" `             //type:string       comment:网页标题LOGO               version:2023-10-13 10:09
}

// 查询用户Logo相关信息
func UserLogoConfigGet(ctx *ctx.Context) (*UserConfig, error){
	var lst *UserConfig
	err := DB(ctx).Debug().Model(&UserConfig{}).Where("id", 1).Find(&lst).Error

	return lst, err
}

// 按id查询
func UserConfigGetById(ctx *ctx.Context, id int64) (*UserConfig, error) {
	var obj *UserConfig
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}
    
	return obj, nil
}

func UserConfigGet(ctx *ctx.Context) (*UserConfig, error) {
	var lst []*UserConfig

	err := ctx.DB.Debug().Where("id",1).Find(&lst).Error

	if err != nil{
		return nil, err
	}

	if lst[0] == nil {
		return nil, err
	}
	return lst[0], err
}

// 按id更新
func UserConfigUpdateById(ctx *ctx.Context, id int64) error {
	var lst *UserConfig
	err := DB(ctx).Debug().Where("id",1).Find(&lst).Error
	if lst.HttpHost != "" {
		lst.HttpHost = ""
	}
	if err != nil {
		return err
	}
	err = DB(ctx).Take(&lst, id).Error
	if err != nil {
		return err
	}
    
	return nil
}

// 更新Logo信息
func (lst *UserConfig) UserLogoUpdate(ctx *ctx.Context) error {
	err := ctx.DB.Debug().Model(&UserConfig{}).Where("ID", 1).Omit("CREATED_BY","CREATED_AT","UPDATED_BY", "UPDATED_AT").Updates(lst).Error
	return err
}

// 查询图片路径
func LogoPathGet(ctx *ctx.Context) (*LogoConfig, error){
	var res *LogoConfig
	err := ctx.DB.Debug().Model(&UserConfig{}).Where("ID", 1).Find(&res).Error
	return res, err
}

// 查询所有
func UserConfigGetsAll(ctx *ctx.Context) ([]UserConfig, error) {
	var lst []UserConfig
	err := DB(ctx).Find(&lst).Error
	
	return lst, err
}

// 增加用户配置
func (u *UserConfig) Add(ctx *ctx.Context) error {
    // 这里写UserConfig的业务逻辑，通过error返回错误
    
    // 实际向库中写入
	return DB(ctx).Create(u).Error
}

// 删除用户配置
func (u *UserConfig) Del(ctx *ctx.Context) error {
    // 这里写UserConfig的业务逻辑，通过error返回错误
	
    // 实际向库中写入
	return DB(ctx).Delete(u).Error
}

// 更新用户配置
func (u *UserConfig) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	// 这里写UserConfig的业务逻辑，通过error返回错误
	
	// 实际向库中写入
	return DB(ctx).Model(u).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(u).Error
}

// 根据条件统计个数
func UserConfigCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&UserConfig{}).Where(where, args...))
}

