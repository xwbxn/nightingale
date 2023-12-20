package license

import (
	"context"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/alert/aconf"
	"github.com/ccfos/nightingale/v6/alert/sender"
	"github.com/ccfos/nightingale/v6/center/ws"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
)

type Scheduler struct {
	// key: hash
	license           *memsto.LicenseCache
	notifyConfigCache *memsto.NotifyConfigCacheType
	isSend            bool
}

func NewScheduler(license *memsto.LicenseCache, notifyConfigCache *memsto.NotifyConfigCacheType) *Scheduler {
	scheduler := &Scheduler{
		license:           license,
		notifyConfigCache: notifyConfigCache,
		isSend:            false,
	}

	go scheduler.LoopSyncLicense(context.Background())
	return scheduler
}

func (s *Scheduler) LoopSyncLicense(ctx context.Context) {
	// if s.frequency == "once" {

	// } else if s.frequency == "days" {
	duration := 24 * time.Hour
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.syncLicense()
		}
	}
	// }
}

func (s *Scheduler) syncLicense() {
	now := time.Now().Unix()
	license := s.license.GetByLicense()
	licenseConfig := s.license.GetByLicenseConfig()
	stmp := s.notifyConfigCache.GetSMTP()
	if license == nil {
		//不存在符合当前时间的license
		ws.SetMessage(758493, "当前证书已过期")
	} else {
		if licenseConfig.Frequency == "once" && !s.isSend {
			//发一次邮件
			s.isSend = true
			send(license, licenseConfig, now, stmp)
		} else if licenseConfig.Frequency == "days" {
			//每天发一次邮件
			send(license, licenseConfig, now, stmp)
		}

	}

}

func send(license *models.License, licenseConfig *models.LicenseConfig, now int64, stmp aconf.SMTPConfig) {
	subject := ""
	content := ""
	if license.PermissionNode-license.UsedNode <= licenseConfig.Nodes {
		subject = "一体化运维系统-节点授权"
		content = "您当前的节点数即将超过了授权节点数，请及时处理"
	} else if license.EndTime-now <= licenseConfig.Days*24*3600 {
		subject = "一体化运维系统-授权到期"
		content = "您当前的授权即将到期，请及时处理"
	}
	tos := strings.Split(licenseConfig.Email, ",")

	var es *sender.EmailSender
	es.SendEmail(subject, content, tos, stmp)
}
