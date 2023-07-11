all: build
.PHONY: all

GO_BUILD_PACKAGES :=./cmd/migrator ./cmd/trigger
GO_TEST_PACKAGES :=./pkg/...

include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/images.mk \
)

# generate image targets
IMAGE_REGISTRY :=registry.svc.ci.openshift.org
$(call build-image,kube-storage-version-migrator,$(IMAGE_REGISTRY)/ocp/4.3:kube-storage-version-migrator,./images/release/Dockerfile,.)

$(call verify-golang-versions,images/release/Dockerfile)
