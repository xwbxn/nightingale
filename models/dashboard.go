package models

import (
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/pkg/errors"
	"github.com/toolkits/pkg/str"
	"gorm.io/gorm"
)

type Dashboard struct {
	Id       int64    `json:"id" gorm:"primaryKey"`
	GroupId  int64    `json:"group_id"`
	Name     string   `json:"name"`
	Tags     string   `json:"-"`
	TagsLst  []string `json:"tags" gorm:"-"`
	Configs  string   `json:"configs"`
	CreateAt int64    `json:"create_at"`
	CreateBy string   `json:"create_by"`
	UpdateAt int64    `json:"update_at"`
	UpdateBy string   `json:"update_by"`
}

type DashboardUser struct {
	Id        int64          `gorm:"column:id;primaryKey" json:"id" `                          //type:BIGINT       comment:主键        version:2023-9-31 09:11
	UserId    int64          `gorm:"column:user_id" json:"user_id" `                           //type:BIGINT       comment:用户id      version:2023-9-31 09:11
	AssetsId  int64          `gorm:"column:assets_id" json:"assets_id" `                       //type:BIGINT       comment:资产id      version:2023-9-31 09:11
	Type      string         `gorm:"column:type" json:"type" `                                 //type:string       comment:资产类型    version:2023-9-31 09:11
	PageName  string         `gorm:"column:page_name" json:"page_name" `                       //type:string       comment:页签        version:2023-9-31 11:18
	Sort      int64          `gorm:"column:sort" json:"sort" `                                 //type:*int         comment:序号        version:2023-9-31 09:11
	CreatedBy string         `gorm:"column:created_by" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人      version:2023-9-31 09:11
	CreatedAt int64          `gorm:"column:created_at" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间    version:2023-9-31 09:11
	UpdatedBy string         `gorm:"column:updated_by" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人      version:2023-9-31 09:11
	UpdatedAt int64          `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间    version:2023-9-31 09:11
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" `                     //type:*time.Time   comment:删除时间    version:2023-9-31 09:11
}

func (d *Dashboard) TableName() string {
	return "dashboard"
}

func (d *Dashboard) Verify() error {
	if d.Name == "" {
		return errors.New("Name is blank")
	}

	if str.Dangerous(d.Name) {
		return errors.New("Name has invalid characters")
	}

	return nil
}

func (d *Dashboard) Add(ctx *ctx.Context) error {
	if err := d.Verify(); err != nil {
		return err
	}

	exists, err := DashboardExists(ctx, "group_id=? and name=?", d.GroupId, d.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to count dashboard")
	}

	if exists {
		return errors.New("Dashboard already exists")
	}

	now := time.Now().Unix()
	d.CreateAt = now
	d.UpdateAt = now

	return Insert(ctx, d)
}

func (d *Dashboard) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	if err := d.Verify(); err != nil {
		return err
	}

	return DB(ctx).Model(d).Select(selectField, selectFields...).Updates(d).Error
}

func (d *Dashboard) Del(ctx *ctx.Context) error {
	cgids, err := ChartGroupIdsOf(ctx, d.Id)
	if err != nil {
		return err
	}

	if len(cgids) == 0 {
		return DB(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id=?", d.Id).Delete(&Dashboard{}).Error; err != nil {
				return err
			}
			return nil
		})
	}

	return DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id in ?", cgids).Delete(&Chart{}).Error; err != nil {
			return err
		}

		if err := tx.Where("dashboard_id=?", d.Id).Delete(&ChartGroup{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id=?", d.Id).Delete(&Dashboard{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func DashboardGet(ctx *ctx.Context, where string, args ...interface{}) (*Dashboard, error) {
	var lst []*Dashboard
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	lst[0].TagsLst = strings.Fields(lst[0].Tags)

	return lst[0], nil
}

func DashboardCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Dashboard{}).Where(where, args...))
}

func DashboardExists(ctx *ctx.Context, where string, args ...interface{}) (bool, error) {
	num, err := DashboardCount(ctx, where, args...)
	return num > 0, err
}

func DashboardGets(ctx *ctx.Context, groupId int64, query string) ([]Dashboard, error) {
	session := DB(ctx).Where("group_id=?", groupId).Order("name")

	arr := strings.Fields(query)
	if len(arr) > 0 {
		for i := 0; i < len(arr); i++ {
			if strings.HasPrefix(arr[i], "-") {
				q := "%" + arr[i][1:] + "%"
				session = session.Where("name not like ? and tags not like ?", q, q)
			} else {
				q := "%" + arr[i] + "%"
				session = session.Where("(name like ? or tags like ?)", q, q)
			}
		}
	}

	var objs []Dashboard
	err := session.Select("id", "group_id", "name", "tags", "create_at", "create_by", "update_at", "update_by").Find(&objs).Error
	if err == nil {
		for i := 0; i < len(objs); i++ {
			objs[i].TagsLst = strings.Fields(objs[i].Tags)
		}
	}

	return objs, err
}

func DashboardGetsByIds(ctx *ctx.Context, ids []int64) ([]Dashboard, error) {
	if len(ids) == 0 {
		return []Dashboard{}, nil
	}

	var lst []Dashboard
	err := DB(ctx).Where("id in ?", ids).Order("name").Find(&lst).Error
	return lst, err
}

func DashboardGetAll(ctx *ctx.Context) ([]Dashboard, error) {
	var lst []Dashboard
	err := DB(ctx).Find(&lst).Error
	return lst, err
}

//删除用户看板
func DeleteByUserAndPageName(tx *gorm.DB, userId int64, pageNam string) error {
	err := tx.Debug().Where("user_id = ? AND page_name = ?", userId, pageNam).Delete(&DashboardUser{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

//添加数据看板
func AddDashBoardUser(tx *gorm.DB, lst []DashboardUser) error {
	err := tx.Debug().Create(&lst).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

//根据userId和pageName分页查询
func DashBoardUserPageNameByUser(ctx *ctx.Context, userId int64) ([]string, error) {
	var lst []string

	err := DB(ctx).Debug().Model(&DashboardUser{}).Where("user_id = ?", userId).Pluck("page_name", &lst).Error
	return lst, err
}

//根据userId和pageName分页查询
func DashBoardUserByUserAndPageName(ctx *ctx.Context, userId int64, pageNam string, limit, offset int) ([]DashboardUser, error) {
	var lst []DashboardUser

	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("sort")
	}

	err := session.Debug().Model(&DashboardUser{}).Where("user_id = ? AND page_name = ?", userId, pageNam).Find(&lst).Error
	return lst, err
}

//根据userId和pageName统计个数
func DashBoardUserCountByUserAndPageName(ctx *ctx.Context, userId int64, pageNam string) (num int64, err error) {
	err = DB(ctx).Debug().Model(&DashboardUser{}).Where("user_id = ? AND page_name = ?", userId, pageNam).Count(&num).Error
	return num, err
}
