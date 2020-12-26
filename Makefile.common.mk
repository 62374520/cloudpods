ifeq ($(__inc_Makefile_common_mk),)
__inc_Makefile_common_mk:=1

DockerImageRegistry?=registry.cn-beijing.aliyuncs.com
DockerImageAlpineBuild?=$(DockerImageRegistry)/yunionio/alpine-build:1.0-5
DockerImageCentOSBuild?=$(DockerImageRegistry)/yunionio/centos-build:1.1-3


EnvIf=$(if $($(1)),$(1)=$($(1)))

define dockerCentOSBuildCmd
set -o xtrace
set -o errexit
set -o pipefail
cd /root/go/src/yunion.io/x/onecloud
env \
	$(call EnvIf,GOARCH) \
	$(call EnvIf,GOOS) \
	$(call EnvIf,CGO_ENABLED) \
	make $(1)
chown -R $(shell id -u):$(shell id -g) _output
endef

docker-centos-build: export dockerCentOSBuildCmd:=$(call dockerCentOSBuildCmd,$(F))
docker-centos-build:
	docker rm --force onecloud-docker-centos-build &>/dev/null || true
	docker run \
		--rm \
		--name onecloud-docker-centos-build \
		-v $(CURDIR):/root/go/src/yunion.io/x/onecloud \
		-v $(CURDIR)/_output/centos-build:/root/go/src/yunion.io/x/onecloud/_output \
		-v $(CURDIR)/_output/centos-build/_cache:/root/.cache \
		$(DockerImageCentOSBuild) \
		/bin/bash -c "$$dockerCentOSBuildCmd"
	ls -lh _output/centos-build/bin

# NOTE we need a way to stop and remove the container started by docker-build.
# No --tty, --stop-signal won't work
docker-centos-build-stop:
	docker stop --time 0 onecloud-docker-centos-build || true

.PHONY: docker-centos-build
.PHONY: docker-centos-build-stop


define dockerAlpineBuildCmd
set -o xtrace
set -o errexit
set -o pipefail
cd /root/go/src/yunion.io/x/onecloud
env \
	$(call EnvIf,GOARCH) \
	$(call EnvIf,GOOS) \
	$(call EnvIf,CGO_ENABLED) \
	make $(1)
chown -R $(shell id -u):$(shell id -g) _output
endef

docker-alpine-build: export dockerAlpineBuildCmd:=$(call dockerAlpineBuildCmd,$(F))
docker-alpine-build:
	docker rm --force onecloud-docker-alpine-build &>/dev/null || true
	docker run \
		--rm \
		--name onecloud-docker-alpine-build \
		-v $(CURDIR):/root/go/src/yunion.io/x/onecloud \
		-v $(CURDIR)/_output/alpine-build:/root/go/src/yunion.io/x/onecloud/_output \
		-v $(CURDIR)/_output/alpine-build/_cache:/root/.cache \
		$(DockerImageAlpineBuild) \
		/bin/sh -c "$$dockerAlpineBuildCmd"
	ls -lh _output/alpine-build/bin

docker-alpine-build-stop:
	docker stop --time 0 onecloud-docker-alpine-build || true

.PHONY: docker-alpine-build
.PHONY: docker-alpine-build-stop

endif # __inc_Makefile_common_mk
