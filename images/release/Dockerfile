FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.19-openshift-4.13 AS builder
WORKDIR /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator
COPY . .
RUN make build

FROM registry.ci.openshift.org/ocp/4.13:base
COPY --from=builder /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator/migrator /usr/bin/
