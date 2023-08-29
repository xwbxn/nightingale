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
	path := fmt.Sprintf("etc/client/categraf-%s-%s-%s", version, goos, arch)
	if goos == "windows" {
		path += ".exe"
	}

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
