.PHONY: prebuild build

ROOT:=$(shell pwd -P)
GIT_COMMIT:=$(shell git --work-tree ${ROOT}  rev-parse 'HEAD^{commit}')
_GIT_VERSION:=$(shell git --work-tree ${ROOT} describe --tags --abbrev=14 "${GIT_COMMIT}^{commit}" 2>/dev/null)
TAG=$(shell echo "${_GIT_VERSION}" |  awk -F"-" '{print $$1}')
RELEASE_VERSION:="$(TAG)"

all: build

build:
	go build -ldflags "-w -s -X github.com/ccfos/nightingale/v6/pkg/version.Version=$(RELEASE_VERSION)" -o n9e ./cmd/center/main.go

build-edge:
	go build -ldflags "-w -s -X github.com/ccfos/nightingale/v6/pkg/version.Version=$(RELEASE_VERSION)" -o n9e-edge ./cmd/edge/

build-alert:
	go build -ldflags "-w -s -X github.com/ccfos/nightingale/v6/pkg/version.Version=$(RELEASE_VERSION)" -o n9e-alert ./cmd/alert/main.go

build-pushgw:
	go build -ldflags "-w -s -X github.com/ccfos/nightingale/v6/pkg/version.Version=$(RELEASE_VERSION)" -o n9e-pushgw ./cmd/pushgw/main.go

build-cli: 
	go build -ldflags "-w -s -X github.com/ccfos/nightingale/v6/pkg/version.Version=$(RELEASE_VERSION)" -o n9e-cli ./cmd/cli/main.go

run:
	nohup ./n9e > n9e.log 2>&1 &

run-alert:
	nohup ./n9e-alert > n9e-alert.log 2>&1 &

run-pushgw:
	nohup ./n9e-pushgw > n9e-pushgw.log 2>&1 &

release: build
	rm -rf ./release
	mkdir -p ./release/etc
	cp ./etc/assets.yaml ./release/etc/assets.yaml
	cp ./etc_prod/config.toml ./release/etc/config.toml
	cp ./docker/initsql ./release -r
	cp ./docker/docker-compose-release.yaml ./release/docker-compose.yaml
	cp ./etc/default ./release/etc -r
	cp ./docker/mysqletc ./release -r
	cp ./docker/redis ./release -r
	cp ./pub ./release -r
	cp ./dataroom ./release -r
	cp ./thirdparty ./release -r
	cp ./categraf-server ./release -r
	cp n9e release
	mv release monset
	tar -czvf monset-$(RELEASE_VERSION).tar.gz monset
	mv monset release

package-test: build
	rm -rf ./test/etc
	mkdir -p ./test/etc
	cp ./etc/assets.yaml ./test/etc/assets.yaml
	cp ./docker/initsql ./test -r
	cp ./etc/default ./test/etc -r
	cp n9e test
	
version: build
	gzip -f n9e