package models

import (
	"fmt"
	"os"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type TargetVersion struct {
	Version string
	Path    string
	Os      string
	Arch    string
}

func TargetVersionGet(ctx *ctx.Context, version string, goos string, arch string) (*TargetVersion, error) {
	win_ext := ""
	if goos == "windows" {
		win_ext = ".exe"
	}
	path := fmt.Sprintf("etc/client/categraf-%s-%s-%s%s.gz", version, goos, arch, win_ext)

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	tv := &TargetVersion{
		Version: version,
		Path:    path,
		Os:      goos,
		Arch:    arch,
	}

	return tv, err
}

func UpdataVersion(ctx *ctx.Context, m map[string]interface{}, ident string) error {
	return DB(ctx).Table("target").Where("ident = ?", ident).Updates(m).Error
}
